package api

import (
	"context"
	"os"
	"time"

	"log/slog"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	_ "github.com/chanzuckerberg/happy/api/pkg/ent/runtime"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/ogen-go/ogen/middleware"
)

type handler struct {
	*ogent.OgentHandler
	db *store.DB
}

func (h handler) Health(_ context.Context) (ogent.HealthRes, error) {
	return &ogent.HealthOK{Status: "ok"}, nil
}

func (h handler) ListAppConfig(ctx context.Context, params ogent.ListAppConfigParams) (ogent.ListAppConfigRes, error) {
	res, err := h.db.GetAppConfigsForStack(ctx, params.AppName, params.Environment, params.Stack.Or(""))
	if err != nil {
		return nil, err
	}

	r := ogent.NewAppConfigLists(res)
	return (*ogent.ListAppConfigOKApplicationJSON)(&r), nil
}

func GetOgentServer(cfg *setup.Configuration) (*ogent.Server, error) {
	db := store.MakeDB(cfg.Database)
	return ogent.NewServer(
		handler{
			OgentHandler: ogent.NewOgentHandler(db.GetDB()),
			db:           db,
		},
		ogent.WithMiddleware(func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: resolveLogLevel(cfg.Api.LogLevel),
			}))
			start := time.Now()
			req.Context = context.WithValue(req.Context, "log", log)

			res, err := next(req)
			// log.Info(fmt.Sprint(res.Type))
			var status int
			if tresp, ok := res.Type.(interface{ GetStatusCode() int }); ok {
				log.Info("here")
				status = tresp.GetStatusCode()
			}
			log.Info("Request complete", "status", status, "duration", time.Since(start), "method", req.Raw.Method, "path", req.Raw.URL.Path, "query", req.Raw.URL.RawQuery)

			return res, err
		}),
	)
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
