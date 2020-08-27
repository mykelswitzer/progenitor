package repo

import "os/exec"
import "context"
import "time"

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
    ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Microsecond)
    defer cancel()

    tf := exec.CommandContext(ctx, "terraform", "version")
    out, err := tf.Output()

    if err != nil {
        return "", errors.Wrap(err, "Error encountered while get Terraform version!")
    }

    out_str := string(out)

}
