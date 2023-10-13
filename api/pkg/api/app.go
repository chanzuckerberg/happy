package api

import (
	"context"
	"fmt"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
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
			ReadTimeout:    60 * time.Second,
			ReadBufferSize: 1024 * 64,
		}),
		Cfg: cfg,
	}
}

func MakeApp(ctx context.Context, cfg *setup.Configuration) *APIApplication {
	db := dbutil.MakeDB(cfg.Database)
	return MakeAppWithDB(ctx, cfg, db)
}

func MakeAppWithDB(ctx context.Context, cfg *setup.Configuration, db *dbutil.DB) *APIApplication {
	apiApp := MakeAPIApplication(cfg).WithDatabase(db)
	apiApp.FiberApp.Use(requestid.New())
	apiApp.FiberApp.Use(cors.New(cors.Config{
		AllowHeaders: "Authorization,Content-Type,x-aws-access-key-id,x-aws-secret-access-key,x-aws-session-token,baggage,sentry-trace",
	}))
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
		verifiers := []request.OIDCVerifier{
			request.MakeGithubVerifier("chanzuckerberg"),
		}
		for _, provider := range cfg.Auth.Providers {
			verifier, err := request.MakeOIDCProvider(ctx, provider.IssuerURL, provider.ClientID, request.DefaultClaimsVerifier)
			if err != nil {
				logrus.Fatalf("failed to create OIDC verifier with error: %s", err.Error())
			}
			verifiers = append(verifiers, verifier)
		}

		v1.Use(request.MakeAuth(request.MakeMultiOIDCVerifier(verifiers...)))
	}

	v1.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	v1.Use(func(c *fiber.Ctx) error {
		user := sentry.User{}
		if email := c.Locals(request.OIDCClaimsEmail{}); email != nil {
			user.Email = email.(string)
		}
		if actor := c.Locals(request.OIDCClaimsGHActor{}); actor != nil {
			user.Username = actor.(string)
		}
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(user)
		})

		txn := sentry.StartSpan(c.Context(), c.Method(), sentry.WithTransactionName(c.Path()))
		defer txn.Finish()
		res := c.Next()
		return res
	})

	RegisterConfigV1(v1, MakeConfigHandler(cmd.MakeConfig(apiApp.DB)))
	RegisterStackListV1(v1, MakeStackHandler(cmd.MakeStack(apiApp.DB)))

	return apiApp
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
