package incus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// run executes an incus command, forwarding stdout to the terminal for
// real-time progress while still capturing it for result inspection.
func (c client) run(args []string, stdin []byte) *Result {
	return c.execCmd(args, stdin, true)
}

// runQuiet executes an incus command capturing all output without forwarding.
func (c client) runQuiet(args []string, stdin []byte) *Result {
	return c.execCmd(args, stdin, false)
}

// execCmd is the shared implementation for run and runQuiet.
func (c client) execCmd(args []string, stdin []byte, forwardStdout bool) *Result {
	ctx := context.Background()
	cancel := func() {}
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, "incus", args...)

	if c.debug {
		fmt.Printf("[debug] %s\n", cmd.String())
	}

	var stdout, stderr bytes.Buffer
	if forwardStdout {
		cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	} else {
		cmd.Stdout = &stdout
	}
	cmd.Stderr = &stderr

	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}

	err := cmd.Run()

	result := &Result{
		Command: "incus " + strings.Join(args, " "),
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			result.Error = fmt.Errorf("command timed out after %s", c.timeout)
			return result
		}
		result.Error = fmt.Errorf("%w: %s", err, strings.TrimSpace(result.Stderr))
	}
	return result
}
