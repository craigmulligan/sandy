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
	Exec("cat", s, ask)

	if call_count != 9 {
		t.Errorf("call_count was incorrect, got: %d, want: %d.", call_count, 10)
	}
}
