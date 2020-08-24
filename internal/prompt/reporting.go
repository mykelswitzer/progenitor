package prompt

import (
    "github.com/caring/go-packages/pkg/errors"
    "github.com/caring/progenitor/internal/config"
    "github.com/manifoldco/promptui"
)

// UseReporting generates a prompt asking the user if the service needs application metrics reporting functionality.
// This user input determines whether or not to render the Terraform plan so it will provision resources used
// application metrics reporting.
func UseReporting(config *config.Config) error {
    output := map[string]bool{"Yes": true, "No": false}
    var keys []string

    for k := range output {
        keys = append(keys, k)
    }

    prompt := promptui.Select{
        Label: "Do you need application metrics reporting?",
        Items: keys,
    }
    _, result, err := prompt.Run()

    if err != nil {
        return errors.Wrap(err, "Error in exporting application reporting prompt")
    }

    config.Set("reportingRequired", output[result])
    return nil
}
