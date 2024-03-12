package bookings

import (
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	mwLogger "booking-schedule/internal/middleware/logger"
	"context"
	"errors"
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
	pathCert        string
	pathKey         string
	serviceProvider *serviceProvider
	router          *chi.Mux
}

// NewApp ...
func NewApp(ctx context.Context, pathConfig string, pathCert string, pathKey string) (*App, error) {
	a := &App{
		pathConfig: pathConfig,
		pathCert:   pathCert,
		pathKey:    pathKey,
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
		//nolint
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
	impl := a.serviceProvider.GetBookingImpl(ctx)

	address, err := a.serviceProvider.GetConfig().GetAddress()
	if err != nil {
		return err
	}
	a.serviceProvider.GetLogger().Info("initializing server", slog.String("address", address))
	a.serviceProvider.GetLogger().Debug("logger debug mode enabled")

	a.router = chi.NewRouter()
	a.router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	//a.router.Use(middleware.Logger)    // Логирование всех запросов
	a.router.Use(mwLogger.New(a.serviceProvider.GetLogger()))
	a.router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть

	a.router.Route("/bookings", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/sign-up", impl.SignUp(a.serviceProvider.GetLogger()))
			r.Get("/sign-in", impl.SignIn(a.serviceProvider.GetLogger()))
			//r.Group(func(r chi.Router) {r.Use(authHandler) r.Get("/me", impl.GetMyUser(a.serviceProvider.GetLogger()))})

		})
		r.Get("/get-vacant-rooms", impl.GetVacantRooms(a.serviceProvider.GetLogger()))
		r.Get("/{suite_id}/get-vacant-dates", impl.GetVacantDates(a.serviceProvider.GetLogger()))
		r.Group(func(r chi.Router) {
			r.Use(auth.Auth(a.serviceProvider.GetLogger(), a.serviceProvider.GetJWTService()))
			r.Post("/add", impl.AddBooking(a.serviceProvider.GetLogger()))
			r.Get("/get-bookings", impl.GetBookings(a.serviceProvider.GetLogger()))
			r.Route("/{booking_id}", func(r chi.Router) {
				r.Get("/get", impl.GetBooking(a.serviceProvider.GetLogger()))
				r.Patch("/update", impl.UpdateBooking(a.serviceProvider.GetLogger()))
				r.Delete("/delete", impl.DeleteBooking(a.serviceProvider.GetLogger()))
			})
		})

	})

	return nil
}

func (a *App) startServer() error {
	srv := a.serviceProvider.getServer(a.router)
	if srv == nil {
		a.serviceProvider.GetLogger().Error("server was not initialized")
		return errors.New("server was not initialized")
	}
	a.serviceProvider.GetLogger().Info("starting server", slog.String("address", srv.Addr))

	done := make(chan os.Signal, 1)
	errChan := make(chan error)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		switch a.serviceProvider.GetConfig().GetEnv() {
		case envProd:
			err := a.initCertificates()
			if err != nil {
				a.serviceProvider.GetLogger().Error("failed to initialize certificates", sl.Err(err))
				errChan <- err
			}

			if err = srv.ListenAndServeTLS(a.pathCert, a.pathKey); err != nil {
				a.serviceProvider.GetLogger().Error("", sl.Err(err))
				errChan <- err
			}
		default:
			if err := srv.ListenAndServe(); err != nil {
				a.serviceProvider.GetLogger().Error("", sl.Err(err))
				errChan <- err
			}
		}
	}()

	a.serviceProvider.GetLogger().Info("server started")

	select {
	case err := <-errChan:
		return err
	case <-done:
		a.serviceProvider.GetLogger().Info("stopping server")
		// TODO: move timeout to config
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			a.serviceProvider.GetLogger().Error("failed to stop server", sl.Err(err))
			return err
		}

		a.serviceProvider.GetLogger().Info("server stopped")
	}

	return nil
}
