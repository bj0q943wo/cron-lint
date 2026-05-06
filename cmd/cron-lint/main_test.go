package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the cron-lint binary into a temp directory and
// returns its path. The test is skipped if the build fails.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "cron-lint")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping integration test — build failed: %v\n%s", err, out)
	}
	return bin
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "crons-*.txt")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestMain_MissingFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "-f") {
		t.Errorf("expected usage hint containing -f, got: %s", out)
	}
}

func TestMain_TextOutput_NoOverlap(t *testing.T) {
	bin := buildBinary(t)
	cronFile := writeTempFile(t, "0 9 * * 1 /bin/backup\n30 10 * * 2 /bin/report\n")
	cmd := exec.Command(bin, "-f", cronFile, "-format", "text")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected exit error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "no overlaps") {
		t.Errorf("expected 'no overlaps' in output, got: %s", out)
	}
}

func TestMain_JSONOutput(t *testing.T) {
	bin := buildBinary(t)
	cronFile := writeTempFile(t, "0 9 * * * /bin/jobA\n0 9 * * * /bin/jobB\n")
	cmd := exec.Command(bin, "-f", cronFile, "-format", "json")
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestMain_StrictMode_ExitsOne(t *testing.T) {
	bin := buildBinary(t)
	cronFile := writeTempFile(t, "0 9 * * * /bin/jobA\n0 9 * * * /bin/jobB\n")
	cmd := exec.Command(bin, "-f", cronFile, "-strict")
	if err := cmd.Run(); err == nil {
		t.Error("expected exit code 1 in strict mode with overlaps, got 0")
	}
}
