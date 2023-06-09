package main

// AWS

//go:generate mockgen -destination=./aws/interfaces/mock_aws_ec2.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces EC2API
//go:generate mockgen -destination=./aws/interfaces/mock_aws_secretsmanager.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces SecretsManagerAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_ssm.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces SSMAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_sts.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces STSAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_sts_presign.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces STSPresignAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_ecr.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces ECRAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_ecs.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces ECSAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_eks.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces EKSAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_ecs_task_stopped_waiter.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces ECSTaskStoppedWaiterAPI
//go:generate mockgen -destination=./aws/interfaces/mock_aws_get_logs.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces GetLogEventsAPIClient
//go:generate mockgen -destination=./aws/interfaces/mock_aws_filter_logs.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces FilterLogEventsAPIClient
//go:generate mockgen -destination=./aws/interfaces/mock_aws_dynamodb.go -package=interfaces github.com/chanzuckerberg/happy/shared/aws/interfaces DynamoDB
//go:generate mockgen -destination=./aws/interfaces/mock_kubernetes.go -package=interfaces -mock_names Interface=MockKubernetesAPI k8s.io/client-go/kubernetes Interface
//go:generate mockgen -destination=./aws/interfaces/mock_kubernetes_corev1.go -package=interfaces -mock_names CoreV1Interface=MockKubernetesCoreV1API,SecretInterface=MockKubernetesSecretAPI k8s.io/client-go/kubernetes/typed/core/v1 CoreV1Interface,SecretInterface
//go:generate mockgen -destination=./backend/aws/interfaces/mock_compute_backend.go -package=interfaces github.com/chanzuckerberg/happy/shared/backend/aws/interfaces ComputeBackend
//go:generate mockgen -destination=./workspace_repo/mock_workspace.go -package=workspace_repo github.com/chanzuckerberg/happy/shared/workspace_repo Workspace
//go:generate mockgen -destination=./workspace_repo/mock_workspace_repo.go -package=workspace_repo github.com/chanzuckerberg/happy/shared/workspace_repo WorkspaceRepoIface
//go:generate mockgen -destination=./stack/mock_stack_service.go -package=stack github.com/chanzuckerberg/happy/shared/stack StackServiceIface
