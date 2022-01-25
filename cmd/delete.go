package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chanzuckerberg/happy-deploy/pkg/backend"
	"github.com/chanzuckerberg/happy-deploy/pkg/config"
	"github.com/chanzuckerberg/happy-deploy/pkg/orchestrator"
	stack_service "github.com/chanzuckerberg/happy-deploy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy-deploy/pkg/workspace_repo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete STACK_NAME",
	Short: "delete an existing stack",
	Long:  "Delete the stack with the given name.",
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	env := "rdev"

	if len(args) != 1 {
		return errors.New("incorrect number of arguments")
	}

	stackName := args[0]

	fmt.Printf("Deleting %s with settings: wait=%v force=%v\n", stackName, wait, force)

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

	taskRunner := backend.GetAwsEcs(happyConfig)

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

	// TODO check env to make sure it allows for stack deletion

	fmt.Printf("Deleting %s\n", stackName)

	stacks, err := stackService.GetStacks()
	if err != nil {
		return err
	}
	stack, ok := stacks[stackName]
	if !ok {
		return fmt.Errorf("stack %s not found", stackName)
	}

	// Run all necessary tasks before deletion
	wait := true
	showLogs := true
	taskOrchestrator := orchestrator.NewOrchestrator(happyConfig, taskRunner)
	err = taskOrchestrator.RunTasks(stack, string(backend.DeletionTask), wait, showLogs)
	if err != nil {

		var ans string
		fmt.Scanln(&ans)
		fmt.Printf("Error running tasks while trying to delete %s; Continue (y/n)? ", stackName)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; !ok {
			return err
		}
	}

	// Destroy the stack
	destroySuccess := true
	if err = stack.Destroy(); err != nil {
		// log error and set a flag, but not do not return
		fmt.Printf("Failed to destroy stack %s", err)
		destroySuccess = false
	}

	doRemoveWorkspace := false
	if !destroySuccess {
		var ans string
		fmt.Scanln(&ans)
		fmt.Printf("Error while destroying %s; resources might remain. Continue to remove workspace (y/n)? ", stackName)
		YES := map[string]bool{"Y": true, "y": true, "yes": true, "YES": true}
		if _, ok := YES[ans]; ok {
			doRemoveWorkspace = true
		}
	}

	// Remove the stack from state
	if destroySuccess || doRemoveWorkspace {
		err = stackService.Remove(stackName)
		if err != nil {
			return err
		}
		fmt.Println("Delete done")
	} else {
		fmt.Println("Delete NOT done")
	}

	return nil
}

// func runTask(stack stack_mgr.StackIface, config *config.HappyConfig, orchestrator *Orchestrator, taskType string, wait bool, showLogs bool) error {
// 	fmt.Printf("Running tasks for %s", taskType)
// 	taskOutputs, ok := config.GetData().Tasks[taskType]
// 	if !ok {
// 		fmt.Printf("Found no tasks for %s", taskType)
// 		taskOutputs = []string{}
// 	}

// 	stackOutput, err := stack.GetOutputs()
// 	if err != nil {
// 		return err
// 	}
// 	tasks := []string{}
// 	for _, taskOutput := range taskOutputs {
// 		task, ok := stackOutput[taskOutput]
// 		if ok {
// 			tasks = append(tasks, task)
// 		} else {
// 			return fmt.Errorf("WARNING task %s not found in stack %s", task, stack.GetName())
// 		}
// 	}
// 	for _, task := range tasks {
// 		orchestrator.RunTask(task, wait, showLogs)
// 	}

// 	return nil
// }
