package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type ECRAPI interface {
	BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
	PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)
	GetAuthorizationToken(ctx context.Context, params *ecr.GetAuthorizationTokenInput, optFns ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)
	DescribeImageScanFindings(context.Context, *ecr.DescribeImageScanFindingsInput, ...func(*ecr.Options)) (*ecr.DescribeImageScanFindingsOutput, error)
	BatchGetRepositoryScanningConfiguration(ctx context.Context, params *ecr.BatchGetRepositoryScanningConfigurationInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetRepositoryScanningConfigurationOutput, error)
}
