package terraform

import (
    "os"
    "strings"
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

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the initTf function successfully executes the command
// `terraform init` inside the Terraform directory of the newly create
// repository
func Test_tfInit(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")
    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TF_DIR is not set!")
    }

    err := initTf(tfDir)
    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }
}

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the newWorkspace function successfully executes the command
// `terraform workspace new <name>` inside the Terraform directory fo the
// newly created repsoitory.
func Test_tfNewWorkspace(t *testing.T) {
   tfDir := os.Getenv("TF_DIR")
   if len(tfDir) < 1 {
       t.Fatal("Aborting test, the environment variable TF_DIR is not set!")
   }

   err := newWorkspace(tfDir, "example")
   if err != nil {
       t.Fatal("Unexpected error encountered!")
   }
}

// TODO: Find a way to mock this so the test can be run in any environment
// Verifies the selectWorkspace function successfully changes the Terraform
// workspace
func Test_tfSelectWorkspace(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")

    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TD_DIR is not set!")
    }

    err := selectWorkspace(tfDir, "caring-dev")
    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }
}

// TODO: Find a way to mock this so the test can be run in any environment
// Verifies the getWorkspace function successfully returns the active
// Terraform workspace.
func Test_tfGetWorkspace(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")

    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TF_DIR is not set!")
    }

    name, err := getWorkspace(tfDir)

    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }

    expectedWorkspace := "caring-dev"
    if strings.Contains(name, expectedWorkspace) == false {
        t.Log(name)
        t.Fatal("Incorrect workspace returned!")
    }
}

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the plan function successfully generates a plan from the
// Terraform files.
func Test_tfPlan(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")

    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TD_DIR is not set!")
    }

    plan, err := plan(tfDir)

    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }

    if len(plan) <1 {
        t.Fatal("Invalid plan returned!")
    }
}

// TODO: Find a way to mock this so the test can run in any environment
// Verifies the apply function successfully applies the Terraform plan
func Test_tfApply(t *testing.T) {
    tfDir := os.Getenv("TF_DIR")

    if len(tfDir) < 1 {
        t.Fatal("Aborting test, the environment variable TF_DIR is not set!")
    }

    err := apply(tfDir)

    if err != nil {
        t.Fatal("Unexpected error encountered!")
    }
}
