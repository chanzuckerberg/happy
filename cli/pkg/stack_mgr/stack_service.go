package stack_mgr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	backend "github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/diagnostics"
	"github.com/chanzuckerberg/happy/shared/options"
	"github.com/chanzuckerberg/happy/shared/util"
	"github.com/chanzuckerberg/happy/shared/util/tf"
	workspacerepo "github.com/chanzuckerberg/happy/shared/workspace_repo"
	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
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

type StackServiceIface interface {
	NewStackMeta(stackName string) *StackMeta
	Add(ctx context.Context, stackName string, options ...workspacerepo.TFERunOption) (*Stack, error)
	Remove(ctx context.Context, stackName string, options ...workspacerepo.TFERunOption) error
	GetStacks(ctx context.Context) (map[string]*Stack, error)
	GetStackWorkspace(ctx context.Context, stackName string) (workspacerepo.Workspace, error)
	GetConfig() *config.HappyConfig
}

type StackService struct {
	// dependencies
	backend       *backend.Backend
	workspaceRepo workspacerepo.WorkspaceRepoIface
	dirProcessor  util.DirProcessor
	executor      util.Executor
	happyConfig   *config.HappyConfig

	// NOTE: creator Workspace is a workspace that creates dependent workspaces with
	// given default values and configuration
	// the derived workspace is then used to launch the actual happy infrastructure
	creatorWorkspaceName string
}

func NewStackService() *StackService {
	return &StackService{
		dirProcessor: util.NewLocalProcessor(),
		executor:     util.NewDefaultExecutor(),
	}
}

func (s *StackService) GetWritePath() string {
	return fmt.Sprintf("/happy/%s/stacklist", s.happyConfig.GetEnv())
}

func (s *StackService) GetNamespacedWritePath() string {
	return fmt.Sprintf("/happy/%s/%s/stacklist", s.happyConfig.App(), s.happyConfig.GetEnv())
}

func (s *StackService) WithBackend(backend *backend.Backend) *StackService {
	creatorWorkspaceName := fmt.Sprintf("env-%s", s.happyConfig.GetEnv())

	s.creatorWorkspaceName = creatorWorkspaceName
	s.backend = backend

	return s
}

func (s *StackService) WithHappyConfig(happyConfig *config.HappyConfig) *StackService {
	s.happyConfig = happyConfig
	return s
}

func (s *StackService) WithExecutor(executor util.Executor) *StackService {
	s.executor = executor
	return s
}

func (s *StackService) WithWorkspaceRepo(workspaceRepo workspacerepo.WorkspaceRepoIface) *StackService {
	s.workspaceRepo = workspaceRepo
	return s
}

func (s *StackService) GetConfig() *config.HappyConfig {
	return s.happyConfig
}

// Invoke a specific TFE workspace that creates/deletes TFE workspaces,
// with prepopulated variables for identifier tokens.
func (s *StackService) resync(ctx context.Context, wait bool, options ...workspacerepo.TFERunOption) error {
	log.Debug("resyncing new workspace...")
	log.Debugf("running creator workspace %s...", s.creatorWorkspaceName)
	creatorWorkspace, err := s.workspaceRepo.GetWorkspace(ctx, s.creatorWorkspaceName)
	if err != nil {
		return errors.Wrapf(err, "unable to get workspace %s", s.creatorWorkspaceName)
	}
	err = creatorWorkspace.Run(ctx, options...)
	if err != nil {
		return errors.Wrapf(err, "error running latest %s workspace version", s.creatorWorkspaceName)
	}
	if wait {
		return creatorWorkspace.Wait(ctx)
	}
	return nil
}

func (s *StackService) GetLatestDeployedTag(ctx context.Context, stackName string) (string, error) {
	stack, err := s.GetStack(ctx, stackName)
	if err != nil {
		return "", errors.Wrap(err, "unable to get the stack")
	}
	stackInfo, err := stack.GetStackInfo(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get the stack info")
	}
	return stackInfo.Tag, nil
}

