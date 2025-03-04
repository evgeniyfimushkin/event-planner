package main

import (
	"registration-service/internal/handler"
	"registration-service/internal/models"
	"registration-service/internal/repository"
	"registration-service/internal/service"
	"fmt"
	"log/slog"
	"net/http"

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
    dbConnection := db.SetupDB(dsn, &models.Registration{})
    registrationRepo := repository.NewRegistrationRepository(dbConnection)
    verifier, err := auth.NewVerifier(cfg.PublicKey)
    if err != nil {
        log.Error("failed to init JWT verifier", logger.Err(err))
        panic("failed to init JWT verifier")
    }

    registrationservice := service.NewRegistrationService(registrationRepo)

    handler := handler.NewRegistrationHandler(registrationservice, verifier)

    router := chi.NewRouter()
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middlewarelogger.New(log))
    router.Use(middleware.Recoverer)
    router.Use(middleware.URLFormat)

    router.Post("/api/v1/registrations", handler.CreateHandler())
    router.Get("/api/v1/registrations", handler.GetAllHandler())
    router.Get("/api/v1/registrations/{id}", handler.GetByIDHandler())
    router.Put("/api/v1/registrations", handler.UpdateHandler())
    router.Delete("/api/v1/registrations/{id}", handler.DeleteHandler())
    router.Delete("/api/v1/registrations", handler.DeleteWhereHandler())
    router.Get("/api/v1/registrations/search", handler.FindHandler())
    router.Get("/api/v1/registrations/search/first", handler.FindFirstHandler())
    router.Get("/api/v1/registrations/count", handler.CountHandler())
    router.Get("/api/v1/registrations/page", handler.GetPageHandler())
    router.Post("/api/v1/registrations/bulk", handler.BulkInsertHandler())
    router.Put("/api/v1/registrations/bulk", handler.BulkUpdateHandler())



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
