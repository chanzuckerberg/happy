package dbutil

import (
	"context"
	"log"
	"sync"

	_ "github.com/chanzuckerberg/happy/api/pkg/ent/runtime"

	"github.com/chanzuckerberg/happy/api/pkg/ent"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
)

type DB struct {
	onceEnt sync.Once
	dbEnt   *ent.Client
	Config  setup.DatabaseConfiguration
}

// MakeDB returns a pointer because we store sync.Once inside it
// so we don't need to initialize the database many times.
func MakeDB(cfg setup.DatabaseConfiguration) *DB {
	db := &DB{
		Config: cfg,
	}
	return db
}

func (d *DB) WithClient(client *ent.Client) *DB {
	d.dbEnt = client
	return d
}

func (d *DB) GetDBEnt() *ent.Client {
	d.onceEnt.Do(func() {
		// if this was set with WithClient we do not want to override it
		if d.dbEnt != nil {
			return
		}

		var err error
		d.dbEnt, err = ent.Open(d.Config.Driver.String(), d.Config.DataSourceName)
		if err != nil {
			log.Fatalf("ENT failed to connect to the DB: %v", err)
		}
	})

	return d.dbEnt
}

func (d *DB) AutoMigrate() error {
	client := d.GetDBEnt()
	ctx := context.Background()

	// Run the auto migration tool.
	return client.Schema.Create(ctx)
}
