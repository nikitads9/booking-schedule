package bookings

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"

	"booking-schedule/internal/app/api/booking"
	"booking-schedule/internal/app/api/user"
	bookingRepository "booking-schedule/internal/app/repository/booking"
	userRepository "booking-schedule/internal/app/repository/user"
	bookingService "booking-schedule/internal/app/service/booking"
	"booking-schedule/internal/app/service/jwt"
	userService "booking-schedule/internal/app/service/user"
	"booking-schedule/internal/config"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/pkg/db"
	"booking-schedule/internal/pkg/db/transaction"
	"booking-schedule/internal/pkg/observability"

	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/trace"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	configPath string
	configType string
	config     *config.BookingConfig

	db        db.Client
	txManager db.TxManager

	server *http.Server
	log    *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter

	bookingRepository bookingRepository.Repository
	bookingService    *bookingService.Service

	userRepository userRepository.Repository
	userService    *userService.Service

	jwtService jwt.Service

	bookingImpl *booking.Implementation
	userImpl    *user.Implementation
}

func newServiceProvider(configType string, configPath string, meter metric.Meter) *serviceProvider {
	return &serviceProvider{
		configType: configType,
		configPath: configPath,
		meter:      meter,
	}
}

func (s *serviceProvider) GetDB(ctx context.Context) db.Client {
	if s.db == nil {
		cfg, err := s.GetConfig().GetDBConfig()
		if err != nil {
			s.log.Error("could not get db config: %s", sl.Err(err))
		}
		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			s.log.Error("coud not connect to db: %s", sl.Err(err))
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.BookingConfig {
	if s.config == nil {
		if s.configType == "env" {
			cfg, err := config.ReadBookingConfigEnv()
			if err != nil {
				log.Fatalf("could not get bookings-api config from env: %s", err)
			}
			s.config = cfg
		} else {
			cfg, err := config.ReadBookingConfigFile(s.configPath)
			if err != nil {
				log.Fatalf("could not get bookings-api config from file: %s", err)
			}
			s.config = cfg
		}
	}

	return s.config
}

func (s *serviceProvider) GetBookingRepository(ctx context.Context) bookingRepository.Repository {
	if s.bookingRepository == nil {
		s.bookingRepository = bookingRepository.NewBookingRepository(s.GetDB(ctx), s.GetLogger(), s.GetTracer(ctx))
		return s.bookingRepository
	}

	return s.bookingRepository
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.GetDB(ctx), s.GetLogger(), s.GetTracer(ctx))
		return s.userRepository
	}

	return s.userRepository
}

func (s *serviceProvider) GetBookingService(ctx context.Context) *bookingService.Service {
	if s.bookingService == nil {
		bookingRepository := s.GetBookingRepository(ctx)
		s.bookingService = bookingService.NewBookingService(bookingRepository, s.GetJWTService(ctx), s.GetLogger(), s.TxManager(ctx), s.GetTracer(ctx))
	}

	return s.bookingService
}

func (s *serviceProvider) GetUserService(ctx context.Context) *userService.Service {
	if s.userService == nil {
		userRepository := s.GetUserRepository(ctx)
		s.userService = userService.NewUserService(userRepository, s.GetJWTService(ctx), s.GetLogger(), s.GetTracer(ctx))
	}

	return s.userService
}

func (s *serviceProvider) GetJWTService(ctx context.Context) jwt.Service {
	if s.jwtService == nil {
		s.jwtService = jwt.NewJWTService(s.GetConfig().GetJWTConfig().Secret, s.GetConfig().GetJWTConfig().Expiration, s.GetLogger(), s.GetTracer(ctx))
	}

	return s.jwtService
}

func (s *serviceProvider) GetBookingImpl(ctx context.Context) *booking.Implementation {
	if s.bookingImpl == nil {
		s.bookingImpl = booking.NewImplementation(s.GetBookingService(ctx), s.GetTracer(ctx))
	}

	return s.bookingImpl
}

func (s *serviceProvider) GetUserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.GetUserService(ctx), s.GetTracer(ctx))
	}

	return s.userImpl
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

func (s *serviceProvider) GetTracer(ctx context.Context) trace.Tracer {
	if s.tracer == nil {
		tracer, err := observability.NewTracer(ctx, s.GetConfig().GetTracerConfig().EndpointURL, "bookings", s.GetConfig().GetTracerConfig().SamplingRate)
		if err != nil {
			s.GetLogger().Error("failed to create tracer: ", sl.Err(err))
			return nil
		}

		s.tracer = tracer

	}

	return s.tracer
}

func (s *serviceProvider) GetMeter(ctx context.Context) metric.Meter {
	if s.meter == nil {
		meter, err := observability.NewMeter(ctx, "bookings")
		if err != nil {
			s.GetLogger().Error("failed to create meter: ", sl.Err(err))
			return nil
		}

		s.meter = meter

	}

	return s.meter
}
