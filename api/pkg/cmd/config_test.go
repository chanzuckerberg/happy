package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/stretchr/testify/require"
)

func TestSetConfigValueSucceed(t *testing.T) {
	testData := []struct {
		input    []*model.AppConfigPayload
		expected []*model.AppConfigPayload
	}{
		{
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different key, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key BUT with a stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key and stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// same key, different stack name, different value
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val3"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different env, different value, everything else the same
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
			},
		},
		{
			input: []*model.AppConfigPayload{
				// different app, different value, everything else the same
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("anotherapp", "rdev", "", "KEY1", "val2"),
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			configs := db.GetDB().AppConfig.Query().AllX(context.Background())

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, model.NewAppConfigPayload(config.AppName, config.Environment, config.Stack, config.Key, config.Value))
			}

			r.EqualValues(results, tc.expected)
		})
	}
}

func TestGetAppConfigSucceed(t *testing.T) {
	testData := []struct {
		input          []*model.AppConfigPayload
		lookup         *model.AppConfigLookupPayload
		expected       *model.AppConfigPayload
		expectedSource string
	}{
		{
			// stack and env configs exist, no stack specified -> return env config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY2"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			expectedSource: "environment",
		},
		{
			// stack and env configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY2", "val-foo2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			expectedSource: "stack",
		},
		{
			// only env configs exist, stack specified -> return env config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			expectedSource: "environment",
		},
		{
			// only stack configs exist, stack specified -> return stack config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:         model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected:       model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			expectedSource: "stack",
		},
		{
			// only stack configs exist, stack not specified -> return null
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack not specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expected: nil,
		},
		{
			// no configs, stack specified -> return null
			input:    []*model.AppConfigPayload{},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected: nil,
		},
		{
			// stack config exists, different stack specified -> return null
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			lookup:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expected: nil,
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			result, err := MakeConfig(db).GetResolvedAppConfig(tc.lookup)
			r.NoError(err)

			if tc.expected == nil {
				r.Nil(result)
			} else {
				r.EqualValues(result.AppConfigPayload, *tc.expected)
				r.EqualValues(result.Source, tc.expectedSource)
			}
		})
	}
}

func TestDeleteAppConfigSucceed(t *testing.T) {
	testData := []struct {
		input            []*model.AppConfigPayload
		deleteCriteria   *model.AppConfigLookupPayload
		expectedResult   *model.AppConfigPayload
		remainingConfigs []*model.AppConfigPayload
	}{
		{
			// no configs exist, no stack specified -> delete nothing
			input:            []*model.AppConfigPayload{},
			deleteCriteria:   model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// no configs exist, stack specified -> delete nothing
			input:            []*model.AppConfigPayload{},
			deleteCriteria:   model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult:   nil,
			remainingConfigs: []*model.AppConfigPayload{},
		},
		{
			// configs exist, no stack specified -> deletes env-level config
			input: []*model.AppConfigPayload{
				// order here is important, stack config needs to be created first to ensure
				// that this test fails when empty stack criteria is not added in the DeleteAppConfig
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			deleteCriteria: model.NewAppConfigLookupPayload("testapp", "rdev", "", "KEY1"),
			expectedResult: model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			remainingConfigs: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			},
		},
		{
			// configs exist, stack specified -> deletes stack-level config
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			deleteCriteria: model.NewAppConfigLookupPayload("testapp", "rdev", "foo", "KEY1"),
			expectedResult: model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
			remainingConfigs: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			result, err := MakeConfig(db).DeleteAppConfig(tc.deleteCriteria)
			r.NoError(err)

			if tc.expectedResult == nil {
				r.Nil(result)
			} else {
				r.EqualValues(result.AppConfigPayload, *tc.expectedResult)
			}

			configs := db.GetDB().AppConfig.Query().AllX(context.Background())

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, model.NewAppConfigPayload(config.AppName, config.Environment, config.Stack, config.Key, config.Value))
			}

			r.EqualValues(results, tc.remainingConfigs)
		})
	}
}

func TestGetAllAppConfigsSucceed(t *testing.T) {
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []*model.AppConfigPayload
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []*model.AppConfigPayload{},
			appMetadata: model.NewAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist, stack specified -> ignores stack and returns all configs for app and env
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			appMetadata: model.NewAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
			},
			appMetadata: model.NewAppMetadata("testapp", "staging", ""),
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, val := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(val)
				r.NoError(err)
			}

			configs, err := MakeConfig(db).GetAllAppConfigs(tc.appMetadata)
			r.NoError(err)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
			}

			r.ElementsMatch(results, tc.expected)
		})
	}
}

func TestGetAppConfigsForEnvSucceed(t *testing.T) {
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []*model.AppConfigPayload
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []*model.AppConfigPayload{},
			appMetadata: model.NewAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist, stack specified -> ignores stack and returns only env-level configs
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			appMetadata: model.NewAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns only configs for specified env
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
			},
			appMetadata: model.NewAppMetadata("testapp", "staging", ""),
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, val := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(val)
				r.NoError(err)
			}

			configs, err := MakeConfig(db).GetAppConfigsForEnv(tc.appMetadata)
			r.NoError(err)

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, &config.AppConfigPayload)
				r.Equal(config.Source, "environment")
			}

			r.ElementsMatch(results, tc.expected)
		})
	}
}

