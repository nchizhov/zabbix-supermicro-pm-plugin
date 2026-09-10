package support

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

func RunCommand(ctx context.Context, prg string, args ...string) (string, int, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, prg, args...)
	output, err := cmd.Output()
	if err != nil {
		if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			return "", 0, err
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", exitErr.ExitCode(), nil
		}
		return "", 0, err
	}
	return strings.TrimSuffix(string(output), "\n"), 0, nil
}
