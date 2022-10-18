package setup

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvToMap(t *testing.T) {
	r := require.New(t)
	os.Clearenv()
	err := os.Setenv("ENV1", "test1")
	r.NoError(err)
	err = os.Setenv("ENV2", "test2")
	r.NoError(err)
	defer os.Unsetenv("ENV1")
	defer os.Unsetenv("ENV2")
	m := envToMap()
	r.Contains(m, "ENV1")
	r.Contains(m, "ENV2")
}

func TesEvaluateConfigWithEnv(t *testing.T) {
	r := require.New(t)
	test := `blah={{.ENV1}}
blah2={{.ENV2}}`
	os.Setenv("ENV1", "test1")
	os.Setenv("ENV2", "test2")
	defer os.Unsetenv("ENV1")
	defer os.Unsetenv("ENV2")
	eval, err := evaluateConfigWithEnv(strings.NewReader(test))
	r.NoError(err)
	expected := `blah=test1
blah2=test2`
	b, err := io.ReadAll(eval)
	r.NoError(err)
	r.Equal(expected, string(b))
}
