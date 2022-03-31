package artifact_builder

import (
	"sort"
	"strconv"
	"testing"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	r := require.New(t)

	bootstrap := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
	}
	happyConfig, err := config.NewHappyConfig(bootstrap)
	r.NoError(err)

	builderConfig := NewBuilderConfig().WithBootstrap(bootstrap).WithHappyConfig(happyConfig).WithExecutor(util.NewDummyExecutor())
	r.NotNil(builderConfig)

	containers, err := builderConfig.GetContainers()
	r.NoError(err)

	expectContainers := []string{
		"database.genepinet.localdev",
		"frontend.genepinet.localdev",
		"localstack.genepinet.localdev",
		"oidc.genepinet.localdev",
		"pangolin.genepinet.localdev",
		"gisaid.genepinet.localdev",
		"nextstrain.genepinet.localdev",
		"backend.genepinet.localdev",
	}

	sort.Strings(containers)
	sort.Strings(expectContainers)

	r.Equal(expectContainers, containers)
}

func TestNewBuilderConfigProfiles(t *testing.T) {
	r := require.New(t)

	bootstrap := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
	}

	happyConfig, err := config.NewHappyConfig(bootstrap)
	r.NoError(err)

	type testcase struct {
		profile          config.Profile
		expectContainers []string
	}

	testCases := []testcase{
		{
			profile: config.Profile("backend"),
			expectContainers: []string{
				"database.genepinet.localdev",
				"localstack.genepinet.localdev",
				"oidc.genepinet.localdev",
				"backend.genepinet.localdev",
			},
		},
		{ // TODO: should we fail if profile not found?
			profile:          config.Profile("frontend"),
			expectContainers: nil,
		},
		{
			profile: config.Profile("jobs"),
			expectContainers: []string{
				"pangolin.genepinet.localdev",
				"gisaid.genepinet.localdev",
				"nextstrain.genepinet.localdev",
			},
		},
	}

	for idx, testCase := range testCases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			r := require.New(t)
			bc := NewBuilderConfig().WithBootstrap(bootstrap).WithHappyConfig(happyConfig).WithProfile(&testCase.profile)
			r.NotNil(bc)
			containers, err := bc.GetContainers()
			r.NoError(err)

			sort.Strings(containers)
			sort.Strings(testCase.expectContainers)

			r.Equal(testCase.expectContainers, containers)
		})
	}
}
