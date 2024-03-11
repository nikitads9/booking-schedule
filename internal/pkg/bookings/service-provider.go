package bookings

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"

	"booking-schedule/internal/app/api/handlers"
	bookingRepository "booking-schedule/internal/app/repository/booking"
	userRepository "booking-schedule/internal/app/repository/user"
	bookingService "booking-schedule/internal/app/service/booking"
	"booking-schedule/internal/app/service/jwt"
	userService "booking-schedule/internal/app/service/user"
	"booking-schedule/internal/config"
	"booking-schedule/internal/pkg/db"
	"booking-schedule/internal/pkg/db/transaction"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	db         db.Client
	txManager  db.TxManager
	configPath string
	config     *config.BookingConfig

	server *http.Server
	log    *slog.Logger

	bookingRepository bookingRepository.Repository
	bookingService    *bookingService.Service

	userRepository userRepository.Repository
	userService    *userService.Service
	jwtService     jwt.Service

	bookingImpl *handlers.Implementation
}

func newServiceProvider(configPath string) *serviceProvider {
	return &serviceProvider{
		configPath: configPath,
	}
}

func (s *serviceProvider) GetDB(ctx context.Context) db.Client {
	if s.db == nil {
		cfg, err := s.GetConfig().GetDBConfig()
		if err != nil {
			s.log.Error("could not get db config: %s", err)
		}
		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			s.log.Error("coud not connect to db: %s", err)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.BookingConfig {
	if s.config == nil {
		cfg, err := config.ReadBookingConfig(s.configPath)
		if err != nil {
			log.Fatalf("could not get bookings-api config: %s", err)
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetBookingRepository(ctx context.Context) bookingRepository.Repository {
	if s.bookingRepository == nil {
		s.bookingRepository = bookingRepository.NewBookingRepository(s.GetDB(ctx), s.GetLogger())
		return s.bookingRepository
	}

	return s.bookingRepository
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.GetDB(ctx), s.GetLogger())
		return s.userRepository
	}

	return s.userRepository
}

func (s *serviceProvider) GetBookingService(ctx context.Context) *bookingService.Service {
	if s.bookingService == nil {
		bookingRepository := s.GetBookingRepository(ctx)
		s.bookingService = bookingService.NewBookingService(bookingRepository, s.GetJWTService(), s.GetLogger(), s.TxManager(ctx))
	}

	return s.bookingService
}

func (s *serviceProvider) GetUserService(ctx context.Context) *userService.Service {
	if s.userService == nil {
		userRepository := s.GetUserRepository(ctx)
		s.userService = userService.NewUserService(userRepository, s.GetJWTService(), s.GetLogger())
	}

	return s.userService
}

func (s *serviceProvider) GetJWTService() jwt.Service {
	if s.jwtService == nil {
		s.jwtService = jwt.NewJWTService(s.GetConfig().GetJWTConfig().Secret, s.GetConfig().GetJWTConfig().Expiration, s.GetLogger())
	}

	return s.jwtService
}

func (s *serviceProvider) GetBookingImpl(ctx context.Context) *handlers.Implementation {
	if s.bookingImpl == nil {
		s.bookingImpl = handlers.NewImplementation(s.GetBookingService(ctx), s.GetUserService(ctx))
	}

	return s.bookingImpl
}

func (s *serviceProvider) getServer(router http.Handler) *http.Server {
	if s.server == nil {
		address, err := s.GetConfig().GetAddress()
		if err != nil {
			s.log.Error("could not get server address: %s", err)
			return nil
		}
		s.server = &http.Server{
			Addr:         address,
			Handler:      router,
			ReadTimeout:  s.GetConfig().GetServerConfig().Timeout,
			WriteTimeout: s.GetConfig().GetServerConfig().Timeout,
			IdleTimeout:  s.GetConfig().GetServerConfig().IdleTimeout,
			TLSConfig: &tls.Config{
				MinVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
			},
		}
	}

	return s.server
}

func (s *serviceProvider) GetLogger() *slog.Logger {
	if s.log == nil {
		env := s.GetConfig().GetEnv()
		switch env {
		case envLocal:
			s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envDev:
			s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envProd:
			s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}

		s.log.With(slog.String("env", env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	}

	return s.log
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDB(ctx).DB())
	}

	return s.txManager
}
