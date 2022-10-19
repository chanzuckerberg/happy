package main

// AWS

//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ec2.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces EC2API
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_secretsmanager.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces SecretsManagerAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ssm.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces SSMAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_sts.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces STSAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_sts_presign.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces STSPresignAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ecr.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces ECRAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ecs.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces ECSAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_eks.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces EKSAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ecs_task_stopped_waiter.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces ECSTaskStoppedWaiterAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_get_logs.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces GetLogEventsAPIClient
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_filter_logs.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces FilterLogEventsAPIClient
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_dynamodb.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces DynamoDB
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_compute_backend.go -package=interfaces github.com/chanzuckerberg/happy/pkg/cli/backend/aws/interfaces ComputeBackend
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_kubernetes.go -package=interfaces -mock_names Interface=MockKubernetesAPI k8s.io/client-go/kubernetes Interface
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_kubernetes_corev1.go -package=interfaces -mock_names CoreV1Interface=MockKubernetesCoreV1API,SecretInterface=MockKubernetesSecretAPI k8s.io/client-go/kubernetes/typed/core/v1 CoreV1Interface,SecretInterface
