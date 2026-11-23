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

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanEBSVolumes(ctx context.Context, ec2Client EC2API, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &ec2.DescribeVolumesInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for {
		page, err := ec2Client.DescribeVolumes(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("describe ec2 volumes: %v", err)
			inv.Manifest.Services["ebs_volume"] = cov
			return
		}
		for _, vol := range page.Volumes {
			tags := map[string]string{}
			for _, t := range vol.Tags {
				if t.Key == nil || t.Value == nil {
					continue
				}
				tags[*t.Key] = *t.Value
			}
			encrypted := false
			if vol.Encrypted != nil {
				encrypted = *vol.Encrypted
			}
			inv.Resources = append(inv.Resources, inventory.Resource{
				ID:     strVal(vol.VolumeId),
				Service: "ec2",
				Type:    string(vol.VolumeType),
				Region:  region,
				State:   string(vol.State),
				Tags:    tags,
				Arn:     strVal(vol.VolumeId),
				Details: map[string]string{
					"size_gb":   fmt.Sprintf("%d", strInt32(vol.Size)),
					"az":        strVal(vol.AvailabilityZone),
					"encrypted": boolVal(&encrypted),
				},
			})
			cov.ResourceCount++
		}
		if page.NextToken == nil || *page.NextToken == "" {
			cov.PaginationComplete = true
			break
		}
		input.NextToken = page.NextToken
	}
	cov.Completed = true
	inv.Manifest.Services["ebs_volume"] = cov
}
