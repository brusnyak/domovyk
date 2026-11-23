package discover

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/brusnyak/domovyk/internal/inventory"
)

// defaultTimeout bounds each AWS service scan so a hung account can't block forever.
const defaultTimeout = 2 * time.Minute

// Options controls discovery behaviour.
// Credentials come from the AWS SDK default credential chain
// (env, shared config ~/.aws, assume-role, IMDS). Named profile is optional.
// Secret keys must never be passed on the command line.
type Options struct {
	Region  string
	Profile string
}

// Discover reads a read-only inventory of AWS resources.
// ctx must be non-nil; it is used for every SDK call and bounds the whole scan.
func Discover(ctx context.Context, opts Options) (*inventory.Inventory, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(opts.Region),
	}
	if opts.Profile != "" {
		loadOpts = append(loadOpts, config.WithSharedConfigProfile(opts.Profile))
	}
	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	inv := inventory.NewInventory(opts.Region)

	// Account identity (best effort; a failed STS call does not abort discovery).
	accountID, _ := accountIdentity(ctx, sts.NewFromConfig(cfg))
	inv.Manifest.AccountID = accountID
	inv.Manifest.Regions = []string{opts.Region}

	ec2Client := ec2.NewFromConfig(cfg)
	rdsClient := rds.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	eksClient := eks.NewFromConfig(cfg)
	lambdaClient := lambda.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)

	// EC2 instances (paginated)
	scanEC2Instances(ctx, ec2Client, opts.Region, inv)
	// EBS volumes (paginated)
	scanEBSVolumes(ctx, ec2Client, opts.Region, inv)
	// RDS instances (paginated)
	scanRDSInstances(ctx, rdsClient, opts.Region, inv)
	// S3 buckets (not paginated)
	scanS3Buckets(ctx, s3Client, opts.Region, inv)
	// EKS clusters (paginated list, then describe each)
	scanEKSClusters(ctx, eksClient, opts.Region, inv)
	// Lambda functions (paginated)
	scanLambdaFunctions(ctx, lambdaClient, opts.Region, inv)
	// IAM roles (paginated)
	scanIAMRoles(ctx, iamClient, opts.Region, inv)

	inv.CountResources()
	return inv, nil
}
