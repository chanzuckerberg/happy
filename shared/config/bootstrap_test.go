package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func cleanup() {
	happyProjectRoot = ""
	happyConfigPath = ""
	dockerComposeConfigPath = ""
	env = ""
	awsProfile = ""
}

func setEnvs(t *testing.T, basedir string, setenv map[string]string) {
	set := setBaseDir(basedir)
	for key, val := range setenv {
		if key == "AWS_PROFILE" || key == "HAPPY_ENV" {
			t.Setenv(key, val)
			continue
		}
		t.Setenv(key, set(val))
	}
}

func setBaseDir(base string) func(string) string {
	return func(other string) string {
		return filepath.Join(base, other)
	}
}

func setFlags(basedir string, setflags map[string]string) {
	set := setBaseDir(basedir)

	if val, ok := setflags[flagHappyProjectRoot]; ok {
		happyProjectRoot = set(val)
	}
	if val, ok := setflags[flagHappyConfigPath]; ok {
		happyConfigPath = set(val)
	}
	if val, ok := setflags[flagDockerComposeConfigPath]; ok {
		dockerComposeConfigPath = set(val)
	}
	if val, ok := setflags[FlagAWSProfile]; ok {
		awsProfile = val
	}
	if val, ok := setflags[flagEnv]; ok {
		env = val
	}
}

func createExpectedFiles(r *require.Assertions, basedir string, b *Bootstrap) {
	applyBasedirWantConfig(basedir, b)

	create := func(p string) {
		if p == "" {
			return
		}

		d := filepath.Dir(p)
		r.NoError(os.MkdirAll(d, 0777))

		logrus.Warnf("creating %s", p)
		_, err := os.Create(p)
		r.NoError(err)
	}

	create(b.DockerComposeConfigPath)
	create(b.HappyConfigPath)
}

func applyBasedirWantConfig(basedir string, expected *Bootstrap) {
	if expected == nil {
		return
	}

	set := setBaseDir(basedir)

	expected.DockerComposeConfigPath = set(expected.DockerComposeConfigPath)
	expected.HappyConfigPath = set(expected.HappyConfigPath)
	expected.HappyProjectRoot = set(expected.HappyProjectRoot)
}

func TestNewBootstrapConfig(t *testing.T) {
	defer cleanup()

	testCases := []struct {
		name     string
		setenvs  map[string]string
		setflags map[string]string

		// set this if you want the test framework
		// to touch files and pass existence validation
		createFiles bool

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
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "",
			},
		},
		{
			name: "just flags, no envs",
			setflags: map[string]string{
				flagHappyConfigPath:         "foo",
				flagHappyProjectRoot:        ".",
				flagDockerComposeConfigPath: "bar",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "",
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
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "flagfoo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "",
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
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "flagenv",
			},
		},
		{
			name: "inferred when possible",
			setflags: map[string]string{
				flagHappyProjectRoot: "/a/b/c",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "/a/b/c/.happy/config.json",
				HappyProjectRoot:        "/a/b/c",
				DockerComposeConfigPath: "/a/b/c/docker-compose.yml",
				Env:                     "",
			},
		},
		{
			name: "set aws profile env",
			setenvs: map[string]string{
				"AWS_PROFILE":                "",
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "",
				AWSProfile:              aws.String(""),
			},
		},
		{
			name: "flagEnv overrides HAPPY_ENV",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
				"HAPPY_ENV":                  "happyEnv",
			},
			setflags: map[string]string{
				flagEnv: "flagEnv",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "flagEnv",
			},
		},
		{
			name: "HAPPY_ENV sets env when flagEnv is absent",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
				"HAPPY_ENV":                  "happyEnv",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "happyEnv",
			},
		},
		{
			name: "flagEnv sets env when HAPPY_ENV is absent",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
			},
			setflags: map[string]string{
				flagEnv: "flagEnv",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "flagEnv",
			},
		},
		{
			name: "no env is set when both flagEnv and HAPPY_ENV are absent",
			setenvs: map[string]string{
				"HAPPY_CONFIG_PATH":          "foo",
				"HAPPY_PROJECT_ROOT":         ".",
				"DOCKER_COMPOSE_CONFIG_PATH": "bar",
			},
			createFiles: true,
			wantConfig: &Bootstrap{
				HappyConfigPath:         "foo",
				HappyProjectRoot:        ".",
				DockerComposeConfigPath: "bar",
				Env:                     "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			defer cleanup()
			// create a tmpdir as root for tests
			basedir, err := os.MkdirTemp("", "")
			r.NoError(err)
			defer os.RemoveAll(basedir)

			logrus.Error(tc.wantConfig)

			if tc.createFiles {
				createExpectedFiles(r, basedir, tc.wantConfig)
			}

			setEnvs(t, basedir, tc.setenvs)
			setFlags(basedir, tc.setflags)

			bc, err := NewBootstrapConfig(&cobra.Command{})
			if tc.wantError {
				r.Error(err)
				return
			}
			r.NoError(err)

			r.Equal(tc.wantConfig, bc)
		})
	}
}

