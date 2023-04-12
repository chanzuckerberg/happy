package model

type AppStack struct {
	CommonDBFields
	AppMetadata // TODO: might want to change this to AppStackPayload but going with minimal columns for now
}

type AppStackResponse struct {
	AppMetadata
	Endpoints       map[string]string `json:"endpoints,omitempty"`
	WorkspaceUrl    string            `json:"workspace_url,omitempty"`
	WorkspaceStatus string            `json:"workspace_status,omitempty"`
	WorkspaceRunUrl string            `json:"workspace_run_url,omitempty"`
	Error           string            `json:"error,omitempty"`
} // @Name response.AppStackResponse

type AppStackPayload struct {
	AppMetadata
	AwsProfile     string `query:"aws_profile"      validate:"required"`
	AwsRegion      string `query:"aws_region"       validate:"required"`
	TaskLaunchType string `query:"task_launch_type" validate:"required,oneof=fargate k8s"`
	K8SNamespace   string `query:"k8s_namespace"    validate:"required_if=TaskLaunchType k8s"`
	K8SClusterId   string `query:"k8s_cluster_id"   validate:"required_if=TaskLaunchType k8s"`
	SecretId       string `query:"secret_id"        validate:"required_if=TaskLaunchType fargate"`
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

func MakeAppStackPayload(appName, env, stack, awsProfile, awsRegion, launchType, k8sNamespace, k8sClusterId string) AppStackPayload {
	return AppStackPayload{
		AppMetadata:    *NewAppMetadata(appName, env, stack),
		AwsProfile:     awsProfile,
		AwsRegion:      awsRegion,
		TaskLaunchType: launchType,
		K8SNamespace:   k8sNamespace,
		K8SClusterId:   k8sClusterId,
	}
}
