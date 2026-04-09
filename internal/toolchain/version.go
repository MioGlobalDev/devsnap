package toolchain

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"devenv-snapshot/internal/runner"
)

func NodeVersion(ctx context.Context) (string, bool, error) {
	if _, err := exec.LookPath("node"); err != nil {
		return "", false, nil
	}
	res, err := runner.Run(ctx, runner.Cmd{
		Name:    "node",
		Args:    []string{"--version"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return "", true, err
	}
	v := strings.TrimSpace(res.Stdout)
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return "", true, nil
	}
	return v, true, nil
}

func PythonVersion(ctx context.Context) (string, bool, error) {
	if _, err := exec.LookPath("python"); err != nil {
		return "", false, nil
	}
	res, err := runner.Run(ctx, runner.Cmd{
		Name:    "python",
		Args:    []string{"--version"},
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return "", true, err
	}
	// python sometimes prints version to stderr; handle both
	out := strings.TrimSpace(res.Stdout)
	if out == "" {
		out = strings.TrimSpace(res.Stderr)
	}
	out = strings.TrimSpace(out)
	const prefix = "Python "
	out = strings.TrimPrefix(out, prefix)
	if out == "" {
		return "", true, nil
	}
	return out, true, nil
}

