package cmd

import (
	"github.com/chanzuckerberg/happy/api/pkg/dbutil"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/stretchr/testify/require"
)

func MakeTestDB(r *require.Assertions) *dbutil.DB {
	config := setup.GetConfiguration()
	return dbutil.MakeDB(config.Database)
}
