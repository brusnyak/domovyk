package discover

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func accountIdentity(ctx context.Context, client STSClient) (string, error) {
	out, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("sts get caller identity: %w", err)
	}
	if out.Account == nil {
		return "", nil
	}
	return *out.Account, nil
}
