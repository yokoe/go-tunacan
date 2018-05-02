package tunacan

import (
	"testing"
)

func TestNoInput(t *testing.T) {
	err := Concat(nil, "test_output.png")
	if err == nil {
		t.Fail()
	}
}

func TestInvalidInput(t *testing.T) {
	err := Concat([]string{"***********"}, "test_output.png")
	if err == nil {
		t.Fail()
	}
}

func TestNoImages(t *testing.T) {
	outputImage := ConcatImages(nil)
	if outputImage == nil {
		t.Fail()
	}
}
