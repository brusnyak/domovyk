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

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanRDSInstances(ctx context.Context, rdsClient RDSAPI, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &rds.DescribeDBInstancesInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for {
		page, err := rdsClient.DescribeDBInstances(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("describe rds instances: %v", err)
			inv.Manifest.Services["rds_instance"] = cov
			return
		}
		for _, db := range page.DBInstances {
			multiAZ := false
			if db.MultiAZ != nil {
				multiAZ = *db.MultiAZ
			}
			tags := map[string]string{}
			for _, t := range db.TagList {
				if t.Key == nil || t.Value == nil {
					continue
				}
				tags[*t.Key] = *t.Value
			}
			inv.Resources = append(inv.Resources, inventory.Resource{
				ID:     strVal(db.DBInstanceIdentifier),
				Service: "rds",
				Type:    strVal(db.DBInstanceClass),
				Region:  region,
				State:   strVal(db.DBInstanceStatus),
				Tags:    tags,
				Arn:     strVal(db.DBInstanceArn),
				Details: map[string]string{
					"engine":      strVal(db.Engine),
					"engine_ver":  strVal(db.EngineVersion),
					"instance_cl": strVal(db.DBInstanceClass),
					"multi_az":    boolVal(&multiAZ),
					"create_date": timeVal(db.InstanceCreateTime),
				},
			})
			cov.ResourceCount++
		}
		if page.Marker == nil || *page.Marker == "" {
			cov.PaginationComplete = true
			break
		}
		input.Marker = page.Marker
	}
	cov.Completed = true
	inv.Manifest.Services["rds_instance"] = cov
}
