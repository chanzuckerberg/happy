package dbutil

import (
	"context"
	"log"
	"sync"

	"github.com/chanzuckerberg/happy/api/pkg/ent"
	_ "github.com/chanzuckerberg/happy/api/pkg/ent/runtime"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	once     sync.Once
	dbClient *ent.Client
	Config   setup.DatabaseConfiguration
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
	d.dbClient = client
	return d
}

func (d *DB) GetDB() *ent.Client {
	d.once.Do(func() {
		// if this was set with WithClient we do not want to override it
		if d.dbClient != nil {
			return
		}

		opts := []ent.Option{}
		if d.Config.LogLevel == "debug" {
			opts = append(opts, ent.Debug())
		}

		var err error
		d.dbClient, err = ent.Open(d.Config.Driver.String(), d.Config.DataSourceName, opts...)
		if err != nil {
			log.Fatalf("ENT failed to connect to the DB: %v", err)
		}
	})

	return d.dbClient
}

func (d *DB) AutoMigrate() error {
	client := d.GetDB()
	ctx := context.Background()

	// Run the auto migration tool.
	return client.Schema.Create(ctx)
}
