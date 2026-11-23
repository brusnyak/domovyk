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

func scanEC2Instances(ctx context.Context, ec2Client EC2API, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &ec2.DescribeInstancesInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for {
		page, err := ec2Client.DescribeInstances(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("describe ec2 instances: %v", err)
			inv.Manifest.Services["ec2_instance"] = cov
			return
		}
		for _, res := range page.Reservations {
			for _, inst := range res.Instances {
				tags := map[string]string{}
				for _, t := range inst.Tags {
					if t.Key == nil || t.Value == nil {
						continue
					}
					tags[*t.Key] = *t.Value
				}
				state := ""
				if inst.State != nil && inst.State.Name != "" {
					state = string(inst.State.Name)
				}
				inv.Resources = append(inv.Resources, inventory.Resource{
					ID:     strVal(inst.InstanceId),
					Service: "ec2",
					Type:    string(inst.InstanceType),
					Region:  region,
					State:   state,
					Tags:    tags,
					Arn:     strVal(inst.InstanceId),
					Details: map[string]string{
						"image_id":      strVal(inst.ImageId),
						"vpc_id":        strVal(inst.VpcId),
						"subnet_id":     strVal(inst.SubnetId),
						"private_ip":    strVal(inst.PrivateIpAddress),
						"architecture":  string(inst.Architecture),
						"kind":          "ec2_instance",
					},
				})
				cov.ResourceCount++
			}
		}
		if page.NextToken == nil || *page.NextToken == "" {
			cov.PaginationComplete = true
			break
		}
		input.NextToken = page.NextToken
	}
	cov.Completed = true
	inv.Manifest.Services["ec2_instance"] = cov
}
