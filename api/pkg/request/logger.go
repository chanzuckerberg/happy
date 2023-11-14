package request

import (
	"context"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/ogen-go/ogen/middleware"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type LoggerKey struct{}

func MakeOgentLoggerMiddleware(cfg *setup.Configuration) ogent.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		logger, err := newLogger(cfg.Api.LogLevel)
		defer func() {
			err := logger.Sync()
			if err != nil {
				logrus.Fatal(err)
			}
		}()

		start := time.Now()
		req.Context = context.WithValue(req.Context, LoggerKey{}, logger)

		res, err := next(req)

		var status int = 200
		if tresp, ok := res.Type.(interface{ GetStatusCode() int }); ok {
			status = tresp.GetStatusCode()
		} else if tresp, ok := res.Type.(interface{ GetCode() int }); ok {
			status = tresp.GetCode()
		}

		if err == nil {
			logger.Info("Success", getLogArgs(status, start, req)...)
		} else {
			if terr, ok := err.(interface{ GetCode() int }); ok {
				status = terr.GetCode()
			}
			args := append([]zap.Field{zap.String("error", err.Error())}, getLogArgs(status, start, req)...)
			logger.Error("Fail", args...)
		}

		return res, err
	}
}

func getLogArgs(status int, start time.Time, req middleware.Request) []zap.Field {
	return []zap.Field{
		zap.Int("status", status),
		zap.Duration("duration", time.Since(start)),
		zap.String("method", req.Raw.Method),
		zap.String("path", req.Raw.URL.Path),
		zap.String("query", req.Raw.URL.RawQuery),
	}
}

func newLogger(logLevel string) (*zap.Logger, error) {
	switch logLevel {
	case "debug":
		return zap.NewDevelopment()
	case "error":
		return zap.NewProduction()
	case "warn":
		return zap.NewProduction()
	case "silent":
		return zap.NewNop(), nil
	default:
		return zap.NewProduction()
	}
}
