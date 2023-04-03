package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	stackservice "github.com/chanzuckerberg/happy/cli/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Printer interface {
	PrintStacks(stackInfos []stackservice.StackInfo) error
	PrintResources(resources []util.ManagedResource) error
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
	Tag         string `header:"Tags"`
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

func Stack2Console(stack stackservice.StackInfo) StackConsoleInfo {
	endpoints := []string{}
	if endpoint, ok := stack.Outputs["frontend_url"]; ok {
		endpoints = append(endpoints, endpoint)
	} else if svc_endpoints, ok := stack.Outputs["service_endpoints"]; ok {
		endpointmap := map[string]string{}
		err := json.Unmarshal([]byte(svc_endpoints), &endpointmap)
		if err != nil {
			logrus.Errorf("Unable to decode endpoints: %s", err.Error())
		} else {
			uniqueMap := map[string]bool{}
			for _, endpoint := range endpointmap {
				uniqueMap[endpoint] = true
			}
			for endpoint := range uniqueMap {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	return StackConsoleInfo{
		Name:        stack.Name,
		Owner:       stack.Owner,
		Tag:         stack.Tag,
		Status:      stack.Status,
		FrontendUrl: strings.Join(endpoints, "\n"),
		LastUpdated: stack.LastUpdated,
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

func (p *TextPrinter) PrintStacks(stackInfos []stackservice.StackInfo) error {
	printer := util.NewTablePrinter()

	stacks := make([]StackConsoleInfo, 0)
	for _, stackInfo := range stackInfos {
		stacks = append(stacks, Stack2Console(stackInfo))
	}
	printer.Print(stacks)

	return nil
}

func (p *TextPrinter) PrintResources(resources []util.ManagedResource) error {
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

func (p *JSONPrinter) PrintStacks(stackInfos []stackservice.StackInfo) error {
	b, err := json.Marshal(stackInfos)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *JSONPrinter) PrintResources(resources []util.ManagedResource) error {
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

func (p *YAMLPrinter) PrintStacks(stackInfos []stackservice.StackInfo) error {
	b, err := yaml.Marshal(stackInfos)
	if err != nil {
		return err
	}
	PrintOutput(string(b))
	return nil
}

func (p *YAMLPrinter) PrintResources(resources []util.ManagedResource) error {
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
