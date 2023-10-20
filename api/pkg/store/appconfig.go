package store

import (
	"context"

	"github.com/chanzuckerberg/happy/api/pkg/ent"
	"github.com/chanzuckerberg/happy/api/pkg/ent/appconfig"
	"github.com/pkg/errors"
)

func (d *DB) GetAppConfigsForStack(ctx context.Context, appName, env, stack string) ([]*ent.AppConfig, error) {
	// get all appconfigs for the app/env and order by key, then by stack DESC. Take the first item for each key
	db := d.GetDB()
	records, err := appEnvScopedQuery(db.AppConfig, appName, env).
		Where(appconfig.StackIn(stack, "")).
		Order(ent.Asc(appconfig.FieldKey), ent.Desc(appconfig.FieldStack)).
		All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[GetAppConfigsForStack] unable to query app configs")
	}

	// we'll get at most 2 config records for each key (one for env and one for stack), so we'll use a map to dedupe
	// and select the stack record if it exists (since we order by stack DESC) and the env record otherwise
	resolvedMap := map[string]*ent.AppConfig{}
	for _, record := range records {
		if _, ok := resolvedMap[record.Key]; !ok {
			resolvedMap[record.Key] = record
		}
	}

	results := []*ent.AppConfig{}
	for _, record := range resolvedMap {
		results = append(results, record)
	}

	return results, nil
}

func appEnvScopedQuery(client *ent.AppConfigClient, appName, env string) *ent.AppConfigQuery {
	return client.Query().Where(
		appconfig.AppName(appName),
		appconfig.Environment(env),
	)
}
