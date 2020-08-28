package repo

import "testing"

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the isTerraformInstalled function doesn't return an error if
// Terraform is installed on the host running the test. Its not exactly
// a unit test since it depends on the host running it.
func Test_isTerraformInstalled(t *testing.T) {
    err := isTerraformInstalled()

    if err != nil {
        t.Fatal("Expected nil but got error instead")
    }
}

// TODO: Find better method of testing that doesn't hard code version
// Verifies the getTerraformVersion function returns the installed
// version of Terraform. Its not exactly a unit test since it depends on
// on the host running it.
func Test_getTerraformVersion(t *testing.T)  {
    version, err := getTerraformVersion()

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
