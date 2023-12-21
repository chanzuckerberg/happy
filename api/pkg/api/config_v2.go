package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	"github.com/chanzuckerberg/happy/api/pkg/request"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/shared/backend/aws"
	"github.com/chanzuckerberg/happy/shared/model"
	"github.com/pkg/errors"
)

func getSecretName(appName, environment, stack string) string {
	appName = strings.ToLower(appName)

	// replace all non-alphanumeric characters with hyphens (-)
	regex := regexp.MustCompile("[^a-z0-9]")
	appName = regex.ReplaceAllString(appName, "-")
	stack = regex.ReplaceAllString(stack, "-")

	var parts []string
	for _, s := range []string{"happy-config", appName, environment, stack} {
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ".")
}

// func getK8sBackend() *aws.K8SComputeBackend {
// 	awsCtx := model.AWSContext{
// 		AWSProfile:     params.AWSProfile,
// 		AWSRegion:      params.AWSRegion,
// 		TaskLaunchType: "k8s",
// 		K8SNamespace:   params.K8sNamespace,
// 		K8SClusterID:   params.K8sClusterID,
// 	}
// 	ctx, err := request.AddAWSAuthToCtx(ctx, params.XAWSAccessKeyID, params.XAWSSecretAccessKey, params.XAWSSessionToken)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "adding aws auth to ctx")
// 	}

// 	happyClient, err := request.MakeHappyClient(ctx, params.AppName, awsCtx.MakeEnvironmentContext(params.Environment))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "making happy client")
// 	}

// 	cb, err := happyClient.AWSBackend.GetComputeBackend(ctx)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "getting compute backend")
// 	}

//		return cb.(*aws.K8SComputeBackend)
//	}
func (h handler) SetAppConfig(ctx context.Context, params ogent.SetAppConfigParams) (ogent.SetAppConfigRes, error) {
	awsCtx := model.AWSContext{
		AWSProfile:     params.AWSProfile,
		AWSRegion:      params.AWSRegion,
		TaskLaunchType: "k8s",
		K8SNamespace:   params.K8sNamespace,
		K8SClusterID:   params.K8sClusterID,
	}
	ctx, err := request.AddAWSAuthToCtx(ctx, params.XAWSAccessKeyID, params.XAWSSecretAccessKey, params.XAWSSessionToken)
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Parsing AWS auth headers").Error())
	}

	happyClient, err := request.MakeHappyClient(ctx, params.AppName, awsCtx.MakeEnvironmentContext(params.Environment))
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Making happy client").Error())
	}

	cb, err := happyClient.AWSBackend.GetComputeBackend(ctx)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrap(err, "Getting compute backend").Error())
	}

	k8sBackend := cb.(*aws.K8SComputeBackend)
	secretName := getSecretName(params.AppName, params.Environment, params.Stack.Or(""))

	res, err := k8sBackend.WriteSecret(ctx, secretName, params.Key, params.Value)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrapf(err, "Writing [%s] to secrets", params.Key).Error())
	}

	source := ogent.AppConfigListSourceStack
	stack := params.Stack.Or("")
	if stack == "" {
		source = ogent.AppConfigListSourceEnvironment
	}
	return &ogent.AppConfigList{
		AppName:     params.AppName,
		Environment: params.Environment,
		Stack:       params.Stack.Or(""),
		Source:      source,
		Key:         params.Key,
		Value:       string(res[params.Key]),
	}, nil
}

func (h handler) DeleteAppConfig(ctx context.Context, params ogent.DeleteAppConfigParams) (ogent.DeleteAppConfigRes, error) {
	awsCtx := model.AWSContext{
		AWSProfile:     params.AWSProfile,
		AWSRegion:      params.AWSRegion,
		TaskLaunchType: "k8s",
		K8SNamespace:   params.K8sNamespace,
		K8SClusterID:   params.K8sClusterID,
	}
	ctx, err := request.AddAWSAuthToCtx(ctx, params.XAWSAccessKeyID, params.XAWSSecretAccessKey, params.XAWSSessionToken)
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Parsing AWS auth headers").Error())
	}

	happyClient, err := request.MakeHappyClient(ctx, params.AppName, awsCtx.MakeEnvironmentContext(params.Environment))
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Making happy client").Error())
	}

	cb, err := happyClient.AWSBackend.GetComputeBackend(ctx)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrap(err, "Getting compute backend").Error())
	}

	k8sBackend := cb.(*aws.K8SComputeBackend)
	secretName := getSecretName(params.AppName, params.Environment, params.Stack.Or(""))

	// delete the key
	_, err = k8sBackend.WriteSecret(ctx, secretName, params.Key, "")
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrapf(err, "Deleting [%s] from secrets", params.Key).Error())
	}

	return &ogent.DeleteAppConfigOK{}, nil
}

