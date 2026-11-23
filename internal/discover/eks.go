/*
 * Copyright 2026 Yegor
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package discover

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanEKSClusters(ctx context.Context, eksClient EKSAPI, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &eks.ListClustersInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var clusterNames []string
	for {
		out, err := eksClient.ListClusters(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("list eks clusters: %v", err)
			inv.Manifest.Services["eks_cluster"] = cov
			return
		}
		clusterNames = append(clusterNames, out.Clusters...)
		if out.NextToken == nil || *out.NextToken == "" {
			cov.PaginationComplete = true
			break
		}
		input.NextToken = out.NextToken
	}

	for _, clusterName := range clusterNames {
		desc, err := eksClient.DescribeCluster(svcCtx, &eks.DescribeClusterInput{Name: stringPtr(clusterName)})
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("describe eks cluster %s: %v", clusterName, err)
			inv.Manifest.Services["eks_cluster"] = cov
			return
		}
		if desc == nil || desc.Cluster == nil {
			continue
		}
		c := desc.Cluster
		status := ""
		if c.Status != "" {
			status = string(c.Status)
		}
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     strVal(c.Name),
			Service: "eks",
			Type:    strVal(c.Version),
			Region:  region,
			State:   status,
			Arn:     strVal(c.Arn),
			Details: map[string]string{
				"version":  strVal(c.Version),
				"role_arn": strVal(c.RoleArn),
				"endpoint": strVal(c.Endpoint),
				"created":  timeVal(c.CreatedAt),
			},
		})
		cov.ResourceCount++
	}
	cov.Completed = true
	inv.Manifest.Services["eks_cluster"] = cov
}
