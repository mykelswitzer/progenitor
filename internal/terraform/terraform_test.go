package terraform

import (
    "testing"
)

// TODO: Find a way to mock this so the test can run in any environment
// Verifies that the isTerraformInstalled function doesn't return an error if
// Terraform is installed on the host running the test. Its not exactly
// a unit test since it depends on the host running it.
func Test_isTerraformInstalled(t *testing.T) {
    _, err := isTerraformInstalled()

    if err != nil {
        t.Fatal("Expected nil but got error instead")
    }
}

// TODO: Find a better method of testing that doesn't hard code version
// Verifies that the getVersion function returns the installed
// version of Terraform. Its not exactly a unit test since it depends on
// on the host running it.
func Test_tfGetVersion(t *testing.T)  {
    version, err := getVersion()

    if err != nil {
       t.Fatal("Unexpected error encountered!")
    }

    if version == nil {
        t.Fatal("No version found!")
    }

    if string(version) != "v0.12.29" {
       t.Log(string(version))
       t.Fatal("Incorrect version returned!")
    }
}
