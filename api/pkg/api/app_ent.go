package api

import (
	"context"
	"net/http"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	_ "github.com/chanzuckerberg/happy/api/pkg/ent/runtime"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/chanzuckerberg/happy/shared/util"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/validate"
	"github.com/pkg/errors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type handler struct {
	*ogent.OgentHandler
	db *store.DB
}

func (h handler) Health(ctx context.Context) (ogent.HealthRes, error) {
	return &ogent.HealthOK{Status: "OK", Version: util.ReleaseVersion, GitSha: util.ReleaseGitSha, Route: "/v2/health"}, nil
}

func MakeOgentServer(ctx context.Context, cfg *setup.Configuration, db *store.DB) (*ogent.Server, error) {
	middlewares := []ogent.Middleware{request.MakeOgentLoggerMiddleware(cfg)}
	if *cfg.Auth.Enable {
		verifier := request.MakeVerifierFromConfig(ctx, cfg)
		middlewares = append(middlewares, request.MakeOgentAuthMiddleware(verifier))
	}

	serverOpts := []ogent.ServerOption{
		ogent.WithPathPrefix("/v2"),
		ogent.WithMiddleware(middlewares...),
		ogent.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
			code := 500
			var customErr response.CustomError
			if errors.As(err, &customErr) {
				code = customErr.GetCode()
			} else if errors.Is(err, validate.ErrFieldRequired) {
				code = 400
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)

			e := jx.GetEncoder()
			e.ObjStart()
			e.FieldStart("code")
			e.Int(code)
			e.FieldStart("errors")
			e.StrEscape(err.Error())
			e.ObjEnd()

			_, _ = w.Write(e.Bytes())
		}),
	}

	if cfg.Sentry.DSN != "" {
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
		)
		serverOpts = append(serverOpts, ogent.WithTracerProvider(tp))
	}

	return ogent.NewServer(
		handler{
			OgentHandler: ogent.NewOgentHandler(db.GetDB()),
			db:           db,
		},
		serverOpts...,
	)
}
