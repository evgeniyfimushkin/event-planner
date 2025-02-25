package main

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handler"
	"auth-service/internal/http-server/middlewarelogger"
	"auth-service/internal/lib/logger"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main(){
    cfg := config.MustLoadConfig()
    log := logger.SetupLogger(cfg.Env)
    fmt.Print(cfg.PrivateKey)

    log.Info("Connecting to db with params: ")
    log.Info("Database: ", slog.String("host", cfg.Database.Host), slog.String("port", cfg.Database.Port))

    dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",cfg.Database.User, cfg.Database.Password, "authdb", cfg.Database.Host, cfg.Database.Port, "disable")
    dbConnection := db.SetupDB(dsn)

    userRepo := repository.NewUserRepository(dbConnection)
    _ = userRepo
    userRepo.Create(&models.User{Username: "ivan",Email: "ivan@gmail",PassHash: "123"})

    loginService, err := service.NewLoginService(userRepo, cfg.PrivateKey, cfg.TokenTTL)
    registerService, err := service.NewRegisterService(userRepo)
    if err != nil {
        log.Error("failed to init login service", logger.Err(err))
    }

    router := chi.NewRouter()
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middlewarelogger.New(log))
    router.Use(middleware.Recoverer)
    router.Use(middleware.URLFormat)
    router.Post("/api/v1/auth/login", handler.Login(loginService))
    registerLimiter := httprate.LimitByIP(5, 1*time.Minute)
    router.With(registerLimiter).Post("/api/v1/auth/register", handler.Register(registerService))

    srv := &http.Server{
        Addr: fmt.Sprintf("%s:%d",cfg.Server.Addr, cfg.Server.Port),
        Handler: router,
        ReadTimeout: cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout: cfg.Server.IdleTimeout,
    }

    log.Info(fmt.Sprintf("Server listening on port %d", cfg.Server.Port))
    if err := srv.ListenAndServe(); err != nil {
        log.Error("failed to start server", logger.Err(err))
    }
    log.Info("server stopped")

}
