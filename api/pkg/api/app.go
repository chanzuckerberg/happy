package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

type APIApplication struct {
	FiberApp *fiber.App
	DB       *dbutil.DB
	Cfg      *setup.Configuration
}

func MakeAPIApplication(cfg *setup.Configuration) *APIApplication {
	return &APIApplication{
		FiberApp: fiber.New(fiber.Config{
			AppName:        "happy-api",
			JSONEncoder:    MarshalJSON,
			ReadTimeout:    60 * time.Second,
			ReadBufferSize: 1024 * 64,
		}),
		Cfg: cfg,
	}
}

func MakeApp(cfg *setup.Configuration) *APIApplication {
	db := dbutil.MakeDB(cfg.Database)
	apiApp := MakeAPIApplication(cfg).WithDatabase(db)
	apiApp.FiberApp.Use(requestid.New())
	apiApp.FiberApp.Use(cors.New())
	apiApp.configureLogger(cfg.Api)
	apiApp.FiberApp.Use(func(c *fiber.Ctx) error {
		err := request.VersionCheckHandler(c)
		if err != nil {
			return err
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		return c.Next()
	})

	apiApp.FiberApp.Get("/", request.HealthHandler)
	apiApp.FiberApp.Get("/health", request.HealthHandler)
	apiApp.FiberApp.Get("/versionCheck", request.VersionCheckHandler)
	apiApp.FiberApp.Get("/swagger/*", swagger.HandlerDefault)

	apiApp.FiberApp.Get("/metrics", request.PrometheusMetricsHandler)

	v1 := apiApp.FiberApp.Group("/v1")

	if *cfg.Auth.Enable {
		verifiers := []request.OIDCVerifier{}
		for _, provider := range cfg.Auth.Providers {
			verifier, err := request.MakeOIDCProvider(context.Background(), provider.IssuerURL, provider.ClientID, request.DefaultClaimsVerifier)
			if err != nil {
				logrus.Fatalf("failed to create OIDC verifier with error: %s", err.Error())
			}
			verifiers = append(verifiers, verifier)
		}

		v1.Use(request.MakeAuth(request.MakeMultiOIDCVerifier(verifiers...)))
	}

	RegisterConfigV1(v1, MakeConfigHandler(cmd.MakeConfig(apiApp.DB)))
	RegisterStackListV1(v1, MakeStackHandler(cmd.MakeStack(apiApp.DB)))

	return apiApp
}

// Copied from https://gist.github.com/Rican7/39a3dc10c1499384ca91
// with a slight tweak to make "ID" convert to "id" instead of "i_d"
func MarshalJSON(val interface{}) ([]byte, error) {
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w{2,})([A-Z])`)
	marshalled, err := json.Marshal(val)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}

func (a *APIApplication) WithDatabase(db *dbutil.DB) *APIApplication {
	a.DB = db
	err := a.DB.AutoMigrate()
	if err != nil {
		logrus.Fatalf("failed to connect to the DB %s", err)
	}
	return a
}

func (a *APIApplication) configureLogger(cfg setup.ApiConfiguration) {
	if cfg.LogLevel == "silent" {
		return
	}

	a.FiberApp.Use(logger.New(logger.Config{
		Format:     "[${date} ${time}] | ${status} | ${latency} | ${method} | ${path} | ${locals:requestid}\n",
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "UTC",
	}))
}

func (a *APIApplication) Listen() error {
	return a.FiberApp.Listen(fmt.Sprintf(":%d", a.Cfg.Api.Port))
}
