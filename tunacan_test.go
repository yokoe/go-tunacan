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

func TestInvalidInput(t *testing.T) {
	err := concat([]string{"***********"}, "test_output.png")
	if err == nil {
		t.Fail()
	}
}
