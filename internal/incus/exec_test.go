package incus

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestExecCmdTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell script based timeout test is not supported on Windows")
	}

	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "incus")
	script := "#!/bin/sh\nsleep 1\nprintf 'delayed'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	result := client{timeout: 10 * time.Millisecond}.runQuiet([]string{"version"}, nil)
	if result.Error == nil {
		t.Fatal("runQuiet() error = nil, want timeout error")
	}
	if !strings.Contains(result.Error.Error(), "timed out") {
		t.Fatalf("timeout error = %q, want to contain %q", result.Error.Error(), "timed out")
	}
}
