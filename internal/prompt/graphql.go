package prompt

import (
	"github.com/mykelswitzer/progenitor/v2/pkg/config"
)

func SetupGraphql(cfg *config.Config) error {
	return boolPrompt("Do you want a graphql interface", config.CFG_GQL_REQ, cfg)
}
