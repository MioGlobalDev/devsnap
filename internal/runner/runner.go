package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Cmd struct {
	Dir     string
	Name    string
	Args    []string
	Timeout time.Duration
}

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type Error struct {
	Cmd    Cmd
	Result Result
}

func (e *Error) Error() string {
	sb := strings.Builder{}
	sb.WriteString("command failed: ")
	sb.WriteString(FormatCmd(e.Cmd))
	if e.Cmd.Dir != "" {
		sb.WriteString(fmt.Sprintf(" (dir=%s)", e.Cmd.Dir))
	}
	if e.Result.ExitCode != 0 {
		sb.WriteString(fmt.Sprintf(" (exit=%d)", e.Result.ExitCode))
	}
	if strings.TrimSpace(e.Result.Stderr) != "" {
		sb.WriteString("\n")
		sb.WriteString(e.Result.Stderr)
	}
	return sb.String()
}

func FormatCmd(c Cmd) string {
	all := append([]string{c.Name}, c.Args...)
	return strings.Join(all, " ")
}

func Run(ctx context.Context, c Cmd) (Result, error) {
	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Dir = c.Dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	res := Result{
		ExitCode: exitCode(err),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}

	if err != nil {
		return res, &Error{Cmd: c, Result: res}
	}
	return res, nil
}

func RunShell(ctx context.Context, dir string, script string, timeout time.Duration) (Result, error) {
	script = strings.TrimSpace(script)
	if script == "" {
		return Result{ExitCode: 0}, nil
	}

	if runtime.GOOS == "windows" {
		return Run(ctx, Cmd{
			Dir:     dir,
			Name:    "cmd.exe",
			Args:    []string{"/C", script},
			Timeout: timeout,
		})
	}

	return Run(ctx, Cmd{
		Dir:     dir,
		Name:    "sh",
		Args:    []string{"-lc", script},
		Timeout: timeout,
	})
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if ok := errors.As(err, &ee); ok && ee.ProcessState != nil {
		return ee.ProcessState.ExitCode()
	}
	// context deadline/cancel or spawn failure
	return 1
}

