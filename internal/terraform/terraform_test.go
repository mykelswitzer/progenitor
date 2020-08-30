package repo

import (
    "os"
    "testing"
)

// TODO: Find a way to mock this so the test can run in any environment
// Verifies that the isTerraformInstalled function doesn't return an error if
// Terraform is installed on the host running the test. Its not exactly
// a unit test since it depends on the host running it.
func Test_isTerraformInstalled(t *testing.T) {
    err := isTerraformInstalled()

    if err != nil {
        t.Fatal("Expected nil but got error instead")
    }
}

// TODO: Find a better method of testing that doesn't hard code version
// Verifies that the tfGetVersion function returns the installed
// version of Terraform. Its not exactly a unit test since it depends on
// on the host running it.
func Test_tfGetVersion(t *testing.T)  {
    version, err := tfGetVersion()

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

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the tfInit function successfully executes the command
// `terraform init` inside the Terraform directory of the newly create
// repository
func Test_tfInit(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")
    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TF_DIR not set!")
    }

    err := tfInit(tfDir)
    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }
}

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the tfNewWorkspace function successfully executes the command
// `terraform workspace new <name>` inside the Terraform directory fo the
// newly created repsoitory.
func Test_tfNewWorkspace(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")
    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TF_DIR not set!")
    }

    err := tfNewWorkspace(tfDir, "example")
    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }
}
