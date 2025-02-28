package main

import (
	"fmt"
	"log/slog"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/config"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/logger"
)


func main(){
    cfg := config.MustLoadConfig()
    log := logger.SetupLogger(cfg.Env)

    log.Info("Connecting to db with params")
    log.Info("Database: ", slog.String("host", cfg.Database.Host), slog.String("port", cfg.Database.Port))

    dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",cfg.Database.User, cfg.Database.Password, "authdb", cfg.Database.Host, cfg.Database.Port, "disable")
    dbConnection := db.SetupDB(dsn)

 
}