func (s *StackService) Remove(ctx context.Context, stackName string, opts ...workspacerepo.TFERunOption) error {
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		return nil
	}
	var err error
	if s.GetConfig().GetFeatures().EnableDynamoLocking {
		err = s.removeFromStacklistWithLock(ctx, stackName)
	} else {
		err = s.removeFromStacklist(ctx, stackName)
	}
	if err != nil {
		return err
	}

	return s.resync(ctx, false, opts...)
}

func (s *StackService) removeFromStacklistWithLock(ctx context.Context, stackName string) error {
	distributedLock, err := s.getDistributedLock()
	if err != nil {
		return err
	}
	defer distributedLock.Close(ctx)

	lockKey := s.GetNamespacedWritePath()
	lock, err := distributedLock.AcquireLock(ctx, lockKey)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.removeFromStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", lockKey))
	}

	return ret
}

func (s *StackService) removeFromStacklist(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Removing stack...")

	stacks, err := s.GetStacks(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get a list of stacks")
	}
	stackNamesList := []string{}
	for name := range stacks {
		if name != stackName {
			stackNamesList = append(stackNamesList, name)
		}
	}

	return s.writeStacklist(ctx, stackNamesList)
}

func (s *StackService) Add(ctx context.Context, stackName string, opts ...workspacerepo.TFERunOption) (*Stack, error) {
	log.WithField("stack_name", stackName).Debug("Adding a new stack...")
	dryRun, ok := ctx.Value(options.DryRunKey).(bool)
	if !ok {
		dryRun = false
	}
	if dryRun {
		log.Debugf("temporarily creating a TFE workspace for stack '%s'", stackName)
	} else {
		log.Debugf("creating stack '%s'", stackName)
	}

	var err error
	if s.GetConfig().GetFeatures().EnableDynamoLocking {
		err = s.addToStacklistWithLock(ctx, stackName)
	} else {
		err = s.addToStacklist(ctx, stackName)
	}
	if err != nil {
		return nil, err
	}

	if !util.IsLocalstackMode() {
		// Create the workspace
		wait := true
		if err := s.resync(ctx, wait, opts...); err != nil {
			return nil, err
		}
	}

	_, err = s.GetStackWorkspace(ctx, stackName)
	if err != nil {
		return nil, err
	}
	return s.createStack(stackName), nil
}

func (s *StackService) addToStacklistWithLock(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Adding new stack with a lock...")
	distributedLock, err := s.getDistributedLock()
	if err != nil {
		return err
	}
	defer distributedLock.Close(ctx)

	lockKey := s.GetNamespacedWritePath()
	lock, err := distributedLock.AcquireLock(ctx, lockKey)
	if err != nil {
		return err
	}

	// don't return if there was an error here, we still need to release the lock so we'll use multierror instead
	ret := s.addToStacklist(ctx, stackName)

	_, err = distributedLock.ReleaseLock(ctx, lock)
	if err != nil {
		ret = multierror.Append(ret, errors.Wrapf(err, "unable to release the lock on %s", lockKey))
	}

	return ret
}

func (s *StackService) addToStacklist(ctx context.Context, stackName string) error {
	log.WithField("stack_name", stackName).Debug("Adding new stack...")
	existStacks, err := s.GetStacks(ctx)
	if err != nil {
		return err
	}

	newStackNames := []string{}
	stackNameExists := false
	for name := range existStacks {
		newStackNames = append(newStackNames, name)
		if name == stackName {
			stackNameExists = true
		}
	}
	if !stackNameExists {
		newStackNames = append(newStackNames, stackName)
	}

	return s.writeStacklist(ctx, newStackNames)
}

func (s *StackService) writeStacklist(ctx context.Context, stackNames []string) error {
	sort.Strings(stackNames)

	stackNamesJson, err := json.Marshal(stackNames)
	if err != nil {
		return errors.Wrap(err, "unable to serialize stack list as json")
	}

	stackNamesStr := string(stackNamesJson)
	log.WithFields(log.Fields{"path": s.GetNamespacedWritePath(), "data": stackNamesStr}).Debug("Writing to paramstore...")
	if err := s.backend.ComputeBackend.WriteParam(ctx, s.GetNamespacedWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}
	log.WithFields(log.Fields{"path": s.GetWritePath(), "data": stackNamesStr}).Debug("Writing to paramstore...")
	if err := s.backend.ComputeBackend.WriteParam(ctx, s.GetWritePath(), stackNamesStr); err != nil {
		return errors.Wrap(err, "unable to write a workspace param")
	}

	return nil
}

