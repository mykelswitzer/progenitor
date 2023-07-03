package cmd

import (
	"testing"

	"github.com/mykelswitzer/progenitor/pkg/prompt"
)

func TestBuildPrompts(t *testing.T) {

	sPrompts := make([]prompt.PromptFunc, 0)

	result := buildPrompts(sPrompts)
	if len(result) != len(defaultPrompts) {
		t.Errorf("Result was incorrect, expected length of: %v, received: %v.", len(defaultPrompts), len(result))
	}
}
