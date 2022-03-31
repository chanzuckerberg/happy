package main

// AWS

//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ec2.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces EC2API
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_secretsmanager.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces SecretsManagerAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ssm.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces SSMAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_sts.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces STSAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ecr.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces ECRAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_ecs.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces ECSAPI
//go:generate mockgen -destination=./pkg/backend/aws/interfaces/mock_aws_logs.go -package=interfaces github.com/chanzuckerberg/happy/pkg/backend/aws/interfaces GetLogEventsAPIClient
