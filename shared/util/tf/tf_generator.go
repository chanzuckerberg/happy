package tf

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type provider struct {
	Name    string
	Source  string
	Version string
}

type variable struct {
	Name        string
	Type        string
	Description string
	Default     cty.Value
}

var requiredProviders []provider = []provider{
	{
		Name:    "aws",
		Source:  "hashicorp/aws",
		Version: ">= 4.45",
	},
	{
		Name:    "kubernetes",
		Source:  "hashicorp/kubernetes",
		Version: ">= 2.16",
	},
	{
		Name:    "datadog",
		Source:  "datadog/datadog",
		Version: ">= 3.20.0",
	},
	{
		Name:    "happy",
		Source:  "chanzuckerberg/happy",
		Version: ">= 0.53.5",
	},
}

var requiredVariables []variable = []variable{
	{
		Name:        "aws_account_id",
		Type:        "string",
		Description: "AWS account ID to apply changes to",
	},
	{
		Name:        "k8s_cluster_id",
		Type:        "string",
		Description: "EKS K8S Cluster ID",
	},
	{
		Name:        "k8s_namespace",
		Type:        "string",
		Description: "K8S namespace for this stack",
	},
	{
		Name:        "aws_role",
		Type:        "string",
		Description: "Name of the AWS role to assume to apply changes",
	},
	{
		Name:        "image_tag",
		Type:        "string",
		Description: "Please provide an image tag",
	},
	{
		Name:        "image_tags",
		Type:        "string",
		Description: "Override the default image tags (json-encoded map)",
		Default:     cty.StringVal("{}"),
	},
	{
		Name:        "stack_name",
		Type:        "string",
		Description: "Happy Path stack name",
	},
	{
		Name:        "wait_for_steady_state",
		Type:        "bool",
		Description: "Should terraform block until k8s deployment reaches a steady state?",
		Default:     cty.BoolVal(true),
	},
}

const requiredTerraformVersion = ">= 1.3"

type TfGenerator struct {
	happyConfig *config.HappyConfig
}

func NewTfGenerator(happyConfig *config.HappyConfig) TfGenerator {
	return TfGenerator{
		happyConfig: happyConfig,
	}
}

