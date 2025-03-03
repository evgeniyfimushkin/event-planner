package main

import (
	"event-service/internal/handler"
	"event-service/internal/models"
	"event-service/internal/repository"
	"event-service/internal/service"
	"fmt"
	"log/slog"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/config"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/db"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/logger"
    "github.com/evgeniyfimushkin/event-planner/services/common/pkg/middlewarelogger"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"


)


func main(){
    cfg := config.MustLoadConfig()
    log := logger.SetupLogger(cfg.Env)

    log.Info("Connecting to db with params")
    log.Info("Database: ", slog.String("host", cfg.Database.Host), slog.String("port", cfg.Database.Port))

    dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port, "disable")
    dbConnection := db.SetupDB(dsn, &models.Event{})
    eventRepo := repository.NewEventRepository(dbConnection)
    verifier, err := auth.NewVerifier(cfg.PublicKey)
    if err != nil {
        log.Error("failed to init JWT verifier", logger.Err(err))
        panic("failed to init JWT verifier")
    }

    eventService := service.NewEventService(verifier, eventRepo)

    handler := handler.NewEventHandler(eventService)

    router := chi.NewRouter()
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middlewarelogger.New(log))
    router.Use(middleware.Recoverer)
    router.Use(middleware.URLFormat)
    
    api := "/api/v1/events"
    
    _, _ = api, handler

    // start server
   
}
