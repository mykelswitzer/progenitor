package prompt

import (
	"github.com/caring/progenitor/internal/config"
)

func SetupGraphql(config *config.Config) error {
	return boolPrompt("Do you want a graphql interface", "gqlRequired", config)
}