func TestSearchHappyRoot(t *testing.T) {
	r := require.New(t)

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	nested := filepath.Join(tmpDir, "/a/b/c/d/e/f/g/h")
	err = os.MkdirAll(nested, 0777)
	r.NoError(err)

	happyRoot := filepath.Join(tmpDir, "/a/b/c/d/e")

	happyDir := filepath.Join(happyRoot, "/.happy")
	err = os.MkdirAll(happyDir, 0777)
	r.NoError(err)

	// if I search from the "root", I should find
	found, err := searchHappyRoot(happyRoot)
	r.NoError(err)
	r.Equal(happyRoot, found)

	// if I search from nested, I should find .happy path
	found, err = searchHappyRoot(nested)
	r.NoError(err)
	r.Equal(happyRoot, found)

	// if I search from happyDir, I should find
	found, err = searchHappyRoot(happyDir)
	r.NoError(err)
	r.Equal(happyRoot, found)

	// if I search from outside the tree, I should not find
	outside := filepath.Join(tmpDir, "/a/b/c/")
	found, err = searchHappyRoot(outside)
	r.EqualError(err, errCouldNotInferFindHappyRoot.Error())
	r.Empty(found)
}

func TestFindFile(t *testing.T) {
	r := require.New(t)
	_, err := findFile("COVERAGE", []string{""})
	r.NoError(err)

	bootstrap := &Bootstrap{
		HappyConfigPath:          "foo",
		HappyProjectRoot:         ".",
		DockerComposeConfigPath:  "bar",
		DockerComposeEnvFilePath: "COVERAGE",
		Env:                      "",
	}

	_, err = findDockerComposeEnvFile(".env.ecr", bootstrap)
	r.NoError(err)
}

func TestBootstrapValidateCustomErrors(t *testing.T) {
	r := require.New(t)
	b := &Bootstrap{}
	err := validate.Struct(b)
	r.Error(err)

	missingFields := []string{
		"HappyConfigPath",
		"HappyProjectRoot",
		"DockerComposeConfigPath",
	}
	var expected error
	for _, mf := range missingFields {
		expected = multierror.Append(expected, errors.Errorf("%s is required but was not set and could not be inferred", mf))
	}

	err = prettyValidationErrors(err)
	r.Error(err)
	r.Equal(expected.Error(), err.Error())
}

// generates a test happy config
// only use in tests
func NewTestHappyConfig(
	t *testing.T,
	testFilePath string,
	env string,
) (*HappyConfig, error) {
	b := &Bootstrap{
		Env:             env,
		HappyConfigPath: testFilePath,
	}
	happyConfig, err := NewHappyConfig(b)
	return happyConfig, err
}