func (tf *TfGenerator) GenerateMain(srcDir, moduleSource string, variables []Variable) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "main.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()
	moduleBlockBody := rootBody.AppendNewBlock("module", []string{"stack"}).Body()

	moduleBlockBody.SetAttributeValue("source", cty.StringVal(moduleSource))

	// Sort module variables alphabetically
	sort.SliceStable(variables, func(i, j int) bool {
		return strings.Compare(variables[i].Name, variables[j].Name) < 0
	})

	for _, variable := range variables {
		switch variable.Name {
		case "image_tag", "stack_name", "k8s_namespace":
			// Assign these module variables to the corresponding stack variables
			tokens := hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: variable.Name},
			})
			moduleBlockBody.SetAttributeRaw(variable.Name, tokens)
		case "image_tags":
			// Assign image_tags variable to "jsondecode(var.image_tag)"
			tokens := hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: variable.Name},
			})
			tokens = hclwrite.TokensForFunctionCall("jsondecode", tokens)
			moduleBlockBody.SetAttributeRaw(variable.Name, tokens)
		case "app_name":
			// Set the app name based on happy config
			moduleBlockBody.SetAttributeValue(variable.Name, cty.StringVal(tf.happyConfig.App()))
		case "deployment_stage":
			// Set the deployment stage to the current happy environment value
			moduleBlockBody.SetAttributeValue(variable.Name, cty.StringVal(tf.happyConfig.GetEnv()))
		case "stack_prefix":
			// Assign stack_prefix variable to "/${var.stack_name}"
			moduleBlockBody.SetAttributeRaw(variable.Name, tokens("\"/${var.stack_name}\""))
		case "routing_method":
			if !variable.Default.IsNull() {
				moduleBlockBody.SetAttributeValue(variable.Name, variable.Default)
			} else {
				moduleBlockBody.SetAttributeValue(variable.Name, cty.StringVal("DOMAIN"))
			}
		case "tasks":
			if !variable.Default.IsNull() {
				moduleBlockBody.SetAttributeValue(variable.Name, variable.Default)
			}
		case "services":
			if !variable.Type.IsMapType() {
				return errors.Errorf("services variable must be an object type")
			}

			values := map[string]cty.Value{}
			defaultValues := variable.TypeDefaults.Children[""].DefaultValues
			for _, service := range tf.happyConfig.GetServices() {
				elem := map[string]cty.Value{}

				// Sort the service attributes alphabetically
				attributeNames := reflect.ValueOf(variable.Type.ElementType().AttributeTypes()).MapKeys()
				sort.SliceStable(attributeNames, func(i, j int) bool {
					return strings.Compare(attributeNames[i].String(), attributeNames[j].String()) < 0
				})

				for i := range attributeNames {
					k := attributeNames[i].String()
					ty := variable.Type.ElementType().AttributeTypes()[k]
					if _, ok := elem[k]; !ok {
						// Forcefully populate the service name
						if ty.IsPrimitiveType() && k == "name" {
							elem[k] = cty.StringVal(service)
						}

						// If default values are known, populate them
						if defaultValue, ok := defaultValues[k]; ok {
							if !defaultValue.IsNull() {
								elem[k] = defaultValue
							}
						}
					}
				}

				values[service] = cty.ObjectVal(elem)
			}

			val := cty.MapVal(values)
			moduleBlockBody.SetAttributeValue(variable.Name, val)
		default:
			if !variable.Default.IsNull() {
				// Newly added variable has a default value, use it
				moduleBlockBody.SetAttributeValue(variable.Name, variable.Default)
			} else {
				// If newly added variables to the module don't have defaults, we won't be albe to autogenerate HCL code
				return errors.Errorf("Unable to find a value for required variable %s", variable.Name)
			}
		}
	}

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func (tf *TfGenerator) GenerateProviders(srcDir string) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "providers.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()
	err = tf.generateAwsProvider(rootBody, "", "${var.aws_account_id}", "${var.aws_role}")
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code for AWS provider")
	}

	err = tf.generateAwsProvider(rootBody, "czi-si", "626314663667", "tfe-si")
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code for AWS provider")
	}

	eksBody := rootBody.AppendNewBlock("data", []string{"aws_eks_cluster", "cluster"}).Body()
	eksBody.SetAttributeRaw("name", tokens("var.k8s_cluster_id"))
	eksAuthBody := rootBody.AppendNewBlock("data", []string{"aws_eks_cluster_auth", "cluster"}).Body()
	eksAuthBody.SetAttributeRaw("name", tokens("var.k8s_cluster_id"))

	kubernetesProviderBody := rootBody.AppendNewBlock("provider", []string{"kubernetes"}).Body()
	kubernetesProviderBody.SetAttributeRaw("host", tokens("data.aws_eks_cluster.cluster.endpoint"))
	kubernetesProviderBody.SetAttributeRaw("cluster_ca_certificate", tokens("base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)"))
	kubernetesProviderBody.SetAttributeRaw("token", tokens("data.aws_eks_cluster_auth.cluster.token"))

	kubeNamespaceBody := rootBody.AppendNewBlock("data", []string{"kubernetes_namespace", "happy-namespace"}).Body()
	kubeNamespaceBody.AppendNewBlock("metadata", nil).Body().SetAttributeRaw("name", tokens("var.k8s_namespace"))

	awsAliasTokens := tokens("aws.czi-si")
	appKeyBody := rootBody.AppendNewBlock("data", []string{"aws_ssm_parameter", "dd_app_key"}).Body()
	appKeyBody.SetAttributeValue("name", cty.StringVal("/shared-infra-prod-datadog/app_key"))
	appKeyBody.SetAttributeRaw("provider", awsAliasTokens)
	apiKeyBody := rootBody.AppendNewBlock("data", []string{"aws_ssm_parameter", "dd_api_key"}).Body()
	apiKeyBody.SetAttributeValue("name", cty.StringVal("/shared-infra-prod-datadog/api_key"))
	apiKeyBody.SetAttributeRaw("provider", awsAliasTokens)

	datadogProviderBody := rootBody.AppendNewBlock("provider", []string{"datadog"}).Body()
	datadogProviderBody.SetAttributeRaw("app_key", tokens("data.aws_ssm_parameter.dd_app_key.value"))
	datadogProviderBody.SetAttributeRaw("api_key", tokens("data.aws_ssm_parameter.dd_api_key.value"))

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func (tf TfGenerator) generateAwsProvider(rootBody *hclwrite.Body, alias, accountIdExpr, roleExpr string) error {
	awsProviderBody := rootBody.AppendNewBlock("provider", []string{"aws"}).Body()
	if alias != "" {
		awsProviderBody.SetAttributeValue("alias", cty.StringVal(alias))
	}
	awsProviderBody.SetAttributeValue("region", cty.StringVal(*tf.happyConfig.AwsRegion()))

	assumeRoleBlockBody := awsProviderBody.AppendNewBlock("assume_role", nil).Body()
	assumeRoleBlockBody.SetAttributeRaw("role_arn", tokens(fmt.Sprintf("\"arn:aws:iam::%s:role/%s\"", accountIdExpr, roleExpr)))
	awsProviderBody.SetAttributeRaw("allowed_account_ids", tokens(fmt.Sprintf("[\"%s\"]", accountIdExpr)))
	return nil
}

