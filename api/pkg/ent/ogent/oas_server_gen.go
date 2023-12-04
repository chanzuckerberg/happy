// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// Health implements Health operation.
	//
	// Simple endpoint to check if the server is up.
	//
	// GET /health
	Health(ctx context.Context) (HealthRes, error)
	// ListAppConfig implements listAppConfig operation.
	//
	// GET /app-configs
	ListAppConfig(ctx context.Context, params ListAppConfigParams) (ListAppConfigRes, error)
	// ReadAppConfig implements readAppConfig operation.
	//
	// Finds the AppConfig with the requested Key and returns it.
	//
	// GET /app-configs/{key}
	ReadAppConfig(ctx context.Context, params ReadAppConfigParams) (ReadAppConfigRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}