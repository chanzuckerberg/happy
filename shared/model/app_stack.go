package model

import (
	"strings"

	"github.com/chanzuckerberg/happy/shared/config"
	"github.com/chanzuckerberg/happy/shared/k8s"
	"github.com/chanzuckerberg/happy/shared/util"
)

type AppStack struct {
	CommonDBFields
	AppMetadata // TODO: might want to change this to AppStackPayload but going with minimal columns for now
}

type StackMetadata struct {
	Owner              string            `json:"owner,omitempty"`
	Tag                string            `json:"tag,omitempty"`
	LastUpdated        string            `json:"last_updated,omitempty"`
	Message            string            `json:"message,omitempty"`
	Outputs            map[string]string `json:"outputs,omitempty"`
	Endpoints          map[string]string `json:"endpoints,omitempty"`
	GitRepo            string            `json:"git_repo,omitempty"`
	GitSHA             string            `json:"git_sha,omitempty"`
	GitBranch          string            `json:"git_branch,omitempty"`
	TFEWorkspaceURL    string            `json:"tfe_workspace_url,omitempty"`
	TFEWorkspaceStatus string            `json:"tfe_workspace_status,omitempty"`
	TFEWorkspaceRunURL string            `json:"tfe_workspace_run_url,omitempty"`
}

type AWSContext struct {
	AWSProfile     string `query:"aws_profile"`
	AWSRegion      string `query:"aws_region"       validate:"required"`
	TaskLaunchType string `query:"task_launch_type" validate:"required,oneof=fargate k8s"`
	K8SNamespace   string `query:"k8s_namespace"    validate:"required_if=TaskLaunchType k8s"`
	K8SClusterID   string `query:"k8s_cluster_id"   validate:"required_if=TaskLaunchType k8s"`
	SecretID       string `query:"secret_id"        validate:"required_if=TaskLaunchType fargate"`
}

func (a *AWSContext) MakeEnvironmentContext(env string) config.EnvironmentContext {
	return config.EnvironmentContext{
		EnvironmentName: env,
		AWSProfile:      &a.AWSProfile,
		AWSRegion:       &a.AWSRegion,
		SecretID:        a.SecretID,
		TaskLaunchType:  util.LaunchType(strings.ToUpper(a.TaskLaunchType)),
		K8S: k8s.K8SConfig{
			Namespace: a.K8SNamespace,
			ClusterID: a.K8SClusterID,
			// we only want to use EKS auth method on the API side
			// since we won't have any local files to work with
			AuthMethod: "eks",
		},
	}
}

type AppStackResponse struct {
	AppMetadata
	StackMetadata
	Error string `json:"error,omitempty"`
} // @Name response.AppStackResponse

type AppStackPayload struct {
	AppMetadata
	AWSContext
	ListAll bool `query:"all"`
} // @Name payload.AppStackPayload

type WrappedAppStacksWithCount struct {
	Records []*AppStackResponse `json:"records"`
	Count   int                 `json:"count" example:"1"`
} // @Name response.WrappedAppStacksWithCount

type WrappedAppStack struct {
	Record *AppStack `json:"record"`
} // @Name response.WrappedAppStack

func MakeAppStackResponse(appName, env, stack string) AppStackResponse {
	return AppStackResponse{
		AppMetadata: *NewAppMetadata(appName, env, stack),
	}
}

func NewAppStackFromAppStackPayload(payload AppStackPayload) *AppStack {
	return &AppStack{
		AppMetadata: *NewAppMetadata(payload.AppName, payload.Environment, payload.Stack),
	}
}

func MakeAppStackPayload(appName, env, stack string, awsContext AWSContext) AppStackPayload {
	meta := NewAppMetadata(appName, env, stack)
	return AppStackPayload{
		AppMetadata: *meta,
		AWSContext:  awsContext,
	}
}
