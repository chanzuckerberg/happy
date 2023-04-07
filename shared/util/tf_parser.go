package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

func ParseServices(dir string) (map[string]bool, error) {
	//var modules []map[string]interface{}
	var services map[string]bool = map[string]bool{}
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
		},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".terraform" {
				return filepath.SkipDir
			}
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			f, diags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return errors.Wrapf(diags.Errs()[0], "failed to parse %s", path)
			}

			content, _, contentDiags := f.Body.PartialContent(schema)
			if contentDiags.HasErrors() {
				return errors.New("Terraform code has errors")
			}
			for _, block := range content.Blocks {
				if block.Type == "module" {
					attrs, _ := block.Body.JustAttributes()
					if sourceAttr, ok := attrs["source"]; ok {
						source, _ := sourceAttr.Expr.(*hclsyntax.TemplateExpr).Parts[0].Value(nil)
						if strings.Contains(source.AsString(), "modules/happy-stack-eks") || strings.Contains(source.AsString(), "modules/happy-stack-ecs") {
							if servicesAttr, ok := attrs["services"]; ok {
								switch servicesAttr.Expr.(type) {
								case *hclsyntax.ObjectConsExpr:
									for _, item := range servicesAttr.Expr.(*hclsyntax.ObjectConsExpr).Items {
										key, _ := item.KeyExpr.Value(nil)
										services[key.AsString()] = true
									}
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return services, errors.Wrap(err, "failed to parse terraform files")
	}
	return services, nil
}
