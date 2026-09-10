//go:build windows

package ipmi

import (
	"context"
	"support"
)

var errorCodesList = exitCodes{}

func ipmiData(ctx context.Context, tool string) (string, error) {
	data, code, err := support.RunCommand(ctx, tool, "-pminfo")
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", errorCodesList.getErrorMessage(code)
	}
	return data, nil
}
