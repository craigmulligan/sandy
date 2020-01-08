package main

import (
	"testing"
)

func TestSum(t *testing.T) {
	s := []string{"sandy.go"}
	reqs, err := Exec("cat", s)

	if err != nil {
		t.Errorf("Something went wrong")
	}

	if len(reqs) != 2 {
		t.Errorf("reqs count was incorrect, got: %d, want: %d.", len(reqs), 2)
	}
}
