package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestExec(t *testing.T) {
	s := []string{"password.txt"}
	patterns := []string{""}
	reqs, err := Exec("cat", s, patterns)

	if err != nil {
		t.Errorf("Something went wrong")
	}

	if len(reqs) != 2 {
		t.Errorf("reqs count was incorrect, got: %d, want: %d.", len(reqs), 2)
	}
}

func TestIntegration(t *testing.T) {
	cmd := exec.Command("./sandy", "--y", "*.so", "--y", "*.txt", "cat", "./password.txt")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Errorf("Something went wrong")
	}
	if out.String() != "123\n" {
		t.Errorf("Expected %s output got %s", "123", out.String())
	}
}
