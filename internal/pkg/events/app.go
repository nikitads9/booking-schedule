package events

import (
	"context"
	"errors"
	"event-schedule/internal/app/api/handlers"
	"event-schedule/internal/logger/sl"
	mwLogger "event-schedule/internal/middleware/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	pathConfig      string
	serviceProvider *serviceProvider
	router          *chi.Mux
}

// NewApp ...
func NewApp(ctx context.Context, pathConfig string) (*App, error) {
	a := &App{
		pathConfig: pathConfig,
	}
	err := a.initDeps(ctx)

	return a, err
}

func (a *App) initDeps(ctx context.Context) error {

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

// Run ...
func (a *App) Run() error {
	defer func() {
		a.serviceProvider.db.Close()
	}()

	err := a.startServer()
	if err != nil {
		a.serviceProvider.log.Error("failed to start server: %s", err)
		return err
	}

	return nil
}

func (a *App) initServer(ctx context.Context) error {
	impl := a.serviceProvider.GetEventImpl(ctx)

	address, err := a.serviceProvider.GetConfig().GetAddress()
	if err != nil {
		return err
	}
	a.serviceProvider.GetLogger().Info("initializing server", slog.String("address", address)) // Вывод параметра с адресом
	a.serviceProvider.GetLogger().Debug("logger debug mode enabled")

	a.setupRouter(impl)

	return nil
}

func (a *App) setupRouter(impl *handlers.Implementation) {
	a.router = chi.NewRouter()
	a.router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	//a.router.Use(middleware.Logger)    // Логирование всех запросов
	a.router.Use(mwLogger.New(a.serviceProvider.GetLogger()))
	a.router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть

	a.router.Route("/events/{user_id}", func(r chi.Router) {
		r.Post("/add", impl.AddEvent(a.serviceProvider.GetLogger()))
		r.Get("/get-events", impl.GetEvents(a.serviceProvider.GetLogger()))
		r.Get("/get-vacant-rooms", impl.GetVacantRooms(a.serviceProvider.GetLogger()))
		r.Get("/{suite_id}/get-vacant-dates", impl.GetVacantDates(a.serviceProvider.GetLogger()))
		r.Route("/{event_id}", func(r chi.Router) {
			r.Get("/get", impl.GetEvent(a.serviceProvider.GetLogger()))
			r.Patch("/update", impl.UpdateEvent(a.serviceProvider.GetLogger()))
			r.Delete("/delete", impl.DeleteEvent(a.serviceProvider.GetLogger()))
		})
	})
}

func (a *App) startServer() error {
	srv := a.serviceProvider.getServer(a.router)
	if srv == nil {
		a.serviceProvider.log.Error("server was not initialized")
		return errors.New("server was not initialized")
	}
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

	a.serviceProvider.log.Info("server stopped")

	return nil
}
