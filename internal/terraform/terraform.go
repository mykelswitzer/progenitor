package terraform

import (
	"context"
	"log"
	"os/exec"
)

import "github.com/caring/go-packages/pkg/errors"
import "github.com/hashicorp/terraform-exec/tfexec"

// TODO: Find a better way to organize these functions. Perhaps make them part of an interface

// isTerraformInstalled checks if Terraform is installed by searching for it in the directories named
//by the PATH environment variable. If found the path of the Terraform executable is returned. If it's
// not found, an error is returned.
func isTerraformInstalled() (string, error) {
	path, err := exec.LookPath("terraform")

	if err != nil {
		return "", errors.Wrap(err, "Terraform not found in PATH")
	}
	return path, nil
}

// TfRun chains together all the steps to run the newly generated project's Terraform
func TfRun(tfDir string) error {
	awsEnvs := []string{"caring-prod", "caring-stg", "caring-dev"}
	installedPath, err := isTerraformInstalled()
	if err != nil {
		return errors.New("Could not find Terraform installed on PATH")
	}

	tf, err := tfexec.NewTerraform(tfDir, installedPath)
	log.Println("Initializing Terraform directory!")
	if tf == nil {
		return errors.New("Error encountered when initializing Terraform!")
	}
	err = tf.Init(context.Background())
	if err != nil {
		return err
	}

	for _, s := range awsEnvs {
		log.Println("Creating Terraform workspace: ", s)
		err := tf.WorkspaceNew(context.Background(), s)
		if err != nil {
			return err
		}
	}

	log.Println("Applying Terraform plan to 'caring-dev' environment")
	err = tf.WorkspaceSelect(context.Background(), awsEnvs[2])
	if err != nil {
		return err
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return err
	}

	return nil
}
