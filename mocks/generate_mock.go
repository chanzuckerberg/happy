package mocks

// Add mocks as necessary
//go:generate mockgen -destination=mock_workspace.go -package=mocks github.com/chanzuckerberg/happy/pkg/workspace_repo Workspace
//go:generate mockgen -destination=mock_workspace_repo.go -package=mocks github.com/chanzuckerberg/happy/pkg/workspace_repo WorkspaceRepoIface
//go:generate mockgen -destination=mock_dir_processor.go -package=mocks github.com/chanzuckerberg/happy/pkg/util DirProcessor
//go:generate mockgen -destination=mock_stack.go -package=mocks github.com/chanzuckerberg/happy/pkg/stack_mgr StackIface
//go:generate mockgen -destination=mock_stack_service.go -package=mocks github.com/chanzuckerberg/happy/pkg/stack_mgr StackServiceIface

// AWS

//go:generate mockgen -destination=mock_aws_ec2.go -package=mocks github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API
//go:generate mockgen -destination=mock_aws_secretsmanager.go -package=mocks github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface SecretsManagerAPI
//go:generate mockgen -destination=mock_aws_ssm.go -package=mocks github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI
//go:generate mockgen -destination=mock_aws_sts.go -package=mocks github.com/aws/aws-sdk-go/service/sts/stsiface STSAPI
//go:generate mockgen -destination=mock_aws_ecr.go -package=mocks github.com/aws/aws-sdk-go/service/ecr/ecriface ECRAPI
//go:generate mockgen -destination=mock_aws_ecs.go -package=mocks github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI
//go:generate mockgen -destination=mock_aws_logs.go -package=mocks github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface CloudWatchLogsAPI
