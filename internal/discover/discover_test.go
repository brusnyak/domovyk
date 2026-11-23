package discover

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/brusnyak/domovyk/internal/inventory"
)

// ---- mock EC2 ----

type mockEC2 struct {
	instances []ec2.DescribeInstancesOutput
	volumes   []ec2.DescribeVolumesOutput
	err       error
}

func (m *mockEC2) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.instances) == 0 {
		return &ec2.DescribeInstancesOutput{}, nil
	}
	out := m.instances[0]
	m.instances = m.instances[1:]
	return &out, nil
}

func (m *mockEC2) DescribeVolumes(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.volumes) == 0 {
		return &ec2.DescribeVolumesOutput{}, nil
	}
	out := m.volumes[0]
	m.volumes = m.volumes[1:]
	return &out, nil
}

// ---- mock RDS ----

type mockRDS struct {
	dbInstances []rds.DescribeDBInstancesOutput
	err         error
}

func (m *mockRDS) DescribeDBInstances(ctx context.Context, params *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.dbInstances) == 0 {
		return &rds.DescribeDBInstancesOutput{}, nil
	}
	out := m.dbInstances[0]
	m.dbInstances = m.dbInstances[1:]
	return &out, nil
}

// ---- mock S3 ----

type mockS3 struct {
	buckets *s3.ListBucketsOutput
	err     error
}

func (m *mockS3) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.buckets, nil
}

// ---- mock EKS ----

type mockEKS struct {
	clusters   []eks.ListClustersOutput
	clustersBy map[string]*eks.DescribeClusterOutput
	err        error
}

func (m *mockEKS) ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.clusters) == 0 {
		return &eks.ListClustersOutput{}, nil
	}
	out := m.clusters[0]
	m.clusters = m.clusters[1:]
	return &out, nil
}

func (m *mockEKS) DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.clustersBy == nil {
		return &eks.DescribeClusterOutput{}, nil
	}
	if out, ok := m.clustersBy[*params.Name]; ok {
		return out, nil
	}
	return &eks.DescribeClusterOutput{}, nil
}

// ---- mock Lambda ----

type mockLambda struct {
	functions []lambda.ListFunctionsOutput
	err       error
}

func (m *mockLambda) ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.functions) == 0 {
		return &lambda.ListFunctionsOutput{}, nil
	}
	out := m.functions[0]
	m.functions = m.functions[1:]
	return &out, nil
}

// ---- mock IAM ----

type mockIAM struct {
	roles []iam.ListRolesOutput
	err   error
}

func (m *mockIAM) ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.roles) == 0 {
		return &iam.ListRolesOutput{}, nil
	}
	out := m.roles[0]
	m.roles = m.roles[1:]
	return &out, nil
}

// ---- mock STS ----

type mockSTS struct {
	accountID *string
	err       error
}

func (m *mockSTS) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &sts.GetCallerIdentityOutput{Account: m.accountID}, nil
}

var _ STSClient = (*mockSTS)(nil)
var _ EC2API = (*mockEC2)(nil)
var _ RDSAPI = (*mockRDS)(nil)
var _ S3API = (*mockS3)(nil)
var _ EKSAPI = (*mockEKS)(nil)
var _ LambdaAPI = (*mockLambda)(nil)
var _ IAMAPI = (*mockIAM)(nil)

func pstr(s string) *string { return &s }

