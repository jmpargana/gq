package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var cliPath = "gq"

func TestMain(m *testing.M) {
	tmp := os.TempDir()
	cliPath = filepath.Join(tmp, "mycli-test")

	cmd := exec.Command("go", "build", "-o", cliPath, "../cmd/gq")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Remove(cliPath)
	os.Exit(code)
}

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command(cliPath, "--version")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Version:") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestCLI_Root(t *testing.T) {
	testCases := []struct {
		desc, stdin, program, wantOut, wantErr string
	}{
		{
			desc:    "root without stdin returns usage",
			stdin:   "",
			program: "",
			wantOut: "",
			wantErr: "Error",
		},
		{
			desc:    "root without stdin returns usage",
			stdin:   `{"a": "b"}`,
			program: "",
			wantOut: "",
			wantErr: "no program provided",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var cmd *exec.Cmd
			if tC.program == "" {
				cmd = exec.Command(cliPath)
			} else {
				cmd = exec.Command(cliPath, tC.program)
			}

			if tC.stdin != "" {
				cmd.Stdin = bytes.NewBuffer([]byte(tC.stdin))
			}

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				got := stderr.String()
				if !strings.Contains(got, tC.wantErr) {
					t.Fatalf("err: %s\nunexpected stderr output:\ngot: %s\nwanted: %s\n", err, got, tC.wantErr)
				}
			} else {
				got := stdout.String()
				if !strings.Contains(got, tC.wantOut) {
					t.Fatalf("unexpected output:\ngot:%s\nwanted:%s\n", got, tC.wantOut)
				}
			}
		})
	}
}
