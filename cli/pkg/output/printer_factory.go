package output

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type Printer interface {
	PrintStacks(ctx context.Context, stackInfos []*model.AppStackResponse) error
	PrintResources(ctx context.Context, resources []util.ManagedResource) error
	Fatal(err error)
}

type TextPrinter struct{}
type JSONPrinter struct{}
type YAMLPrinter struct{}

func NewPrinter(outputFormat string) Printer {
	switch outputFormat {
	case "json":
		return &JSONPrinter{}
	case "yaml":
		return &YAMLPrinter{}
	default:
		return &TextPrinter{}
	}
}

type StackConsoleInfo struct {
	Name        string `header:"Name"`
	Owner       string `header:"Owner"`
	App         string `header:"App"`
	Repo        string `header:"Repo"`
	Branch      string `header:"Branch"`
	Hash        string `header:"Hash"`
	Status      string `header:"Status"`
	FrontendUrl string `header:"URLs"`
	LastUpdated string `header:"LastUpdated"`
}

type ResourceConsoleInfo struct {
	Module    string   `header:"Module"`
	Name      string   `header:"Name"`
	Type      string   `header:"Type"`
	ManagedBy string   `header:"ManagedBy"`
	Instances []string `header:"Instances"`
}

func Stack2Console(ctx context.Context, stack *model.AppStackResponse) StackConsoleInfo {
	endpoints := []string{}
	stackEndpoints := stack.Endpoints
	uniqueMap := map[string]bool{}
	for _, endpoint := range stackEndpoints {
		uniqueMap[endpoint] = true
	}
	for endpoint := range uniqueMap {
		// filter out the k8s cluster endpoints
		if strings.Contains(endpoint, "svc.cluster") {
			continue
		}
		endpoints = append(endpoints, endpoint)
	}

	abbrevRepo := strings.TrimPrefix(strings.TrimSuffix(stack.GitRepo, ".git"), "git@github.com:")
	abbrevOwner := strings.TrimSuffix(stack.Owner, "@chanzuckerberg.com")
	updatedTime, err := time.Parse("2006-01-02T15:04:05-07:00", stack.LastUpdated)
	abbrevLastUpdated := stack.LastUpdated
	if err == nil {
		abbrevLastUpdated = time.Since(updatedTime).Truncate(time.Second * 1).String()
	}
	return StackConsoleInfo{
		Name:        stack.Stack,
		Owner:       abbrevOwner,
		App:         stack.AppName,
		Repo:        abbrevRepo,
		Branch:      stack.GitBranch,
		Hash:        stack.GitSHA,
		Status:      stack.TFEWorkspaceStatus,
		FrontendUrl: strings.Join(endpoints, "\n"),
		LastUpdated: abbrevLastUpdated,
	}
}

func Resource2Console(resource util.ManagedResource) ResourceConsoleInfo {
	return ResourceConsoleInfo{
		Name:      resource.Name,
		Module:    resource.Module,
		Type:      resource.Type,
		ManagedBy: resource.ManagedBy,
		Instances: resource.Instances,
	}
}

func (p *TextPrinter) PrintStacks(ctx context.Context, stackInfos []*model.AppStackResponse) error {
	if len(stackInfos) == 0 {
		logrus.Info("No stacks found")
	}
	printer := util.NewTablePrinter()

	stacks := make([]StackConsoleInfo, 0)
	for _, stackInfo := range stackInfos {
		stacks = append(stacks, Stack2Console(ctx, stackInfo))
	}
	printer.Print(stacks)

	return nil
}

func (p *TextPrinter) PrintResources(ctx context.Context, resources []util.ManagedResource) error {
	printer := util.NewTablePrinter()

	resourceInfos := make([]ResourceConsoleInfo, 0)
	for _, resource := range resources {
		resourceInfos = append(resourceInfos, Resource2Console(resource))
	}
	printer.Print(resourceInfos)

	return nil
}

func (p *TextPrinter) Fatal(err error) {
	logrus.Fatal(err)
}

func (p *JSONPrinter) PrintStacks(ctx context.Context, stackInfos []*model.AppStackResponse) error {
	b, err := json.Marshal(stackInfos)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *JSONPrinter) PrintResources(ctx context.Context, resources []util.ManagedResource) error {
	b, err := json.Marshal(resources)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *JSONPrinter) Fatal(err error) {
	PrintError(err)
}

func (p *YAMLPrinter) PrintStacks(ctx context.Context, stackInfos []*model.AppStackResponse) error {
	b, err := yaml.Marshal(stackInfos)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *YAMLPrinter) PrintResources(ctx context.Context, resources []util.ManagedResource) error {
	b, err := yaml.Marshal(resources)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *YAMLPrinter) Fatal(err error) {
	PrintError(err)
}

func PrintError(err error) {
	os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
}

func PrintOutput(output string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", output))
}
