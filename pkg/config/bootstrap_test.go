package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func cleanup() {
	happyProjectRoot = ""
	happyConfigPath = ""
	dockerComposeConfigPath = ""
	env = ""
}

func setEnvs(t *testing.T, setenv map[string]string) {
	for key, val := range setenv {
		t.Setenv(key, val)
	}
}

func setFlags(setflags map[string]string) {
	if val, ok := setflags[flagHappyProjectRoot]; ok {
		happyProjectRoot = val
	}
	if val, ok := setflags[flagHappyConfigPath]; ok {
		happyConfigPath = val
	}
	if val, ok := setflags[flagDockerComposeConfigPath]; ok {
		dockerComposeConfigPath = val
	}
	if val, ok := setflags[flagEnv]; ok {
		env = val
	}
}

func TestNewBootstrapConfig(t *testing.T) {
	defer cleanup()

	testCases := []struct {
		name     string
		setenvs  map[string]string
		setflags map[string]string

		wantError  bool
		wantConfig *Bootstrap
	}{
		{
			name:      "error if fields missing",
			wantError: true,
		},
		{
			name: "just envs, no flags",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
			},
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "rdev",
			},
		},
		{
			name: "just flags, no envs",
			setflags: map[string]string{
				flagHappyConfigPath:         "foo",
				flagHappyProjectRoot:        ".",
				flagDockerComposeConfigPath: "bar",
			},
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "rdev",
			},
		},
		{
			name: "flags override some envs",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
			},
			setflags: map[string]string{
				flagHappyConfigPath: "flagfoo",
			},
			wantConfig: &Bootstrap{
				HappyConfigPath:         "flagfoo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "rdev",
			},
		},
		{
			name: "override env",
			setflags: map[string]string{
				flagHappyConfigPath:         "foo",
				flagHappyProjectRoot:        ".",
				flagDockerComposeConfigPath: "bar",
				flagEnv:                     "flagenv",
			},
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "flagenv",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			defer cleanup()

			setEnvs(t, tc.setenvs)
			setFlags(tc.setflags)

			bc, err := NewBootstrapConfig()
			if tc.wantError {
				r.Error(err)
				return
			}
			r.NoError(err)

			r.Equal(tc.wantConfig, bc)
		})
	}
}

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
