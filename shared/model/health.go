package model

// @Description Health info
type HealthResponse struct {
	Status string `json:"status" example:"OK"`
	Route  string `json:"route" example:"/"`
} // @Name response.Health
