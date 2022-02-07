package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func cleanup() {
	happyProjectRoot = ""
	happyConfigPath = ""
	dockerComposeConfigPath = ""
	env = ""
}

func setEnvs(t *testing.T, basedir string, setenv map[string]string) {
	set := setBaseDir(basedir)
	for key, val := range setenv {
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
			createFiles: true,
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
			createFiles: true,
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
				Env:                     "rdev",
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

func TestSearchHappyRoot(t *testing.T) {
	r := require.New(t)

	tmpDir, err := os.MkdirTemp("", "")
	r.NoError(err)
	defer os.RemoveAll(tmpDir)

	nested := filepath.Join(tmpDir, "/a/b/c/d/e/f/g/h")
	err = os.MkdirAll(nested, 0777)
	r.NoError(err)

	happyDir := filepath.Join(tmpDir, "/a/b/c/d/e/.happy")
	err = os.MkdirAll(happyDir, 0777)
	r.NoError(err)

	// if I search from nested, I should find .happy path
	found, err := searchHappyRoot(nested)
	r.NoError(err)
	r.Equal(happyDir, found)

	// if I search from happyDir, I should find
	found, err = searchHappyRoot(happyDir)
	r.NoError(err)
	r.Equal(happyDir, found)

	// if I search from outside the tree, I should not find
	outside := filepath.Join(tmpDir, "/a/b/c/")
	found, err = searchHappyRoot(outside)
	r.EqualError(err, errCouldNotInferFindHappyRoot.Error())
	r.Empty(found)
}

// generates a test happy config
// only use in tests
func NewTestHappyConfig(
	t *testing.T,
	testFilePath string,
	env string,
	awsSecretMgr SecretsBackend,
) (HappyConfig, error) {
	b := &Bootstrap{
		Env:             env,
		HappyConfigPath: testFilePath,
	}
	return NewHappyConfigWithSecretsBackend(b, awsSecretMgr)
}
