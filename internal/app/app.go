package app

import (
	"context"
	"event-schedule/internal/api/handlers"
	mwLogger "event-schedule/internal/app/middleware/logger"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	pathConfig      string
	serviceProvider *serviceProvider
	//server
	router *chi.Mux
}

func Start(ctx context.Context, pathConfig string) error {
	a := &App{
		pathConfig: pathConfig,
	}

	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.pathConfig)

	return nil
}

func (a *App) initServer(ctx context.Context) error {
	log := a.serviceProvider.setupLogger()
	impl := a.serviceProvider.GetScheduleImpl(ctx)
	defer a.serviceProvider.db.Close()

	address, err := a.serviceProvider.config.GetAddress()
	if err != nil {
		return err
	}
	log.Info("initializing server", slog.String("address", address)) // Вывод параметра с адресом
	log.Debug("logger debug mode enabled")

	a.setupRouter(impl)

	srv := a.serviceProvider.getServer(a.router)
	err = a.startServer(srv)
	if err != nil {
		log.Error("failed to start server %s", err)
		return err
	}

	return nil
}

func (a *App) setupRouter(impl *handlers.Implementation) {
	a.router = chi.NewRouter()
	a.router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	a.router.Use(middleware.Logger)    // Логирование всех запросов
	a.router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	a.router.Use(mwLogger.New(a.serviceProvider.log))
	//r.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// RESTy routes for "events" resource
	a.router.Route("/events/{user_id}", func(r chi.Router) {
		r.Post("/add", impl.AddEvent(a.serviceProvider.log))                                 // POST /events/u123
		r.Get("/{interval}", impl.GetEvents(a.serviceProvider.log))                          // GET /events/u123/get/{interval}
		r.Get("/get-vacant-rooms/{start}-{end}", impl.GetVacantRooms(a.serviceProvider.log)) // GET /events/u123/get-vacant-rooms
		r.Get("/{suite_id}/get-vacant-dates", impl.GetVacantDates(a.serviceProvider.log))    // GET /events/u123/get-vacant-dates

		r.Route("/{event_id}", func(r chi.Router) {
			r.Use(impl.EventCtx(a.serviceProvider.log)) // Load the *Event on the request context
			r.Get("/get", impl.GetEvent(a.serviceProvider.log))
			r.Patch("/update", impl.UpdateEvent(a.serviceProvider.log))  // PATCH /event/123/update
			r.Delete("/delete", impl.DeleteEvent(a.serviceProvider.log)) // DELETE /event/123/delete
		})

		// GET /articles/whats-up
		//r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle) */
	})
}

func (a *App) startServer(srv *http.Server) error {
	a.serviceProvider.log.Info("starting server", slog.String("address", srv.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() error {
		if err := srv.ListenAndServe(); err != nil {
			a.serviceProvider.log.Error("failed to start listener")
			return err
		}

		return nil
	}()

	a.serviceProvider.log.Info("server started")

	<-done
	a.serviceProvider.log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.serviceProvider.log.Error("failed to stop server", sl.Err(err))

		return err
	}

	// TODO: close storage

	a.serviceProvider.log.Info("server stopped")

	return nil
}
