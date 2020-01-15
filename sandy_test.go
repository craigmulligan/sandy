package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	s := []string{"password.txt"}
	patterns := []string{""}
	reqs, err := Exec("cat", s, patterns, patterns)

	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}

	if len(reqs) != 2 {
		t.Errorf("reqs count was incorrect, got: %d, want: %d.", len(reqs), 2)
	}
}

func TestInput(t *testing.T) {
	cmd := exec.Command("./sandy", "cat", "./password.txt")
	var out bytes.Buffer
	var in bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdout = &in
	in.Write([]byte("n\n\r"))

	err := cmd.Run()
	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}
	if strings.Contains(out.String(), "Blocked READ on ...") {
		t.Errorf("Expected %s output got %s", "Blocked READ on ...", out.String())
	}
}

// TODO we probably should instead just pass a mock reader for stdin into the Exec function and then call the fn
// directly rather that full bin tests
func TestAllowList(t *testing.T) {
	cmd := exec.Command("./sandy", "--y", "*.so", "--y", "*.txt", "cat", "./password.txt")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}
	if out.String() != "123\n" {
		t.Errorf("Expected %s output got %s", "123", out.String())
	}
}

func TestBlockList(t *testing.T) {
	cmd := exec.Command("./sandy", "--y", "*.so", "--n", "*.txt", "cat", "./password.txt")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}
	if !strings.Contains(out.String(), "Blocked READ on ") {
		t.Errorf("Expected %s output got %s", "Blocked READ on", out.String())
	}
}

func TestHelp(t *testing.T) {
	cmd := exec.Command("./sandy", "-h")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		t.Errorf("Something went wrong: %v", err)
	}

	if strings.Contains(out.String(), "Usage of ./sandy:") {
		t.Errorf("Expected %s output got %s", "123", out.String())
	}
}