func (h handler) ListAppConfig(ctx context.Context, params ogent.ListAppConfigParams) (ogent.ListAppConfigRes, error) {
	awsCtx := model.AWSContext{
		AWSProfile:     params.AWSProfile,
		AWSRegion:      params.AWSRegion,
		TaskLaunchType: "k8s",
		K8SNamespace:   params.K8sNamespace,
		K8SClusterID:   params.K8sClusterID,
	}
	ctx, err := request.AddAWSAuthToCtx(ctx, params.XAWSAccessKeyID, params.XAWSSecretAccessKey, params.XAWSSessionToken)
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Parsing AWS auth headers").Error())
	}

	happyClient, err := request.MakeHappyClient(ctx, params.AppName, awsCtx.MakeEnvironmentContext(params.Environment))
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Making happy client").Error())
	}

	cb, err := happyClient.AWSBackend.GetComputeBackend(ctx)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrap(err, "Getting compute backend").Error())
	}

	k8sBackend := cb.(*aws.K8SComputeBackend)

	envSecretName := getSecretName(params.AppName, params.Environment, "")
	envSecrets, err := k8sBackend.GetSecret(ctx, envSecretName)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrapf(err, "Getting env secret [%s]", envSecretName).Error())
	}

	results := make(map[string]struct {
		source ogent.AppConfigListSource
		value  []byte
	})
	for key, value := range envSecrets {
		if len(value) == 0 {
			continue
		}
		results[key] = struct {
			source ogent.AppConfigListSource
			value  []byte
		}{
			source: ogent.AppConfigListSourceEnvironment,
			value:  value,
		}
	}

	stack := params.Stack.Or("")
	if stack != "" {
		stackSecretName := getSecretName(params.AppName, params.Environment, stack)
		stackSecrets, err := k8sBackend.GetSecret(ctx, stackSecretName)
		if err != nil {
			return nil, response.NewBadRequestError(errors.Wrapf(err, "Getting stack secret [%s]", stackSecretName).Error())
		}
		for key, value := range stackSecrets {
			if len(value) == 0 {
				continue
			}
			results[key] = struct {
				source ogent.AppConfigListSource
				value  []byte
			}{
				source: ogent.AppConfigListSourceStack,
				value:  value,
			}
		}
	}

	var configs []ogent.AppConfigList
	for key, secret := range results {
		configs = append(configs, ogent.AppConfigList{
			AppName:     params.AppName,
			Environment: params.Environment,
			Stack:       stack,
			Source:      secret.source,
			Key:         key,
			Value:       string(secret.value),
		})
	}

	return (*ogent.ListAppConfigOKApplicationJSON)(&configs), nil
}

func (h handler) ReadAppConfig(ctx context.Context, params ogent.ReadAppConfigParams) (ogent.ReadAppConfigRes, error) {
	awsCtx := model.AWSContext{
		AWSProfile:     params.AWSProfile,
		AWSRegion:      params.AWSRegion,
		TaskLaunchType: "k8s",
		K8SNamespace:   params.K8sNamespace,
		K8SClusterID:   params.K8sClusterID,
	}
	ctx, err := request.AddAWSAuthToCtx(ctx, params.XAWSAccessKeyID, params.XAWSSecretAccessKey, params.XAWSSessionToken)
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Parsing AWS auth headers").Error())
	}

	happyClient, err := request.MakeHappyClient(ctx, params.AppName, awsCtx.MakeEnvironmentContext(params.Environment))
	if err != nil {
		return nil, response.NewForbiddenError(errors.Wrap(err, "Making happy client").Error())
	}

	cb, err := happyClient.AWSBackend.GetComputeBackend(ctx)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrap(err, "Getting compute backend").Error())
	}

	k8sBackend := cb.(*aws.K8SComputeBackend)

	envSecretName := getSecretName(params.AppName, params.Environment, "")
	envSecrets, err := k8sBackend.GetSecret(ctx, envSecretName)
	if err != nil {
		return nil, response.NewBadRequestError(errors.Wrapf(err, "Getting env secret [%s]", envSecretName).Error())
	}
	result := envSecrets[params.Key]
	source := ogent.AppConfigListSourceEnvironment

	stack := params.Stack.Or("")
	if stack != "" {
		stackSecretName := getSecretName(params.AppName, params.Environment, stack)
		stackSecrets, err := k8sBackend.GetSecret(ctx, stackSecretName)
		if err != nil {
			return nil, response.NewBadRequestError(errors.Wrapf(err, "Getting stack secret [%s]", stackSecretName).Error())
		}
		stackSecretValue, ok := stackSecrets[params.Key]

		if ok && len(stackSecretValue) > 0 {
			result = stackSecretValue
			source = ogent.AppConfigListSourceStack
		}
	}

	if len(result) == 0 {
		return nil, response.NewNotFoundError(fmt.Sprintf("The specified app config was not found: %s", params.Key))
	}

	r := ogent.AppConfigList{
		AppName:     params.AppName,
		Environment: params.Environment,
		Stack:       stack,
		Source:      source,
		Key:         params.Key,
		Value:       string(result),
	}
	return (ogent.ReadAppConfigRes)(&r), nil
}
