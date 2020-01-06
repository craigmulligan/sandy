package main

import (
	"testing"
)

var call_count int

func ask(path string) (string, error) {
	call_count += 1
	return path, nil
}

func TestSum(t *testing.T) {
	// introspect :)
	call_count = 0
	s := []string{"sandy.go"}
	reqs, err := Exec("cat", s)

	if err != nil {
		t.Errorf("Something went wrong")
	}

	if len(reqs) != 2 {
		t.Errorf("reqs count was incorrect, got: %d, want: %d.", len(reqs), 2)
	}
}