func (s *StackService) GetStacks(ctx context.Context) (map[string]*Stack, error) {
	defer diagnostics.AddProfilerRuntime(ctx, time.Now(), "GetStacks")
	log.WithField("path", s.GetNamespacedWritePath()).Debug("Reading stacks from paramstore at path...")
	paramOutput, err := s.backend.ComputeBackend.GetParam(ctx, s.GetNamespacedWritePath())
	if err != nil && strings.Contains(err.Error(), "ParameterNotFound") {
		log.WithField("path", s.GetWritePath()).Debug("Reading stacks from paramstore at path...")
		paramOutput, err = s.backend.ComputeBackend.GetParam(ctx, s.GetWritePath())
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stacks")
	}

	log.WithField("output", paramOutput).Debug("read stacks info from param store")

	var stacklist []string
	err = json.Unmarshal([]byte(paramOutput), &stacklist)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json")
	}

	log.WithField("output", stacklist).Debug("marshalled json output to string slice")

	stacks := map[string]*Stack{}
	for _, stackName := range stacklist {
		stacks[stackName] = s.createStack(stackName)
	}

	return stacks, nil
}

func (s *StackService) CollectStackInfo(ctx context.Context, listAll bool, app string) ([]StackInfo, error) {
	g, ctx := errgroup.WithContext(ctx)
	stacks, err := s.GetStacks(ctx)
	if err != nil {
		return nil, err
	}
	// Iterate in order
	stackNames := maps.Keys(stacks)
	stackInfos := make([]*StackInfo, len(stackNames))
	sort.Strings(stackNames)
	for i, name := range stackNames {
		i, name := i, name // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			stackInfo, err := stacks[name].GetStackInfo(ctx)
			if err != nil {
				log.Warnf("unable to get stack info for %s: %s (likely means the deploy failed the first time)", name, err)
				if !diagnostics.IsInteractiveContext(ctx) {
					stackInfos[i] = &StackInfo{
						Name:    name,
						Status:  "error",
						Message: err.Error(),
					}
				}
				// we still want to show the other stacks if this errors
				return nil
			}

			// only show the stacks that belong to this app or they want to list all
			if listAll || (stackInfo != nil && stackInfo.App == app) {
				stackInfos[i] = stackInfo
			}

			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stack infos")
	}

	// remove empties
	nonEmptyStackInfos := []StackInfo{}
	for _, stackInfo := range stackInfos {
		if stackInfo == nil {
			continue
		}
		nonEmptyStackInfos = append(nonEmptyStackInfos, *stackInfo)
	}
	return nonEmptyStackInfos, g.Wait()
}

func (s *StackService) GetStack(ctx context.Context, stackName string) (*Stack, error) {
	existingStacks, err := s.GetStacks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stacks")
	}
	stack, ok := existingStacks[stackName]
	if !ok {
		return nil, errors.Errorf("stack %s doesn't exist", stackName)
	}

	return stack, nil
}

// pre-format stack name and call workspaceRepo's GetWorkspace method
func (s *StackService) GetStackWorkspace(ctx context.Context, stackName string) (workspacerepo.Workspace, error) {
	workspaceName := fmt.Sprintf("%s-%s", s.happyConfig.GetEnv(), stackName)

	ws, err := s.workspaceRepo.GetWorkspace(ctx, workspaceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workspace")
	}

	return ws, nil
}

func (s *StackService) createStack(stackName string) *Stack {
	return &Stack{
		stackService: s,
		Name:         stackName,
		dirProcessor: s.dirProcessor,
		executor:     s.executor,
	}
}

func (s *StackService) HasState(ctx context.Context, stackName string) (bool, error) {
	workspace, err := s.GetStackWorkspace(ctx, stackName)
	if err != nil {
		if errors.Is(err, tfe.ErrInvalidWorkspaceValue) || errors.Is(err, tfe.ErrResourceNotFound) {
			// Workspace doesn't exist, thus no state
			return false, nil
		}
		return true, errors.Wrap(err, "Cannot get the stack workspace")
	}
	return workspace.HasState(ctx)
}

