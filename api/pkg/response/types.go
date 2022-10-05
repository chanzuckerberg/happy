package response

// @Description Health info
type HealthResponse struct {
	Status string `json:"status" example:"OK"`
	Route  string `json:"route" example:"/"`
} // @Name response.Health

type ValidationError struct {
	FailedField string `json:"failed_field"` // the field that failed to be validated
	Tag         string `json:"tag" swaggerignore:"true"`
	Value       string `json:"value" swaggerignore:"true"`
	Type        string `json:"type" swaggerignore:"true"`
	Message     string `json:"message"` // a description of the error that occured
}

// @Description Object denoting which app config keys are missing from the destination env/stack
type ConfigDiffResponse struct {
	MissingKeys []string `json:"missing_keys" example:"SOME_KEY,ANOTHER_KEY"`
} // @Name response.ConfigDiff
