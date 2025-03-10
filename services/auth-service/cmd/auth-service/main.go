package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"time"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/config"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/db"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/logger"
    "github.com/evgeniyfimushkin/event-planner/services/common/pkg/middlewarelogger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main(){
    cfg := config.MustLoadConfig()
    log := logger.SetupLogger(cfg.Env)

    log.Info("Connecting to db with params: ")
    log.Info("Database: ", slog.String("host", cfg.Database.Host), slog.String("port", cfg.Database.Port))

    dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port, "disable")
    //TODO: configure sllmode with postgres
    dbConnection := db.SetupDB(dsn, &models.User{})

    userRepo := repository.NewUserRepository(dbConnection)
    _ = userRepo

    loginService, err := service.NewLoginService(userRepo, cfg.PrivateKey, cfg.TokenTTL)
    registerService, err := service.NewRegisterService(userRepo)
    if err != nil {
        log.Error("failed to init login service", logger.Err(err))
        panic("failed to init login service")
    }

    refreshService, err := service.NewRefreshService(userRepo, cfg.PrivateKey, cfg.PublicKey)
    if err != nil {
        log.Error("failed to init refresh service", logger.Err(err))
        panic("failed to init refresh service")
    }

    router := chi.NewRouter()
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middlewarelogger.New(log))
    router.Use(middleware.Recoverer)
    router.Use(middleware.URLFormat)
    router.Post("/api/v1/auth/login", handler.Login(loginService))
    router.Get("/api/v1/auth/refresh", handler.Refresh(refreshService))
    registerLimiter := httprate.LimitByIP(5, 1*time.Minute)
    router.With(registerLimiter).Post("/api/v1/auth/register", handler.Register(registerService))


    // TODO: Oauth

    srv := &http.Server{
        Addr: fmt.Sprintf("%s:%d",cfg.Server.Addr, cfg.Server.Port),
        Handler: router,
        ReadTimeout: cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout: cfg.Server.IdleTimeout,
    }

    // go prometheus metrics
    go func (){
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":9100", nil)
    }()

    log.Info(fmt.Sprintf("Server listening on port %d", cfg.Server.Port))
    if err := srv.ListenAndServe(); err != nil {
        log.Error("failed to start server", logger.Err(err))
    }
    log.Info("server stopped")

}
