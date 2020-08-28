package repo

import (
    "os/exec"
    "regexp"
)

import "github.com/caring/go-packages/pkg/errors"

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

func getTerraformVersion() (string, error) {
    tf := exec.Command("terraform", "version")
    out, err := tf.Output()

    if err != nil {
        return "", errors.Wrap(err, "Error encountered while get Terraform version!")
    }

    OutStr := string(out)
    pattern := regexp.MustCompile(`^Terraform (v[0-9]\.[0-9]+\.[0-9]+$)`)
    return pattern.FindString(OutStr), nil
}
