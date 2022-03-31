package testbackend

// AWS

//go:generate mockgen -destination=mock_aws_ec2.go -package=testbackend github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface EC2API
//go:generate mockgen -destination=mock_aws_secretsmanager.go -package=testbackend github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface SecretsManagerAPI
//go:generate mockgen -destination=mock_aws_ssm.go -package=testbackend github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface SSMAPI
//go:generate mockgen -destination=mock_aws_sts.go -package=testbackend github.com/aws/aws-sdk-go-v2/service/sts/stsiface STSAPI
//go:generate mockgen -destination=mock_aws_ecr.go -package=testbackend github.com/aws/aws-sdk-go-v2/service/ecr/ecriface ECRAPI
//go:generate mockgen -destination=mock_aws_ecs.go -package=testbackend github.com/aws/aws-sdk-go-v2/service/ecs/ecsiface ECSAPI
//go:generate mockgen -destination=mock_aws_logs.go -package=testbackend github.com/aws/aws-sdk-go-v2-v2/service/cloudwatchlogs GetLogEventsAPIClient
