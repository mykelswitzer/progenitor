package prompt

import (
    "github.com/mykelswitzer/progenitor/v2/pkg/config"
)

// RunTerraform generates a prompt asking the user if they would like Progenitor to go
// ahead and run the rendered project's Terraform plan
func RunTerraform(cfg *config.Config) error {
    return boolPrompt("Would you like to run the project's Terraform plan?", config.CFG_TF_RUN, cfg)
}
