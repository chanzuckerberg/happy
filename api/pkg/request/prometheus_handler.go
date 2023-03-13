package request

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var fasthttpPrometheusHandler = fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())

func PrometheusMetricsHandler(c *fiber.Ctx) error {
	fasthttpPrometheusHandler(c.Context())
	return nil
}
