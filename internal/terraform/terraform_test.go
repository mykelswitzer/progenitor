package repo

import "testing"

// Verifies the isTerraformInstalled function doesn't return an error if
// Terraform is installed on the host running the test. Its not exactly
// a unit test since it depends on the host running it.
func Test_isTerraformInstalled(t *testing.T) {
    err := isTerraformInstalled()

    if err != nil {
        t.Fatal("Expected nil but got error instead")
    }
}

func Test_getTerraformVersion(t *testing.T)  {
    version, err := getTerraformVersion()

    if err != nil {
       t.Fatal("Unexpected error encountered!")
    }

    if version != "v0.12.29" {
       t.Fatal("Incorrect version returned!")
    }
}