func (s *StackService) getDistributedLock() (*backend.DistributedLock, error) {
	lockConfig := backend.DistributedLockConfig{DynamodbTableName: s.backend.Conf().GetDynamoLocktableName()}
	return backend.NewDistributedLock(&lockConfig, s.backend.GetDynamoDBClient())
}

func (s *StackService) Generate(ctx context.Context) error {
	moduleSource := "git@github.com:chanzuckerberg/happy//terraform/modules/happy-stack-%s?ref=main"
	if s.GetConfig().TaskLaunchType() == util.LaunchTypeK8S {
		moduleSource = fmt.Sprintf(moduleSource, "eks")
	} else {
		moduleSource = fmt.Sprintf(moduleSource, "ecs")
	}

	_, modulePath, _, err := s.parseModuleSource(moduleSource)
	if err != nil {
		return errors.Wrap(err, "Unable to parse module path out")
	}
	modulePathParts := strings.Split(modulePath, "/")
	moduleName := modulePathParts[len(modulePathParts)-1]

	tempDir, err := os.MkdirTemp("", moduleName)
	if err != nil {
		return errors.Wrap(err, "Unable to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	// Download the module source
	err = getter.GetAny(tempDir, moduleSource)
	if err != nil {
		return errors.Wrap(err, "Unable to download module source")
	}

	// Extract variable information from the module
	variables, err := tf.ParseVariables(tempDir)
	if err != nil {
		return errors.Wrap(err, "Unable to parse out variables from the module")
	}

	outputs, err := tf.ParseOutputs(tempDir)
	if err != nil {
		return errors.Wrap(err, "Unable to parse out variables from the module")
	}

	tfDirPath := s.GetConfig().TerraformDirectory()

	happyProjectRoot := s.GetConfig().GetProjectRoot()
	srcDir := filepath.Join(happyProjectRoot, tfDirPath)

	// Generate main.tf
	err = s.generateMain(srcDir, moduleSource, variables)
	if err != nil {
		return errors.Wrap(err, "Unable to generate main.tf")
	}

	// TODO: Generate variables.tf
	// TODO: Generate outputs.tf
	// TODO: Generate versions.tf
	// TODO: Generate providers.tf

	err = s.generateProviders(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate providers.tf")
	}

	err = s.generateVersions(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate versions.tf")
	}

	err = s.generateOutputs(srcDir, outputs)
	if err != nil {
		return errors.Wrap(err, "Unable to generate outputs.tf")
	}

	err = s.generateVariables(srcDir)
	if err != nil {
		return errors.Wrap(err, "Unable to generate variables.tf")
	}

	return nil
}

func (s *StackService) generateMain(srcDir, moduleSource string, variables []tf.Variable) error {
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
			moduleBlockBody.SetAttributeValue(variable.Name, cty.StringVal(s.happyConfig.App()))
		case "deployment_stage":
			// Set the deployment stage to the current happy environment value
			moduleBlockBody.SetAttributeValue(variable.Name, cty.StringVal(s.happyConfig.GetEnv()))
		case "stack_prefix":
			// Assign stack_prefix variable to "/${var.stack_name}"
			moduleBlockBody.SetAttributeRaw(variable.Name, s.tokens("\"/${var.stack_name}\""))
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
			for _, service := range s.happyConfig.GetServices() {
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

func (s *StackService) generateProviders(srcDir string) error {
	tfFile, err := os.Create(filepath.Join(srcDir, "providers.tf"))
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code")
	}
	defer tfFile.Close()
	hclFile := hclwrite.NewEmptyFile()

	rootBody := hclFile.Body()
	err = s.generateAwsProvider(rootBody, "", "${var.aws_account_id}", "${var.aws_role}")
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code for AWS provider")
	}

	err = s.generateAwsProvider(rootBody, "czi-si", "626314663667", "tfe-si")
	if err != nil {
		return errors.Wrap(err, "Unable to generate HCL code for AWS provider")
	}

	eksBody := rootBody.AppendNewBlock("data", []string{"aws_eks_cluster", "cluster"}).Body()
	eksBody.SetAttributeRaw("name", s.tokens("var.k8s_cluster_id"))
	eksAuthBody := rootBody.AppendNewBlock("data", []string{"aws_eks_cluster_auth", "cluster"}).Body()
	eksAuthBody.SetAttributeRaw("name", s.tokens("var.k8s_cluster_id"))

	kubernetesProviderBody := rootBody.AppendNewBlock("provider", []string{"kubernetes"}).Body()
	kubernetesProviderBody.SetAttributeRaw("host", s.tokens("data.aws_eks_cluster.cluster.endpoint"))
	kubernetesProviderBody.SetAttributeRaw("cluster_ca_certificate", s.tokens("base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)"))
	kubernetesProviderBody.SetAttributeRaw("token", s.tokens("data.aws_eks_cluster_auth.cluster.token"))

	kubeNamespaceBody := rootBody.AppendNewBlock("data", []string{"kubernetes_namespace", "happy-namespace"}).Body()
	kubeNamespaceBody.AppendNewBlock("metadata", nil).Body().SetAttributeRaw("name", s.tokens("var.k8s_namespace"))

	awsAliasTokens := s.tokens("aws.czi-si")
	appKeyBody := rootBody.AppendNewBlock("data", []string{"aws_ssm_parameter", "dd_app_key"}).Body()
	appKeyBody.SetAttributeValue("name", cty.StringVal("/shared-infra-prod-datadog/app_key"))
	appKeyBody.SetAttributeRaw("provider", awsAliasTokens)
	apiKeyBody := rootBody.AppendNewBlock("data", []string{"aws_ssm_parameter", "dd_api_key"}).Body()
	apiKeyBody.SetAttributeValue("name", cty.StringVal("/shared-infra-prod-datadog/api_key"))
	apiKeyBody.SetAttributeRaw("provider", awsAliasTokens)

	datadogProviderBody := rootBody.AppendNewBlock("provider", []string{"datadog"}).Body()
	datadogProviderBody.SetAttributeRaw("app_key", s.tokens("data.aws_ssm_parameter.dd_app_key.value"))
	datadogProviderBody.SetAttributeRaw("api_key", s.tokens("data.aws_ssm_parameter.dd_api_key.value"))

	_, err = tfFile.Write(hclFile.Bytes())

	return err
}

func (s *StackService) generateAwsProvider(rootBody *hclwrite.Body, alias, accountIdExpr, roleExpr string) error {
	awsProviderBody := rootBody.AppendNewBlock("provider", []string{"aws"}).Body()
	if alias != "" {
		awsProviderBody.SetAttributeValue("alias", cty.StringVal(alias))
	}
	awsProviderBody.SetAttributeValue("region", cty.StringVal(*s.happyConfig.AwsRegion()))

	assumeRoleBlockBody := awsProviderBody.AppendNewBlock("assume_role", nil).Body()
	assumeRoleBlockBody.SetAttributeRaw("role_arn", s.tokens(fmt.Sprintf("\"arn:aws:iam::%s:role/%s\"", accountIdExpr, roleExpr)))
	awsProviderBody.SetAttributeRaw("allowed_account_ids", s.tokens(fmt.Sprintf("[\"%s\"]", accountIdExpr)))
	return nil
}

func (s *StackService) generateVersions(srcDir string) error {
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

func (s *StackService) generateOutputs(srcDir string, outputs []tf.Output) error {
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

func (s *StackService) generateVariables(srcDir string) error {
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

func (s *StackService) parseModuleSource(moduleSource string) (gitUrl string, modulePath string, ref string, err error) {

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

func (s *StackService) stringToTokens(value string) (hclwrite.Tokens, error) {
	file, diags := hclwrite.ParseConfig([]byte("attr = "+value), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags.Errs()[0]
	}
	attr := file.Body().GetAttribute("attr")
	return attr.Expr().BuildTokens(hclwrite.Tokens{}), nil
}

func (s *StackService) tokens(value string) hclwrite.Tokens {
	tokens, err := s.stringToTokens(value)
	if err != nil {
		logrus.Errorf("Unable to parse an HCL expression: %s: %s", value, err.Error())
	}
	return tokens
}
