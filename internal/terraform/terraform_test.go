package terraform

import (
	"testing"
)

// TODO: Find a way to mock this so the test can run in any environment
// Verifies that the isTerraformInstalled function doesn't return an error if
// Terraform is installed on the host running the test. Its not exactly
// a unit test since it depends on the host running it.
func Test_isTerraformInstalled(t *testing.T) {
	_, err := isInstalled()

	if err != nil {
		t.Fatal("Expected nil but got error instead")
	}
}
