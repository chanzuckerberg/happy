package output

import (
	"encoding/json"
	"fmt"
	"os"

	stackservice "github.com/chanzuckerberg/happy/pkg/stack_mgr"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Printer interface {
	PrintStacks(stackInfos []stackservice.StackInfo) error
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

func (p *TextPrinter) PrintStacks(stackInfos []stackservice.StackInfo) error {
	headings := []string{"Name", "Owner", "Tags", "Status", "URLs", "LastUpdated"}
	tablePrinter := util.NewTablePrinter(headings)

	for _, stackInfo := range stackInfos {
		tablePrinter.AddRow(stackInfo.Name, stackInfo.Owner, stackInfo.Tag, stackInfo.Status, stackInfo.Outputs["frontend_url"], stackInfo.LastUpdated)
	}
	tablePrinter.Print()
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

func (p *YAMLPrinter) Fatal(err error) {
	PrintError(err)
}

func PrintError(err error) {
	os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
}

func PrintOutput(output string) {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", output))
}
