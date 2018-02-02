package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetError(t *testing.T) {
	setError("Testing error.", 400)

	assert.Equal(t, getError(), "Testing error.", "they should be equal")
}

func TestUnsetError(t *testing.T) {
	unsetError()

	assert.Equal(t, getError(), "Testing error.", "they should be equal")
}

func TestIsImageTrue(t *testing.T) {
	result := isImage("image/jpg")
	if result == false {
		t.Error("It is supposed true.")
	}
}

func TestIsImageFalse(t *testing.T) {
	result := isImage("text/jpg")
	if result {
		t.Error("It is supposed false.")
	}
}
