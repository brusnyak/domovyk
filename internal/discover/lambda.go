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
