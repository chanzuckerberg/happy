package cmd

import (
	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/setup"
	"github.com/stretchr/testify/require"
)

func MakeTestDB(r *require.Assertions) *dbutil.DB {
	config, err := setup.GetConfiguration()
	r.NoError(err)
	return dbutil.MakeDB(config.Database)
}
