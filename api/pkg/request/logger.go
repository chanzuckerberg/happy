package request

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/ogen-go/ogen/middleware"
)

type LoggerKey struct{}

func MakeOgentLoggerMiddleware(cfg *setup.Configuration) ogent.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: resolveLogLevel(cfg.Api.LogLevel),
		}))
		start := time.Now()
		req.Context = context.WithValue(req.Context, LoggerKey{}, log)

		res, err := next(req)
		var status int
		if tresp, ok := res.Type.(interface{ GetStatusCode() int }); ok {
			log.Info("here")
			status = tresp.GetStatusCode()
		}
		log.Info("Request complete", "status", status, "duration", time.Since(start), "method", req.Raw.Method, "path", req.Raw.URL.Path, "query", req.Raw.URL.RawQuery)

		return res, err
	}
}

func resolveLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "silent":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
