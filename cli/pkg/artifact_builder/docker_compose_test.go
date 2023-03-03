package artifact_builder

import (
	"strconv"
	"testing"

	"github.com/chanzuckerberg/happy/cli/pkg/config"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/stretchr/testify/require"
)

func TestInvokeDockerComposeConfig(t *testing.T) {
	type testcase struct {
		profile           *config.Profile
		expectedServices  []string
		expectedPlatforms map[string]string
	}

	profileBackend := config.Profile("backend")
	tcases := []testcase{
		{
			profile:          nil,
			expectedServices: []string{"database", "frontend", "backend", "localstack", "oidc", "gisaid", "pangolin", "nextstrain"},
			expectedPlatforms: map[string]string{
				"frontend": "linux/amd64",
			},
		},
		{
			profile:           &profileBackend,
			expectedServices:  []string{"database", "backend", "localstack", "oidc"},
			expectedPlatforms: map[string]string{},
		},
	}

	for idx, tcase := range tcases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			r := require.New(t)
			bc := &BuilderConfig{
				composeFile: testDockerComposePath,
				Profile:     tcase.profile,
			}

			bc.WithExecutor(util.NewDefaultExecutor())

			data, err := bc.DockerComposeConfig()
			r.NoError(err)

			for _, k := range tcase.expectedServices {
				r.Contains(data.Services, k)
			}
			for k, service := range data.Services {
				r.Contains(tcase.expectedServices, k)
				platform, ok := tcase.expectedPlatforms[k]
				if ok {
					r.Equal(platform, service.Platform)
				}
			}
		})
	}
}
