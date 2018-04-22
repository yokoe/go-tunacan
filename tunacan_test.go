package main

import (
	"testing"
)

func TestNoInput(t *testing.T) {
	err := concat(nil, "test_output.png")
	if err == nil {
		t.Fail()
	}
}
