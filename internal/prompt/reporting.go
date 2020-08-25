package prompt

import (
    "github.com/caring/progenitor/internal/config"
)

// UseReporting generates a prompt asking the user if the service needs application metrics reporting functionality.
// This user input determines whether or not to render the Terraform plan so that it will provision resources used
// for application metrics reporting.
func UseReporting(config *config.Config) error {
    return boolPrompt("Do you need application metrics reporting?", "reportingRequired", config)
}