func TestScanEC2InstancesNilFields(t *testing.T) {
	m := &mockEC2{
		instances: []ec2.DescribeInstancesOutput{
			{
				Reservations: []ec2types.Reservation{
					{
						Instances: []ec2types.Instance{
							{InstanceId: nil, State: &ec2types.InstanceState{}},
						},
					},
				},
			},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanEC2Instances(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(inv.Resources))
	}
	r := inv.Resources[0]
	if r.ID != "" {
		t.Errorf("id = %q, want empty", r.ID)
	}
	if r.State != "" {
		t.Errorf("state = %q, want empty", r.State)
	}
	cov := inv.Manifest.Services["ec2_instance"]
	if !cov.Completed {
		t.Errorf("coverage not completed")
	}
	if !cov.PaginationComplete {
		t.Errorf("pagination not complete")
	}
}

func TestScanEC2InstancesPagination(t *testing.T) {
	m := &mockEC2{
		instances: []ec2.DescribeInstancesOutput{
			{
				Reservations: []ec2types.Reservation{
					{Instances: []ec2types.Instance{{InstanceId: pstr("i-1")}}},
				},
				NextToken: pstr("tok1"),
			},
			{
				Reservations: []ec2types.Reservation{
					{Instances: []ec2types.Instance{{InstanceId: pstr("i-2")}}},
				},
			},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanEC2Instances(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 2 {
		t.Fatalf("expected 2 resources across pages, got %d", len(inv.Resources))
	}
	cov := inv.Manifest.Services["ec2_instance"]
	if !cov.PaginationComplete {
		t.Errorf("expected pagination complete")
	}
	if cov.ResourceCount != 2 {
		t.Errorf("resource_count = %d, want 2", cov.ResourceCount)
	}
}

func TestScanEC2InstancesAccessDenied(t *testing.T) {
	m := &mockEC2{err: errors.New("AccessDenied")}
	inv := inventory.NewInventory("eu-central-1")
	scanEC2Instances(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 0 {
		t.Errorf("expected 0 resources on error, got %d", len(inv.Resources))
	}
	cov := inv.Manifest.Services["ec2_instance"]
	if cov.Completed {
		t.Errorf("expected not completed on access denied")
	}
	if cov.FirstError == "" {
		t.Errorf("expected first_error recorded")
	}
}

func TestScanRDSNilFields(t *testing.T) {
	m := &mockRDS{
		dbInstances: []rds.DescribeDBInstancesOutput{
			{DBInstances: []rdstypes.DBInstance{{}}},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanRDSInstances(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(inv.Resources))
	}
	r := inv.Resources[0]
	if r.ID != "" || r.State != "" {
		t.Errorf("unexpected nil-field values: %+v", r)
	}
}

func TestScanRDSPagination(t *testing.T) {
	m := &mockRDS{
		dbInstances: []rds.DescribeDBInstancesOutput{
			{DBInstances: []rdstypes.DBInstance{{DBInstanceIdentifier: pstr("db-1")}}, Marker: pstr("m1")},
			{DBInstances: []rdstypes.DBInstance{{DBInstanceIdentifier: pstr("db-2")}}},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanRDSInstances(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 2 {
		t.Fatalf("expected 2 resources across pages, got %d", len(inv.Resources))
	}
}

func TestScanS3AccessDenied(t *testing.T) {
	m := &mockS3{err: errors.New("AccessDenied")}
	inv := inventory.NewInventory("eu-central-1")
	scanS3Buckets(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 0 {
		t.Errorf("expected 0 resources on error, got %d", len(inv.Resources))
	}
	cov := inv.Manifest.Services["s3_bucket"]
	if cov.Completed {
		t.Errorf("expected not completed on access denied")
	}
}

func TestScanEKSNilCluster(t *testing.T) {
	m := &mockEKS{
		clusters:   []eks.ListClustersOutput{{Clusters: []string{"c1"}}},
		clustersBy: map[string]*eks.DescribeClusterOutput{"c1": nil},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanEKSClusters(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 0 {
		t.Fatalf("expected 0 resources for nil cluster, got %d", len(inv.Resources))
	}
	cov := inv.Manifest.Services["eks_cluster"]
	if !cov.Completed {
		t.Errorf("expected coverage completed despite nil cluster")
	}
}

func TestScanEKSNilFields(t *testing.T) {
	m := &mockEKS{
		clusters: []eks.ListClustersOutput{{Clusters: []string{"c1"}}},
		clustersBy: map[string]*eks.DescribeClusterOutput{
			"c1": {Cluster: &ekstypes.Cluster{Name: pstr("c1")}},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanEKSClusters(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(inv.Resources))
	}
	r := inv.Resources[0]
	if r.ID != "c1" {
		t.Errorf("id = %q, want c1", r.ID)
	}
}

func TestScanEKSAccessDenied(t *testing.T) {
	m := &mockEKS{err: errors.New("AccessDenied")}
	inv := inventory.NewInventory("eu-central-1")
	scanEKSClusters(context.Background(), m, "eu-central-1", inv)
	cov := inv.Manifest.Services["eks_cluster"]
	if cov.Completed {
		t.Errorf("expected not completed on access denied")
	}
	if cov.FirstError == "" {
		t.Errorf("expected first_error recorded")
	}
}

func TestScanLambdaNilFields(t *testing.T) {
	m := &mockLambda{
		functions: []lambda.ListFunctionsOutput{
			{Functions: []lambdatypes.FunctionConfiguration{{FunctionName: pstr("f-1")}}},
		},
	}
	inv := inventory.NewInventory("eu-central-1")
	scanLambdaFunctions(context.Background(), m, "eu-central-1", inv)
	if len(inv.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(inv.Resources))
	}
}

func TestScanIAMAccessDenied(t *testing.T) {
	m := &mockIAM{err: errors.New("AccessDenied")}
	inv := inventory.NewInventory("eu-central-1")
	scanIAMRoles(context.Background(), m, "eu-central-1", inv)
	cov := inv.Manifest.Services["iam_role"]
	if cov.Completed {
		t.Errorf("expected not completed on access denied")
	}
	if cov.FirstError == "" {
		t.Errorf("expected first_error recorded")
	}
}

func TestAccountIdentitySTSError(t *testing.T) {
	m := &mockSTS{err: errors.New("expired token")}
	id, err := accountIdentity(context.Background(), m)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if id != "" {
		t.Errorf("expected empty id on error, got %q", id)
	}
}

func TestDiscoverNilContext(t *testing.T) {
	_, err := Discover(nil, Options{Region: "eu-central-1"})
	if err == nil {
		t.Errorf("expected error for nil context")
	}
}

func TestInventoryCoveragePartialFailure(t *testing.T) {
	ec2m := &mockEC2{err: errors.New("AccessDenied")}
	s3m := &mockS3{buckets: &s3.ListBucketsOutput{Buckets: []s3types.Bucket{{Name: pstr("b-1")}}}}

	inv := inventory.NewInventory("eu-central-1")
	scanEC2Instances(context.Background(), ec2m, "eu-central-1", inv)
	scanS3Buckets(context.Background(), s3m, "eu-central-1", inv)

	if len(inv.Resources) != 1 || inv.Resources[0].Service != "s3" {
		t.Errorf("expected only s3 resource, got %d", len(inv.Resources))
	}
	ec2Cov := inv.Manifest.Services["ec2_instance"]
	if ec2Cov.Completed {
		t.Errorf("expected ec2 not completed")
	}
	s3Cov := inv.Manifest.Services["s3_bucket"]
	if !s3Cov.Completed {
		t.Errorf("expected s3 completed")
	}
}