func TestGetAppConfigsForStackSucceed(t *testing.T) {
	type expected struct {
		payload *model.AppConfigPayload
		source  string
	}
	testData := []struct {
		input       []*model.AppConfigPayload
		expected    []expected
		appMetadata *model.AppMetadata
	}{
		{
			// no configs exist -> return empty array
			input:       []*model.AppConfigPayload{},
			expected:    []expected{},
			appMetadata: model.NewAppMetadata("testapp", "rdev", ""),
		},
		{
			// stack configs exist but query is without stack -> return empty array
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			expected:    []expected{},
			appMetadata: model.NewAppMetadata("testapp", "rdev", ""),
		},
		{
			// configs exist -> returns stack and env-level configs with overrides applied
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
			},
			expected: []expected{
				{
					payload: model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
					source:  "stack",
				},
				{
					payload: model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
					source:  "environment",
				},
			},
			appMetadata: model.NewAppMetadata("testapp", "rdev", "foo"),
		},
		{
			// configs exist, different env specified -> returns correct stack and env-level configs with overrides applied
			input: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val-foo"),
				model.NewAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val-bar"),
			},
			expected: []expected{
				{
					payload: model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val2"),
					source:  "environment",
				},
				{
					payload: model.NewAppConfigPayload("testapp", "staging", "stg", "KEY2", "val-stg"),
					source:  "stack",
				},
			},
			appMetadata: model.NewAppMetadata("testapp", "staging", "stg"),
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, val := range tc.input {
				_, err := MakeConfig(db).SetConfigValue(val)
				r.NoError(err)
			}

			configs, err := MakeConfig(db).GetAppConfigsForStack(tc.appMetadata)
			r.NoError(err)

			results := []expected{}
			for _, config := range configs {
				results = append(results, expected{
					payload: &config.AppConfigPayload,
					source:  config.Source,
				})
			}

			r.ElementsMatch(results, tc.expected)
		})
	}
}

func TestCopyAppConfigSucceed(t *testing.T) {
	testData := []struct {
		seeds       []*model.AppConfigPayload
		copyPayload *model.CopyAppConfigPayload
		expected    []*model.AppConfigPayload
	}{
		{
			// no configs exist -> nothing gets copied
			seeds:       []*model.AppConfigPayload{},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected:    []*model.AppConfigPayload{},
		},
		{
			// configs exist but don't match -> nothing gets copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// configs exist but don't match -> nothing gets copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// matching env config exists -> copy env config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
			},
		},
		{
			// matching stack config exists -> copy stack config
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "staging", "", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
			},
		},
		{
			// test copying to different stack in same env
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val2"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "rdev", "foo", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// test copying to same stack in same env
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
			copyPayload: model.NewCopyAppConfigPayload("testapp", "rdev", "bar", "rdev", "bar", "KEY1"),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY1", "val1"),
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			_, err = MakeConfig(db).CopyAppConfig(tc.copyPayload)
			r.NoError(err)

			configs := db.GetDB().AppConfig.Query().AllX(context.Background())

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, model.NewAppConfigPayload(config.AppName, config.Environment, config.Stack, config.Key, config.Value))
			}

			r.EqualValues(tc.expected, results)
		})
	}
}

func TestAppConfigDiffSucceed(t *testing.T) {
	testData := []struct {
		seeds       []*model.AppConfigPayload
		diffPayload *model.AppConfigDiffPayload
		expected    []string
	}{
		{
			// no configs -> no diff
			seeds:       []*model.AppConfigPayload{},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected:    []string{},
		},
		{
			// config exists only for stack and no stack specified -> no diff
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected:    []string{},
		},
		{
			// config exists only for env, stack specified -> env config is part of stack -> key returned as diff
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "foo", "staging", ""),
			expected:    []string{"KEY1"},
		},
		{
			// same configs in each -> no diff
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "rdev-foo-val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "staging-val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected:    []string{},
		},
		{
			// configs exists only in source -> key returned as diff
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "foo-val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY2", "bar-val2"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "foo", "staging", ""),
			expected:    []string{"KEY1"},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			results, err := MakeConfig(db).AppConfigDiff(tc.diffPayload)
			r.NoError(err)

			r.EqualValues(tc.expected, results)
		})
	}
}

func TestCopyAppConfigDiffSucceed(t *testing.T) {
	testData := []struct {
		seeds       []*model.AppConfigPayload
		diffPayload *model.AppConfigDiffPayload
		expected    []*model.AppConfigPayload
	}{
		{
			// no configs -> no diff
			seeds:       []*model.AppConfigPayload{},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected:    []*model.AppConfigPayload{},
		},
		{
			// config exists only for stack and no stack specified -> no copies
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
		},
		{
			// config exists only for env, stack specified -> env config is part of stack -> config copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "foo", "staging", ""),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
			},
		},
		{
			// same configs in each -> no copies
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "rdev-foo-val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "staging-val1"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "foo", "staging", ""),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "rdev-foo-val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "staging-val1"),
			},
		},
		{
			// configs exists only in source -> configs copied
			seeds: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "foo-val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY2", "bar-val2"),
			},
			diffPayload: model.NewAppConfigDiffPayload("testapp", "rdev", "", "staging", ""),
			expected: []*model.AppConfigPayload{
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "rdev", "", "KEY2", "val2"),
				model.NewAppConfigPayload("testapp", "rdev", "foo", "KEY1", "foo-val1"),
				model.NewAppConfigPayload("testapp", "rdev", "bar", "KEY2", "bar-val2"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY1", "val1"),
				model.NewAppConfigPayload("testapp", "staging", "", "KEY2", "val2"),
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(t)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeConfig(db).SetConfigValue(input)
				r.NoError(err)
			}

			_, err = MakeConfig(db).CopyAppConfigDiff(tc.diffPayload)
			r.NoError(err)

			configs := db.GetDB().AppConfig.Query().AllX(context.Background())

			results := []*model.AppConfigPayload{}
			for _, config := range configs {
				results = append(results, model.NewAppConfigPayload(config.AppName, config.Environment, config.Stack, config.Key, config.Value))
			}

			r.ElementsMatch(tc.expected, results)
		})
	}
}
