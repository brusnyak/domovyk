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

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanS3Buckets(ctx context.Context, s3Client S3API, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	out, err := s3Client.ListBuckets(svcCtx, &s3.ListBucketsInput{})
	if err != nil {
		cov.Completed = false
		cov.FirstError = fmt.Sprintf("list s3 buckets: %v", err)
		inv.Manifest.Services["s3_bucket"] = cov
		return
	}

	for _, b := range out.Buckets {
		inv.Resources = append(inv.Resources, inventory.Resource{
			ID:     strVal(b.Name),
			Service: "s3",
			Type:    "bucket",
			Region:  region,
			State:   "active",
			Details: map[string]string{
				"creation_date": timeVal(b.CreationDate),
			},
		})
		cov.ResourceCount++
	}
	cov.Completed = true
	inv.Manifest.Services["s3_bucket"] = cov
}
