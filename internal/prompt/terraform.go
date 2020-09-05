package prompt

import (
    "github.com/caring/progenitor/internal/config"
)

// RunTerraform generates a prompt asking the user if they would like Progenitor to go
// ahead and run the rendered project's Terraform plan
func RunTerraform (config *config.Config) error {
    return boolPrompt("Would you like to run the project's Terraform plan?", "runTerraform", config)
}
