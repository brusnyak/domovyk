package discover

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanIAMRoles(ctx context.Context, iamClient IAMAPI, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &iam.ListRolesInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for {
		page, err := iamClient.ListRoles(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("list iam roles: %v", err)
			inv.Manifest.Services["iam_role"] = cov
			return
		}
		for _, r := range page.Roles {
			inv.Resources = append(inv.Resources, inventory.Resource{
				ID:     strVal(r.RoleName),
				Service: "iam",
				Type:    "role",
				Region:  region,
				State:   "active",
				Arn:     strVal(r.Arn),
				Details: map[string]string{
					"path":        strVal(r.Path),
					"create_date": timeVal(r.CreateDate),
				},
			})
			cov.ResourceCount++
		}
		if page.IsTruncated && strVal(page.Marker) != "" {
			input.Marker = page.Marker
			continue
		}
		cov.PaginationComplete = true
		break
	}
	cov.Completed = true
	inv.Manifest.Services["iam_role"] = cov
}
