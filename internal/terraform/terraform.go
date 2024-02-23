package terraform

import (
	"context"
	"log"
	"os/exec"

	"github.com/mykelswitzer/progenitor/pkg/config"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
)

// isInstalled checks if Terraform is installed by searching for it in the directories named
// by the PATH environment variable. If found the path of the Terraform executable is returned.
// Otherwise, an error is returned.
func isInstalled() (string, error) {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return "", errors.Wrap(err, "Terraform not found in PATH")
	}
	return path, nil
}

// Run chains together all the steps to run the newly generated project's Terraform
func Run(tfCfg config.TerraformSettings, tfDir string) error {
	installedPath, err := isInstalled()
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

	if tfCfg.IsDefined() {
		// awsEnvs := []string{"mykelswitzer-prod", "mykelswitzer-stg", "mykelswitzer-dev"}
		for _, s := range tfCfg.Workspaces {
			log.Println("Creating Terraform workspace: ", s)
			err := tf.WorkspaceNew(context.Background(), s)
			if err != nil {
				log.Println("WARNING: ", err.Error())
			}
		}

		log.Println("Applying Terraform plan to environment: ", tfCfg.Workspaces[0])
		err = tf.WorkspaceSelect(context.Background(), tfCfg.Workspaces[0])
		if err != nil {
			return err
		}
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return err
	}

	return nil
}
