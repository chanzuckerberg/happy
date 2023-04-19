package artifact_builder

import (
	"context"
	"sort"
	"strconv"
	"testing"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/stretchr/testify/require"
)

func TestNewBuilderConfig(t *testing.T) {
	ctx := diagnostics.BuildDiagnosticContext(context.Background(), true)
	r := require.New(t)

	bootstrap := &config.Bootstrap{
		HappyConfigPath:         testFilePath,
		DockerComposeConfigPath: testDockerComposePath,
	}
	happyConfig, err := config.NewHappyConfig(bootstrap)
	r.NoError(err)

	builderConfig := NewBuilderConfig().WithBootstrap(bootstrap).WithHappyConfig(happyConfig)
	r.NotNil(builderConfig)

	containers, err := builderConfig.GetContainers(ctx)
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
	ctx := diagnostics.BuildDiagnosticContext(context.Background(), true)
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
			bc := NewBuilderConfig().WithBootstrap(bootstrap).WithHappyConfig(happyConfig)
			bc.Profile = &testCase.profile
			r.NotNil(bc)
			containers, err := bc.GetContainers(ctx)
			r.NoError(err)

			sort.Strings(containers)
			sort.Strings(testCase.expectContainers)

			r.Equal(testCase.expectContainers, containers)
		})
	}
}
