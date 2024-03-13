package auth

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"

	"booking-schedule/internal/app/api/auth"
	userRepository "booking-schedule/internal/app/repository/user"
	"booking-schedule/internal/app/service/jwt"
	userService "booking-schedule/internal/app/service/user"
	"booking-schedule/internal/config"
	"booking-schedule/internal/pkg/db"
	"booking-schedule/internal/pkg/db/transaction"
	tracer "booking-schedule/internal/pkg/trace"

	"go.opentelemetry.io/otel/trace"
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
	config     *config.AppConfig

	server *http.Server
	log    *slog.Logger
	tracer trace.Tracer

	userRepository userRepository.Repository
	userService    *userService.Service
	jwtService     jwt.Service

	authImpl *auth.Implementation
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

func (s *serviceProvider) GetConfig() *config.AppConfig {
	if s.config == nil {
		cfg, err := config.ReadAppConfig(s.configPath)
		if err != nil {
			log.Fatalf("could not get bookings-api config: %s", err)
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.GetDB(ctx), s.GetLogger())
		return s.userRepository
	}

	return s.userRepository
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

func (s *serviceProvider) GetAuthImpl(ctx context.Context) (*auth.Implementation, error) {
	if s.authImpl == nil {
		tracer, err := s.GetTracer(ctx)
		if err != nil {
			return nil, err
		}
		s.authImpl = auth.NewImplementation(s.GetUserService(ctx), tracer)
	}

	return s.authImpl, nil
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

func (s *serviceProvider) GetTracer(ctx context.Context) (trace.Tracer, error) {
	if s.tracer == nil {
		tracer, err := tracer.NewTracer(ctx, "http://jaeger:14268/api/traces", "auth")
		if err != nil {
			s.log.Error("failed to create tracer: %s", err)
			return nil, err
		}
		s.tracer = tracer
	}

	return s.tracer, nil
}
