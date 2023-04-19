package diagnostics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTfeRunInfo(t *testing.T) {
	type TestCase struct {
		url         string
		org         string
		workspace   string
		runId       string
		expectedUrl string
		expectMatch bool
		expectError bool
	}

	testCases := []TestCase{
		{
			url:         "https://example.com",
			org:         "happy-ig",
			workspace:   "rdev-dtsai",
			runId:       "run-deadboof",
			expectedUrl: "https://example.com/app/happy-ie/workspaces/rdev-dtsai/runs/run-deadbeef",
			expectMatch: false,
			expectError: false,
		},
		{
			url:         "https://example.com",
			org:         "happy-ie",
			workspace:   "rdev-dtsai",
			runId:       "run-deadbeef",
			expectedUrl: "https://example.com/app/happy-ie/workspaces/rdev-dtsai/runs/run-deadbeef",
			expectMatch: true,
			expectError: false,
		},
		{
			url:         "https://example.com",
			org:         "",
			workspace:   "rdev-dtsai",
			runId:       "run-deadbeet",
			expectedUrl: "",
			expectMatch: true,
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.runId, func(t *testing.T) {
			r := require.New(t)
			info := NewTfeRunInfo()
			info.AddTfeUrl(testCase.url)
			info.AddOrg(testCase.org)
			info.AddWorkspace(testCase.workspace)
			info.AddRunId(testCase.runId)
			link, err := info.MakeTfeRunLink()
			if testCase.expectError {
				r.Error(err)
			} else {
				r.NoError(err)
			}
			match := testCase.expectedUrl == link
			r.Equal(testCase.expectMatch, match)
		})
	}
}
