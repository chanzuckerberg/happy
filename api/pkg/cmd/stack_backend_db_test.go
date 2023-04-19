package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/stretchr/testify/require"
)

func MakeAppStack(appName, env, stack string) model.AppStack {
	return model.AppStack{
		AppMetadata: *model.NewAppMetadata(appName, env, stack),
	}
}

func TestCreateStackSuccess(t *testing.T) {
	testData := []struct {
		seeds    []model.AppStackPayload
		expected []model.AppMetadata
	}{
		{
			// should create one stack
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			},
			expected: []model.AppMetadata{
				*model.NewAppMetadata("testapp", "rdev", "mystack"),
			},
		},
		{
			// should create multiple stacks
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", "", "", "", "", ""),
				model.MakeAppStackPayload("testapp", "staging", "mystack2", "", "", "", "", ""),
			},
			expected: []model.AppMetadata{
				*model.NewAppMetadata("testapp", "rdev", "mystack"),
				*model.NewAppMetadata("testapp", "rdev", "mystack2"),
				*model.NewAppMetadata("testapp", "staging", "mystack2"),
			},
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(r)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeStackBackendDB(db).CreateOrUpdateAppStack(input)
				r.NoError(err)
			}

			stacks := []*model.AppStack{}
			db.GetDB().Find(&stacks)

			results := []model.AppMetadata{}
			for _, stack := range stacks {
				results = append(results, stack.AppMetadata)
			}

			r.EqualValues(tc.expected, results)
		})
	}
}

func TestDeleteStackSuccess(t *testing.T) {
	testData := []struct {
		seeds         []model.AppStackPayload
		stackPayload  model.AppStackPayload
		expectDeleted bool
	}{
		{
			// should return nil when no stacks matched
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			},
			stackPayload:  model.MakeAppStackPayload("testapp", "rdev", "mystack2", "", "", "", "", ""),
			expectDeleted: false,
		},
		{
			// should delete a matching stack
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			},
			stackPayload:  model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			expectDeleted: true,
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(r)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeStackBackendDB(db).CreateOrUpdateAppStack(input)
				r.NoError(err)
			}

			res, err := MakeStackBackendDB(db).DeleteAppStack(tc.stackPayload)
			r.NoError(err)

			if tc.expectDeleted {
				r.Equal(tc.stackPayload.AppMetadata, res.AppMetadata)
			} else {
				r.Nil(res)
			}
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
			seeds:        []model.AppStackPayload{},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			expected:     0,
		},
		{
			// should return an empty list if no stacks match
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			},
			stackPayload: model.MakeAppStackPayload("misspelled app name", "rdev", "mystack", "", "", "", "", ""),
			expected:     0,
		},
		{
			// should return a single item
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "mystack", "", "", "", "", ""),
			expected:     1,
		},
		{
			// should return all the items (with the empty string provided)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", "", "", "", "", ""),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", "", "", "", "", ""),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "rdev", "", "", "", "", "", ""),
			expected:     2,
		},
		{
			// should return all the items (without the stack provided)
			seeds: []model.AppStackPayload{
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", "", "", "", "", ""),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", "", "", "", "", ""),
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
				model.MakeAppStackPayload("testapp", "rdev", "mystack1", "", "", "", "", ""),
				model.MakeAppStackPayload("testapp", "rdev", "mystack2", "", "", "", "", ""),
			},
			stackPayload: model.MakeAppStackPayload("testapp", "staging", "", "", "", "", "", ""),
			expected:     0,
		},
	}

	for idx, testCase := range testData {
		tc := testCase
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			db := MakeTestDB(r)
			err := db.AutoMigrate()
			r.NoError(err)

			for _, input := range tc.seeds {
				_, err := MakeStackBackendDB(db).CreateOrUpdateAppStack(input)
				r.NoError(err)
			}

			stacks, err := MakeStackBackendDB(db).GetAppStacks(
				context.Background(),
				model.MakeAppStackPayload(tc.stackPayload.AppName, tc.stackPayload.Environment, "", "czi-si", "us-west-2", "fargate", "", ""),
			)
			r.NoError(err)
			r.Len(stacks, tc.expected)
		})
	}
}