func (tf *TfGenerator) GenerateVersions(srcDir string) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "versions.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()
	terraformBlockBody := rootBody.AppendNewBlock("terraform", nil).Body()
	terraformBlockBody.SetAttributeValue("required_version", cty.StringVal(requiredTerraformVersion))
	requiredProvidersBody := terraformBlockBody.AppendNewBlock("required_providers", nil).Body()

	for _, provider := range requiredProviders {
		p := cty.ObjectVal(map[string]cty.Value{
			"source":  cty.StringVal(provider.Source),
			"version": cty.StringVal(provider.Version),
		})
		requiredProvidersBody.SetAttributeValue(provider.Name, p)

	}

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func (tf *TfGenerator) GenerateOutputs(srcDir string, outputs []Output) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "outputs.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()

	sort.SliceStable(outputs, func(i, j int) bool {
		return strings.Compare(outputs[i].Name, outputs[j].Name) < 0
	})

	for _, output := range outputs {
		moduleOutputBody := rootBody.AppendNewBlock("output", []string{output.Name}).Body()
		if len(output.Description) > 0 {
			moduleOutputBody.SetAttributeValue("description", cty.StringVal(output.Description))
		}
		moduleOutputBody.SetAttributeValue("sensitive", cty.BoolVal(output.Sensitive))
		tokens := hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{Name: "module"},
			hcl.TraverseAttr{Name: "stack"},
			hcl.TraverseAttr{Name: output.Name},
		})
		moduleOutputBody.SetAttributeRaw("value", tokens)
	}

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func (tf *TfGenerator) GenerateVariables(srcDir string) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "variables.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()
	for _, variable := range requiredVariables {
		variableBody := rootBody.AppendNewBlock("variable", []string{variable.Name}).Body()
		tokens := hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{Name: variable.Type},
		})
		variableBody.SetAttributeRaw("type", tokens)
		variableBody.SetAttributeValue("description", cty.StringVal(variable.Description))
		if !variable.Default.IsNull() {
			variableBody.SetAttributeValue("default", variable.Default)
		}
	}

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func ParseModuleSource(moduleSource string) (gitUrl string, modulePath string, ref string, err error) {

	parts := strings.Split(moduleSource, "//")
	if len(parts) < 2 {
		return "", "", "", errors.Errorf("invalid module source %s", moduleSource)
	}

	gitUrl = parts[0]
	modulePathAndRef := parts[1]

	modulePathAndRefParts := strings.Split(modulePathAndRef, "?ref=")

	if len(modulePathAndRefParts) < 2 {
		return "", "", "", errors.Errorf("invalid module source, reference is missing %s", moduleSource)
	}

	modulePath = modulePathAndRefParts[0]
	ref = modulePathAndRefParts[1]

	return gitUrl, modulePath, ref, nil
}

func stringToTokens(value string) (hclwrite.Tokens, error) {
	file, diags := hclwrite.ParseConfig([]byte("attr = "+value), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags.Errs()[0]
	}
	attr := file.Body().GetAttribute("attr")
	return attr.Expr().BuildTokens(hclwrite.Tokens{}), nil
}

func tokens(value string) hclwrite.Tokens {
	tokens, err := stringToTokens(value)
	if err != nil {
		logrus.Errorf("Unable to parse an HCL expression: %s: %s", value, err.Error())
	}
	return tokens
}
