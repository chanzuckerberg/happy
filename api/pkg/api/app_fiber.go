package api

import (
	"context"
	"time"

	"github.com/chanzuckerberg/happy/api/pkg/cmd"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

type FiberServer struct {
	FiberApp *fiber.App
	DB       *store.DB
	Cfg      *setup.Configuration
}

func MakeFiberServer(cfg *setup.Configuration) *FiberServer {
	return &FiberServer{
		FiberApp: fiber.New(fiber.Config{
			AppName:        "happy-api",
			ReadTimeout:    60 * time.Second,
			ReadBufferSize: 1024 * 64,
		}),
		Cfg: cfg,
	}
}

func MakeFiberApp(ctx context.Context, cfg *setup.Configuration, db *store.DB) *FiberServer {
	apiApp := MakeFiberServer(cfg).WithDatabase(db)
	apiApp.FiberApp.Use(requestid.New())
	apiApp.configureLogger(cfg.Api)
	apiApp.FiberApp.Use(func(c *fiber.Ctx) error {
		err := request.VersionCheckHandlerFiber(c)
		if err != nil {
			return err
		}
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		return c.Next()
	})

	v1 := apiApp.FiberApp.Group("/v1")
	v1.Get("/health", request.HealthHandlerFiber)
	v1.Get("/swagger/*", swagger.HandlerDefault)
	v1.Get("/metrics", request.PrometheusMetricsHandler)

	if *cfg.Auth.Enable {
		verifier := request.MakeVerifierFromConfig(ctx, cfg)
		v1.Use(request.MakeFiberAuthMiddleware(verifier))
	}

	v1.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	v1.Use(func(c *fiber.Ctx) error {
		user := sentry.User{}
		oidcValues := c.Locals(request.OIDCAuthKey{})
		if oidcValues != nil {
			oidcValues := oidcValues.(*request.OIDCAuthValues)
			if len(oidcValues.Email) > 0 {
				user.Email = oidcValues.Email
			}
			if len(oidcValues.Actor) > 0 {
				user.Username = oidcValues.Actor
			}
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

func (a *FiberServer) WithDatabase(db *store.DB) *FiberServer {
	a.DB = db
	err := a.DB.AutoMigrate()
	if err != nil {
		logrus.Fatalf("failed to connect to the DB %s", err)
	}
	return a
}

func (a *FiberServer) configureLogger(cfg setup.ApiConfiguration) {
	if cfg.LogLevel == "silent" {
		return
	}

	a.FiberApp.Use(logger.New(logger.Config{
		Format:     "[${date} ${time}] | ${status} | ${latency} | ${method} | ${path} | ${locals:requestid}\n",
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "UTC",
	}))
}
