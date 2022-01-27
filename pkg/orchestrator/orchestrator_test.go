package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOrchestrator(t *testing.T) {
	r := require.New(t)
	r.NotNil(NewOrchestrator(nil, nil))
}
