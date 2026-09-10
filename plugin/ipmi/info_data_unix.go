//go:build !windows

package ipmi

import (
	"context"
	"support"
)

var errorCodesList = exitCodes{
	42: "incorrect sudo configured",
}

func ipmiData(ctx context.Context, tool string) (string, error) {
	data, code, err := support.RunCommand(ctx, "sudo", tool, "-pminfo")
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", errorCodesList.getErrorMessage(code)
	}
	return data, nil
}
