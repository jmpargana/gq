package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
	cmd := exec.Command(cliPath, "version")

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
		{
			desc:    "iter",
			stdin:   `[1,2]`,
			program: `.[]`,
			wantOut: `1
2
`,
			wantErr: "",
		},
		{
			desc:    "capture iter",
			stdin:   `[1,2]`,
			program: `[.[]]`,
			wantOut: `[
  1,
  2
]
`,
			wantErr: "",
		},
		{
			desc:    "pipe iter",
			stdin:   `[1,2]`,
			program: `.[] | [.]`,
			wantOut: `[
  1
]
[
  2
]
`,
			wantErr: "",
		},
		{
			desc:    "cartesian product",
			stdin:   `[1, 2]`,
			program: `{a:.[], b:.[]}`,
			wantOut: `{
  "a": 1,
  "b": 1
}
{
  "a": 1,
  "b": 2
}
{
  "a": 2,
  "b": 1
}
{
  "a": 2,
  "b": 2
}
`,
			wantErr: "",
		},
		{
			desc:    "piped cartesian product",
			stdin:   `[1, 2]`,
			program: `{a:.[], b:.[]} | .[]`,
			wantOut: `1
1
1
2
2
1
2
2
`,
			wantErr: "piped flattening",
		},
		{
			desc:    "piped iter",
			stdin:   `[[1,2], [3]]`,
			program: `[.[].[]]`,
			wantOut: `[
  1,
  2,
  3
]
`,
			wantErr: "",
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
				var gotVal, wantVal any
				json.Unmarshal([]byte(got), &gotVal)
				json.Unmarshal([]byte(tC.wantOut), &wantVal)
				if !reflect.DeepEqual(gotVal, wantVal) {
					t.Fatalf("unexpected output:\ngot:%s\nwanted:%s\n", got, tC.wantOut)
				}
				// if !strings.Contains(got, tC.wantOut) {
				// 	t.Fatalf("unexpected output:\ngot:%s\nwanted:%s\n", got, tC.wantOut)
				// }
			}
		})
	}
}

func TestCLI_RootTestDataString(t *testing.T) {
	testCases := []struct {
		desc, query, file string
	}{
		{
			desc:  "root without stdin returns usage",
			query: `.[] | .parents[].sha`,
			file:  `output1.json`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmd := exec.Command(cliPath, tC.query)
			f, err := os.Open("./testdata/test.json")
			if err != nil {
				t.Fatalf("failed opening file: %v", err)
			}
			expectedFile, err := os.Open(fmt.Sprintf("./testdata/%s", tC.file))
			if err != nil {
				t.Fatalf("failed opening file: %v", err)
			}

			cmd.Stdin = f
			expected, _ := io.ReadAll(expectedFile)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()
			if err != nil {
				t.Fatalf("err: %s", err)
			}
			got := stdout.String()
			if !reflect.DeepEqual(got, string(expected)) {
				t.Fatalf("unexpected output:\ngot:%s\nwanted:%s\n", got, expected)
			}
		})
	}
}

func TestCLI_RootTestDataJSON(t *testing.T) {
	testCases := []struct {
		desc, query, file string
	}{
		{
			desc:  "tutorial example 1",
			query: `.[0] | {message: .commit.message, name: .commit.committer.name}`,
			file:  `output2.json`,
		},
		{
			desc:  "tutorial example 2",
			query: `[.[] | {message: .commit.message, name: .commit.committer.name}]`,
			file:  `output3.json`,
		},
		{
			desc:  "tutorial example 3",
			query: `[.[] | {message: .commit.message, name: .commit.committer.name, parents: [.parents[].html_url]}]`,
			file:  `output4.json`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmd := exec.Command(cliPath, tC.query)
			f, err := os.Open("./testdata/test.json")
			if err != nil {
				t.Fatalf("failed opening file: %v", err)
			}
			expectedFile, err := os.Open(fmt.Sprintf("./testdata/%s", tC.file))
			if err != nil {
				t.Fatalf("failed opening file: %v", err)
			}

			cmd.Stdin = f
			expected, _ := io.ReadAll(expectedFile)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			var gotJSON, wantJSON any
			got := stdout.String()
			json.Unmarshal([]byte(got), &gotJSON)
			json.Unmarshal(expected, &wantJSON)
			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Fatalf("unexpected output:\ngot:%s\nwanted:%s\n", got, expected)
			}
		})
	}
}
