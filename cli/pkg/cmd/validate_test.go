package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStackNameIsInDnsCharset(t *testing.T) {
	type TestCase struct {
		stackName   string
		expectError bool
	}

	testCases := []TestCase{
		{
			stackName:   "foobar",
			expectError: false,
		},
		{
			stackName:   "f00bar",
			expectError: false,
		},
		{
			stackName:   "00fbar",
			expectError: false,
		},
		{
			stackName:   "0w0",
			expectError: false,
		},
		{
			stackName:   "ew0",
			expectError: false,
		},
		{
			stackName:   "f00b@r",
			expectError: true,
		},
		{
			stackName:   "1234",
			expectError: true,
		},
		{
			stackName:   "-foobar",
			expectError: true,
		},
		{
			stackName:   "foobar-",
			expectError: true,
		},
		{
			stackName:   "-",
			expectError: true,
		},
		{
			stackName:   "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar",
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.stackName, func(t *testing.T) {
			r := require.New(t)
			err := IsStackNameDNSCharset(nil, []string{testCase.stackName})
			if testCase.expectError {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
