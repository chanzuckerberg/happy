package cmd

import (
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy/pkg/backend"
	"github.com/chanzuckerberg/happy/pkg/config"
	stack_service "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/chanzuckerberg/happy/pkg/workspace_repo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	createTag string
	wait      bool
	force     bool
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createTag, "tag", "", "Tag name for docker image. Leave empty to generate one")
	createCmd.Flags().BoolVar(&wait, "wait", true, "Wait for this cmd to complete")
	createCmd.Flags().BoolVar(&force, "force", false, "Ignore the already-exists errors")
}

var createCmd = &cobra.Command{
	Use:   "create STACK_NAME",
	Short: "create new stack",
	Long:  "Create a new stack with a given tag.",
	RunE:  runCreate,
	Args:  cobra.ExactArgs(1),
}

func runCreate(cmd *cobra.Command, args []string) error {
	env := "rdev"

	stackName := args[0]

	fmt.Printf("Creating %s with settings: wait=%v force=%v\n", stackName, wait, force)

	happyConfigPath, ok := os.LookupEnv("HAPPY_CONFIG_PATH")
	if !ok {
		return errors.New("please set env var HAPPY_CONFIG_PATH")
	}

	_, ok = os.LookupEnv("HAPPY_PROJECT_ROOT")
	if !ok {
		return errors.New("please set env var HAPPY_PROJECT_ROOT")
	}

	happyConfig, err := config.NewHappyConfig(happyConfigPath, env)
	if err != nil {
		return err
	}

	url, err := happyConfig.TfeUrl()
	if err != nil {
		return err
	}
	org, err := happyConfig.TfeOrg()
	if err != nil {
		return err
	}
	workspaceRepo, err := workspace_repo.NewWorkspaceRepo(url, org)
	if err != nil {
		return err
	}
	paramStoreBackend := backend.GetAwsBackend(happyConfig)
	stackService := stack_service.NewStackService(happyConfig, paramStoreBackend, workspaceRepo)

	existingStacks, err := stackService.GetStacks()
	if err != nil {
		return err
	}
	if _, ok := existingStacks[stackName]; ok {
		if !force {
			return errors.Errorf("stack %s already exists", stackName)
		}
	}

	stackMeta := stackService.NewStackMeta(stackName)
	secretArn := happyConfig.GetSecretArn()
	if err != nil {
		return err
	}
	metaTag := map[string]string{"happy/meta/configsecret": secretArn}
	stackMeta.Load(metaTag)

	if createTag == "" {
		createTag, err = util.GenerateTag(happyConfig)
		if err != nil {
			return err
		}

		// invoke push cmd
		fmt.Printf("Pushing images with tags %s...\n", createTag)
		err := runPush(createTag)
		if err != nil {
			return errors.Errorf("failed to push image: %s", err)
		}
	}
	stackMeta.Update(createTag, stackService)
	fmt.Printf("Creating %s\n", stackName)

	stack, err := stackService.Add(stackName)
	if err != nil {
		return err
	}
	fmt.Printf("setting stackMeta %v\n", stackMeta)
	stack.SetMeta(stackMeta)

	err = stack.Apply()
	if err != nil {
		return err
	}

	autoRunMigration := happyConfig.AutoRunMigration()
	if err != nil {
		fmt.Println("WARNING autoRunMigration flag not set, defaulting to false")
	}

	if autoRunMigration {
		runMigrate(stackName)
	}
	// TODO migrate db here

	stack.PrintOutputs()

	return nil
}
