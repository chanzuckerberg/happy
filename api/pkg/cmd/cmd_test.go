package cmd

import (
	"fmt"
	"sync"
	"testing"

	"github.com/chanzuckerberg/happy/api/pkg/ent/enttest"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/chanzuckerberg/happy/api/pkg/store"
	"github.com/google/uuid"
)

var (
	mu sync.Mutex
)

func MakeTestDB(t *testing.T) *store.DB {
	config := setup.GetConfiguration()

	// Even with a UUID in the data source name this is not thread safe so we need to use a mutex to prevent concurrent access
	mu.Lock()
	client := enttest.Open(t, "sqlite3", fmt.Sprintf("file:memdb%s?mode=memory&cache=shared&_fk=1", uuid.NewString()))
	mu.Unlock()

	return store.MakeDB(config.Database).WithClient(client)
}
