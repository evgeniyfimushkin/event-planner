package logger

import (
	"log/slog"
	"os"
)

// SetupLogger get env var and configure returns configured loggers that depends on env var
func SetupLogger(env string) *slog.Logger {
    var log *slog.Logger
    switch env {
    case "local":
        log = setupPrettySlog()
    case "dev":
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case "prod":
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }
    return log
}

func setupPrettySlog() *slog.Logger {
    opts := PrettyHandlerOptions{
        SlogOpts: &slog.HandlerOptions{
            Level: slog.LevelDebug,
        },
    }
    handler := opts.NewPrettyHandler(os.Stdout)
    return slog.New(handler)
}
