package logger

import (
	"log/slog"
	"os"
	"github.com/samber/slog-multi"
)

func SetupLogger(env string, logFile string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = setupPrettySlog(logFile)
	case "dev":
		log = setupFileAndStdoutLogger(logFile, slog.LevelDebug)
	case "prod":
		log = setupFileAndStdoutLogger(logFile, slog.LevelInfo)
	}
	return log
}

func setupPrettySlog(logFileName string) *slog.Logger {
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(
		slogmulti.Fanout(
			stdoutHandler,  
			fileHandler,    
		),
	)
}

func setupFileAndStdoutLogger(logFileName string, level slog.Level) *slog.Logger {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(
		slogmulti.Fanout(
			stdoutHandler, 
			fileHandler,   
		),
	)
}

