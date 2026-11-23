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

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func scanLambdaFunctions(ctx context.Context, lambdaClient LambdaAPI, region string, inv *inventory.Inventory) {
	cov := inventory.ServiceCoverage{Attempted: true}
	input := &lambda.ListFunctionsInput{}
	svcCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	for {
		page, err := lambdaClient.ListFunctions(svcCtx, input)
		if err != nil {
			cov.Completed = false
			cov.FirstError = fmt.Sprintf("list lambda functions: %v", err)
			inv.Manifest.Services["lambda_function"] = cov
			return
		}
		for _, f := range page.Functions {
			timeout := int64(0)
			if f.Timeout != nil {
				timeout = int64(*f.Timeout)
			}
			memory := int64(0)
			if f.MemorySize != nil {
				memory = int64(*f.MemorySize)
			}
			inv.Resources = append(inv.Resources, inventory.Resource{
				ID:     strVal(f.FunctionName),
				Service: "lambda",
				Type:    string(f.Runtime),
				Region:  region,
				State:   string(f.State),
				Tags:    map[string]string{},
				Arn:     strVal(f.FunctionArn),
				Details: map[string]string{
					"runtime":    string(f.Runtime),
					"handler":    strVal(f.Handler),
					"size_bytes": fmt.Sprintf("%d", f.CodeSize),
					"timeout":    fmt.Sprintf("%d", timeout),
					"memory":     fmt.Sprintf("%d", memory),
					"last_mod":   strVal(f.LastModified),
				},
			})
			cov.ResourceCount++
		}
		if page.NextMarker == nil || *page.NextMarker == "" {
			cov.PaginationComplete = true
			break
		}
		input.Marker = page.NextMarker
	}
	cov.Completed = true
	inv.Manifest.Services["lambda_function"] = cov
}
