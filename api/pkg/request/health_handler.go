package request

import (
	"encoding/json"
	"net/http"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary Check that the server is up.
// @Tags    root
// @Accept  */*
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router  /health [get]
func HealthHandlerFiber(c *fiber.Ctx) error {
	status := HealthOKResponse(c.Path())
	return c.JSON(status)
}

type HealthHandler struct{}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := HealthOKResponse(r.URL.Path)
	b, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to convert HealthResponse to json"))
	}
	w.Write(b)
}

func HealthOKResponse(path string) model.HealthResponse {
	return model.HealthResponse{
		Status:  "OK",
		Route:   path,
		Version: util.ReleaseVersion,
		GitSha:  util.ReleaseGitSha,
	}
}
