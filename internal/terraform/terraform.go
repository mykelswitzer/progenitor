package terraform

import (
    "log"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

import "github.com/caring/go-packages/pkg/errors"

// TODO: Find a better way to organize these functions. Perhaps make them part of an interface

// isTerraformInstalled checks if Terraform is installed by searching for it
// in the directories named by the PATH environment variable. If its not found
// an error is returned, otherwise nil is returned.
func isTerraformInstalled() error {
    _, err := exec.LookPath("terraform")

    if err != nil {
        return errors.Wrap(err, "Terraform not found in PATH")
    }
    return nil
}
// TODO: Add timeout and better OS signal handling
// tfGetVersion executes the command `terraform version` and
// returns the version parsed from the output
func tfGetVersion() ([]byte, error) {
    tf := exec.Command("terraform", "version")
    out, err := tf.Output()

    if err != nil {
        return nil, errors.Wrap(err, "Error encountered while get Terraform version!")
    }

    pattern := regexp.MustCompile(`v[0-9]\.[0-9]+\.[0-9]+`)
    return pattern.Find(out), nil
}

// TODO: Add timeout and better OS signal handling
// tfInit initializes the Terraform directory of the newly
// created project repository. The var tfDir is the absolute path
// of the Terraform directory.
func tfInit(tfDir string) error {
    tf := exec.Command("terraform", "init")
    tf.Dir = tfDir
    tf.Stdout = os.Stdout
    tf.Stderr = os.Stdout
    err := tf.Run()

    if err != nil {
        return errors.Wrap(err, "Error encountered while initializing Terraform directory!")
    }
    return nil
}

// TODO: Add timeout and better OS signal handling
// tfNewWorkspace adds a new workspace.
// tfDir: the absolute path of the Terraform directory to run the command in
// name: the name of the new Terraform workspace
func tfNewWorkspace(tfDir string, name string) error {
    tf := exec.Command("terraform", "workspace", "new", name)
    tf.Dir = tfDir
    tf.Stdout = os.Stdout
    tf.Stderr = os.Stdout
    err := tf.Run()

    if err != nil {
        return errors.Wrap(err, "Error encountered while creating Terraform workspace!")
    }
    return nil
}

// tfGetWorkspace returns the current Terraform workspace
// tfDir: the absolute path of the Terraform directory to run the command in
func tfGetWorkspace(tfDir string) (string, error) {
    tf := exec.Command("terraform", "workspace", "show")
    tf.Dir = tfDir
    out, err := tf.Output()

    if err != nil {
        return "", errors.Wrap(err, "Error encountered while getting active Terraform workspace name!")
    }

    if len(out) < 1 {
        return "", errors.New("Could not get active Terraform workspace name!")
    }
    return string(out), nil
}

// TODO: Add timeout and better OS signal handling
// tfSelectWorkspace set the current workspace to the supplied workspace name
// tfDir: the absolute path of the Terraform directory to run the command in
// name: the name of the Terraform workspace to select
func tfSelectWorkspace(tfDir string, name string) error {
    tf := exec.Command("terraform", "workspace", "select", name)
    tf.Dir = tfDir
    tf.Stdout = os.Stdout
    tf.Stderr = os.Stdout
    err := tf.Run()

    if err != nil {
        return errors.Wrap(err, "Error encountered while selecting Terraform workspace!")
    }
    return nil
}

// TODO: Add timeout and better OS signal handling
// tfPlan generates a plan via the command 'terraform plan'
// tfDir: the absolute path of the Terraform directory to run the command in
func tfPlan(tfDir string) (string, error) {
    tf := exec.Command("terraform", "plan")
    tf.Dir = tfDir
    tf.Stdout = os.Stdout
    tf.Stderr = os.Stdout
    out, err := tf.Output()

    if err != nil {
       return "", errors.Wrap(err, "Error encountered while running terraform plan!")
    }

    if len(out) < 1 {
        return "", errors.New("Terraform failed to generate a plan!")
    }

    plan := string(out)
    return strings.TrimSpace(plan), nil
}

// TODO: Add timeout and better OS signal handling
// tfApply runs the Terraform plan for the newly generated project
// tfDir: the absolute path of the Terraform directory to run the command in
func tfApply(tfDir string) error {
    tf := exec.Command("terraform", "apply", "-auto-approve")
    tf.Dir = tfDir
    tf.Stdout = os.Stdout
    tf.Stderr = os.Stdout
    err := tf.Run()

    if err != nil {
        return errors.Wrap(err, "Error encountering while applying Terraform plan!")
    }

    return nil
}

// TfRun chains together all the steps to run the newly generated project's Terraform
func TfRun(tfDir string) error {
    awsEnvs := []string{"caring-prod", "caring-stg", "caring-dev"}

    log.Println("Initializing Terraform directory!")
    err := tfInit(tfDir)
    if err != nil {
        return err
    }

    //log.Println("Creating Terraform  workspace: caring-prod.")
    //err = tfNewWorkspace(tfDir, "caring-prod")
    //if err != nil {
    //    return err
    //}
    //
    //log.Println("Creating Terraform workspace: caring-stg.")
    //err = tfNewWorkspace(tfDir, "caring-stg")
    //if err != nil {
    //    return err
    //}
    //
    //log.Println("Creating Terraform workspace: caring-dev.")
    //err = tfNewWorkspace(tfDir, "caring-dev")
    //if err != nil {
    //    return err
    //}

    for _, s := range awsEnvs {
       log.Println("Creating Terraform workspace: ", s)
       err := tfNewWorkspace(tfDir, s)
       if err != nil {
           return err
       }
    }

    log.Println("Applying Terraform plan to 'caring-dev' environment")
    err = tfSelectWorkspace(tfDir, "caring-dev")
    if err != nil {
        return err
    }

    err = tfApply(tfDir)
    if err != nil {
        return err
    }
    return nil
}
