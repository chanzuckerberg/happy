package config

import "testing"

// generates a test happy config
// only use in tests
func NewTestHappyConfig(
	t *testing.T,
	testFilePath string,
	env string,
) (HappyConfig, error) {
	b := &Bootstrap{
		Env:             env,
		HappyConfigPath: testFilePath,
	}
	return NewHappyConfig(b)
}
