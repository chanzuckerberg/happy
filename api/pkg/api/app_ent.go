package api

import (
	"context"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	_ "github.com/chanzuckerberg/happy/api/pkg/ent/runtime"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/ogen-go/ogen/middleware"
	"github.com/sirupsen/logrus"
)

type handler struct {
	*ogent.OgentHandler
	db *dbutil.DB
}

func GetOgentServer(cfg *setup.Configuration) (*ogent.Server, error) {
	db := dbutil.MakeDB(cfg.Database)
	return ogent.NewServer(
		handler{
			OgentHandler: ogent.NewOgentHandler(db.GetDBEnt()),
			db:           db,
		},
		ogent.WithMiddleware(func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
			log := logrus.New()
			log.Level = resolveLogLevel(cfg.Api.LogLevel)
			start := time.Now()
			req.Context = context.WithValue(req.Context, "log", log)
			defer func() {
				log.Infof("[%s]\t%s\t%s\t%s\t%s", time.Now().UTC().Format(time.RFC3339), req.Raw.Method, req.Raw.URL.Path, req.Raw.URL.RawQuery, time.Since(start))
			}()
			return next(req)
		}),
	)
}

func resolveLogLevel(logLevel string) logrus.Level {
	switch logLevel {
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "silent":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
