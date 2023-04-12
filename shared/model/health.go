package model

// @Description Health info
type HealthResponse struct {
	Status  string `json:"status" example:"OK"`
	Route   string `json:"route" example:"/"`
	Version string `json:"version" example:"v0.50.2"`
	GitSha  string `json:"git_sha" example:"b746df8"`
} // @Name response.Health
