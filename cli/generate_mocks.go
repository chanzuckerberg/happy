package main

//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_compute_backend.go -package=interfaces github.com/chanzuckerberg/happy/cli/pkg/backend/aws/interfaces ComputeBackend
