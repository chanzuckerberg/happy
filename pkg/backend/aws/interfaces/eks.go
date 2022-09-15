package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type EKSAPI interface {
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
}
