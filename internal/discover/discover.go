package discover

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtr(s string) *string {
	return &s
}

// Options controls discovery behaviour.
type Options struct {
	Region    string
	Profile   string
	AccessKey string
	SecretKey string
	Session   string
}

// Discover reads a read-only inventory of AWS resources.
func Discover(opts Options) (*inventory.Inventory, error) {
	ctx := context.Background()

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(opts.Region),
	}
	if opts.AccessKey != "" && opts.SecretKey != "" {
		loadOpts = append(loadOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, opts.Session),
		))
	}
	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	inv := inventory.NewInventory(opts.Region)

	// EC2 instances
	ec2Client := ec2.NewFromConfig(cfg)
	instances, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe ec2 instances: %w", err)
	}
	for _, res := range instances.Reservations {
		for _, inst := range res.Instances {
			tags := map[string]string{}
			for _, t := range inst.Tags {
				tags[*t.Key] = *t.Value
			}
			inv.Resources = append(inv.Resources, inventory.Resource{
				ID:     *inst.InstanceId,
				Service: "ec2",
				Type:    string(inst.InstanceType),
				Region:  opts.Region,
				State:   string(inst.State.Name),
				Tags:    tags,
				Arn:     *inst.InstanceId,
				Details: map[string]string{
					"image_id":      *inst.ImageId,
					"vpc_id":        *inst.VpcId,
					"subnet_id":     *inst.SubnetId,
					"private_ip":    *inst.PrivateIpAddress,
					"architecture":  string(inst.Architecture),
				},
			})
			inv.Counts["ec2_instance"]++
		}
	}

	// EBS volumes
	volumes, err := ec2Client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe ec2 volumes: %w", err)
	}
	for _, vol := range volumes.Volumes {
		tags := map[string]string{}
		for _, t := range vol.Tags {
			tags[*t.Key] = *t.Value
		}
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     *vol.VolumeId,
			Service: "ec2",
			Type:    string(vol.VolumeType),
			Region:  opts.Region,
			State:   string(vol.State),
			Tags:    tags,
			Arn:     *vol.VolumeId,
			Details: map[string]string{
				"size_gb":   fmt.Sprintf("%d", *vol.Size),
				"az":        *vol.AvailabilityZone,
				"encrypted": fmt.Sprintf("%t", *vol.Encrypted),
			},
		})
		inv.Counts["ebs_volume"]++
	}

	// RDS instances
	rdsClient := rds.NewFromConfig(cfg)
	dbInstances, err := rdsClient.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("describe rds instances: %w", err)
	}
	for _, db := range dbInstances.DBInstances {
		tags := map[string]string{}
		for _, t := range db.TagList {
			tags[*t.Key] = *t.Value
		}
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     strVal(db.DBInstanceIdentifier),
			Service: "rds",
			Type:    strVal(db.DBInstanceClass),
			Region:  opts.Region,
			State:   strVal(db.DBInstanceStatus),
			Tags:    tags,
			Arn:     strVal(db.DBInstanceArn),
			Details: map[string]string{
				"engine":      strVal(db.Engine),
				"engine_ver":  strVal(db.EngineVersion),
				"instance_cl": strVal(db.DBInstanceClass),
				"multi_az":    fmt.Sprintf("%t", *db.MultiAZ),
			},
		})
		inv.Counts["rds_instance"]++
	}

	// S3 buckets
	s3Client := s3.NewFromConfig(cfg)
	buckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("list s3 buckets: %w", err)
	}
	for _, b := range buckets.Buckets {
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     *b.Name,
			Service: "s3",
			Type:    "bucket",
			Region:  opts.Region,
			State:   "active",
			Details: map[string]string{
				"creation_date": b.CreationDate.String(),
			},
		})
		inv.Counts["s3_bucket"]++
	}

	// EKS clusters
	eksClient := eks.NewFromConfig(cfg)
	clustersOut, err := eksClient.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("list eks clusters: %w", err)
	}
	for _, clusterName := range clustersOut.Clusters {
		desc, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: strPtr(clusterName)})
		if err != nil {
			return nil, fmt.Errorf("describe eks cluster %s: %w", clusterName, err)
		}
		c := desc.Cluster
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     strVal(c.Name),
			Service: "eks",
			Type:    strVal(c.Version),
			Region:  opts.Region,
			State:   string(c.Status),
			Arn:     strVal(c.Arn),
			Details: map[string]string{
				"version":  strVal(c.Version),
				"role_arn": strVal(c.RoleArn),
				"endpoint": strVal(c.Endpoint),
			},
		})
		inv.Counts["eks_cluster"]++
	}

	// Lambda functions (names only)
	lambdaClient := lambda.NewFromConfig(cfg)
	functions, err := lambdaClient.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, fmt.Errorf("list lambda functions: %w", err)
	}
	for _, f := range functions.Functions {
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     *f.FunctionName,
			Service: "lambda",
			Type:    string(f.Runtime),
			Region:  opts.Region,
			State:   string(f.State),
			Arn:     strVal(f.FunctionArn),
			Details: map[string]string{
				"runtime":    string(f.Runtime),
				"handler":    *f.Handler,
				"size_bytes": fmt.Sprintf("%d", f.CodeSize),
				"timeout":    fmt.Sprintf("%d", f.Timeout),
				"memory":     fmt.Sprintf("%d", f.MemorySize),
			},
		})
		inv.Counts["lambda_function"]++
	}

	// IAM roles (names only)
	iamClient := iam.NewFromConfig(cfg)
	roles, err := iamClient.ListRoles(ctx, &iam.ListRolesInput{})
	if err != nil {
		return nil, fmt.Errorf("list iam roles: %w", err)
	}
	for _, r := range roles.Roles {
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     *r.RoleName,
			Service: "iam",
			Type:    "role",
			Region:  opts.Region,
			State:   "active",
			Arn:     strVal(r.Arn),
			Details: map[string]string{
				"path":        strVal(r.Path),
				"create_date": r.CreateDate.String(),
			},
		})
		inv.Counts["iam_role"]++
	}

	return inv, nil
}


