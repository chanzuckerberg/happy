package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type APIApplication struct {
	cfg *setup.Configuration
	mux *http.ServeMux
	DB  *store.DB
}

func MakeAPIApplication(ctx context.Context, cfg *setup.Configuration, db *store.DB) *APIApplication {
	// create a mux to route requests to the correct app
	rootMux := http.NewServeMux()
	rootMux.Handle("/", request.HealthHandler{})
	rootMux.Handle("/health", request.HealthHandler{})
	rootMux.Handle("/versionCheck", request.VersionCheckHandler{})

	// create the Fiber app
	app := MakeFiberApp(ctx, cfg, db)
	nativeHandler := adaptor.FiberApp(app.FiberApp)
	rootMux.Handle("/v1/", http.StripPrefix("/v1", nativeHandler))

	// create the Ogent app
	svr, err := MakeOgentServer(ctx, cfg, db)
	if err != nil {
		logrus.Fatal(err)
	}
	rootMux.Handle("/v2/", svr)

	return &APIApplication{
		cfg: cfg,
		DB:  db,
		mux: rootMux,
	}
}

func (a *APIApplication) Listen() error {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Content-Type", "x-aws-access-key-id", "x-aws-secret-access-key", "x-aws-session-token", "baggage", "sentry-trace"},
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Api.Port), c.Handler(a.mux))
}
