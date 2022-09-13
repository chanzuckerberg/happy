package cmd

import (
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy-api/pkg/dbutil"
	"github.com/chanzuckerberg/happy-api/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestCreateStackFailures(t *testing.T) {
	testData := []struct {
		seeds        []model.AppStackPayload
		stackPayload model.AppStackPayload
	}{
		{
			// should throw an error if trying to create a duplicate stack
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
		},
		{
			// should throw an error if trying to create a duplicate stack (even if disabled)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "mystack", false),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			db := dbutil.MakeDB(dbutil.WithErrorLogLevel(), dbutil.WithInMemorySQLDriver())
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range testCase.seeds {
				_, err := MakeStack(db).CreateAppStack(input)
				r.NoError(err)
			}

			_, err = MakeStack(db).CreateAppStack(testCase.stackPayload)
			r.Error(err)
		})
	}
}

func TestUpdateStackFailures(t *testing.T) {
	testData := []struct {
		seeds        []model.AppStackPayload
		stackPayload model.AppStackPayload
	}{
		{
			// should throw an error if trying to update a stack that doesn't exist
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			},
			stackPayload: model.MakeAppStackPayload("misspelled app name", "rdev", "mystack", true),
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			db := dbutil.MakeDB(dbutil.WithInfoLogLevel(), dbutil.WithInMemorySQLDriver())
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range testCase.seeds {
				_, err := MakeStack(db).CreateAppStack(input)
				r.NoError(err)
			}

			_, err = MakeStack(db).UpdateAppStack(testCase.stackPayload)
			r.Error(err)
		})
	}
}

func TestGetStackFailures(t *testing.T) {
	testData := []struct {
		seeds        []model.AppStackPayload
		stackPayload model.AppStackPayload
		expected     int
	}{
		{
			// should return an empty list if no stacks match
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			},
			stackPayload: model.MakeAppStackPayload("misspelled app name", "rdev", "mystack", true),
			expected:     0,
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			db := dbutil.MakeDB(dbutil.WithInfoLogLevel(), dbutil.WithInMemorySQLDriver())
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range testCase.seeds {
				_, err := MakeStack(db).CreateAppStack(input)
				r.NoError(err)
			}

			stacks, err := MakeStack(db).GetAppStacks(testCase.stackPayload)
			r.NoError(err)
			r.Len(stacks, testCase.expected)
		})
	}
}

func TestGetStackSuccesses(t *testing.T) {
	testData := []struct {
		seeds        []model.AppStackPayload
		stackPayload model.AppStackPayload
		expected     int
	}{
		{
			// should return a single item
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "mystack", true),
			expected:     1,
		},
		{
			// should return all the items (with the empty string provided)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", true),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", true),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "", true),
			expected:     2,
		},
		{
			// should return all the items (without the stack provided)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", true),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", true),
			},
			stackPayload: model.AppStackPayload{
				AppMetadata: model.AppMetadata{
					App:         model.App{AppName: "testapp"},
					Environment: "rdev",
				},
			},
			expected: 2,
		},
		{
			// should return no items
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", true),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", true),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "", false),
			expected:     0,
		},
		{
			// should return only enabled items even when enabled is not set (if not set, the default should be true)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", true),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", true),
				model.MakeAppStackPayload("testapp", "rdev", "mystack3", false),
			},
			stackPayload: model.AppStackPayload{
				AppMetadata: *model.NewAppMetadata("testapp", "rdev", ""),
			},
			expected: 2,
		},
	}

	for idx, testCase := range testData {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := require.New(t)
			db := dbutil.MakeDB(dbutil.WithInfoLogLevel(), dbutil.WithInMemorySQLDriver())
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range testCase.seeds {
				_, err := MakeStack(db).CreateAppStack(input)
				r.NoError(err)
			}

			stacks, err := MakeStack(db).GetAppStacks(testCase.stackPayload)
			r.NoError(err)
			r.Len(stacks, testCase.expected)
		})
	}
}
