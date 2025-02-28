package logger

import "log/slog"

// Err gets err and returns this wrapped by slog.Attr
func Err(err error) slog.Attr{
    return slog.Attr{
        Key: "error",
        Value: slog.StringValue(err.Error()),
    }
}
