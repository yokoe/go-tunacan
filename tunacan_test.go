package main

import (
	"testing"

	"./concatenator"
)

func TestNoInput(t *testing.T) {
	err := concatenator.Concat(nil, "test_output.png")
	if err == nil {
		t.Fail()
	}
}
