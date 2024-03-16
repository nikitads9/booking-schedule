package bookings

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	mwLogger "booking-schedule/internal/middleware/logger"
	"booking-schedule/internal/pkg/certificates"
	"context"

	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
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
		a.initserviceProvider,
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

func (a *App) initserviceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.pathConfig)

	return nil
}

// Run ...
func (a *App) Run() error {
	defer func() {
		a.serviceProvider.db.Close() //nolint:errcheck
	}()

	err := a.startServer()
	if err != nil {
		a.serviceProvider.GetLogger().Error("failed to start server: %s", err)
		return err
	}

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
	defer close(errChan)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		switch a.serviceProvider.GetConfig().GetEnv() {
		case envProd:
			err := certificates.InitCertificates()
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

func (a *App) initServer(ctx context.Context) error {
	bookingImpl := a.serviceProvider.GetBookingImpl(ctx)
	userImpl := a.serviceProvider.GetUserImpl(ctx)

	address, err := a.serviceProvider.GetConfig().GetAddress()
	if err != nil {
		return err
	}
	a.serviceProvider.GetLogger().Info("initializing server", slog.String("address", address))
	a.serviceProvider.GetLogger().Debug("logger debug mode enabled")

	a.router = chi.NewRouter()
	a.router.Use(middleware.RequestID)
	//a.router.Use(middleware.Logger)    // Логирование всех запросов
	a.router.Use(otelchi.Middleware("bookings-api", otelchi.WithChiRoutes(a.router)))
	a.router.Use(mwLogger.New(a.serviceProvider.GetLogger()))
	a.router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть

	a.router.Route("/bookings", func(r chi.Router) {
		r.Get("/ping", api.HandlePingCheck())
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(auth.Auth(a.serviceProvider.GetLogger(), a.serviceProvider.GetJWTService(ctx)))
				r.Get("/me", userImpl.GetMyProfile(a.serviceProvider.GetLogger()))
				r.Delete("/delete", userImpl.DeleteMyProfile(a.serviceProvider.GetLogger()))
				r.Patch("/edit", userImpl.EditMyProfile(a.serviceProvider.GetLogger()))
			})

		})
		r.Get("/get-vacant-rooms", bookingImpl.GetVacantRooms(a.serviceProvider.GetLogger()))
		r.Get("/{suite_id}/get-vacant-dates", bookingImpl.GetVacantDates(a.serviceProvider.GetLogger()))
		r.Group(func(r chi.Router) {
			r.Use(auth.Auth(a.serviceProvider.GetLogger(), a.serviceProvider.GetJWTService(ctx)))
			r.Post("/add", bookingImpl.AddBooking(a.serviceProvider.GetLogger()))
			r.Get("/get-bookings", bookingImpl.GetBookings(a.serviceProvider.GetLogger()))
			r.Route("/{booking_id}", func(r chi.Router) {
				r.Get("/get", bookingImpl.GetBooking(a.serviceProvider.GetLogger()))
				r.Patch("/update", bookingImpl.UpdateBooking(a.serviceProvider.GetLogger()))
				r.Delete("/delete", bookingImpl.DeleteBooking(a.serviceProvider.GetLogger()))
			})
		})

	})

	return nil
}
