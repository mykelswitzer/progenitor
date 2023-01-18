package scaffolding

import (
	"log"

	"github.com/caring/go-packages/pkg/errors"
	gogrpc "github.com/mykelswitzer/progenitor-tmpl-go-grpc"
	"github.com/mykelswitzer/progenitor/v2/internal/filesys"
	"github.com/mykelswitzer/progenitor/v2/pkg/config"
	"github.com/mykelswitzer/progenitor/v2/pkg/scaffold"
)

var scaffoldingTypes = map[string]scaffold.ScaffoldDS{
	"go-grpc": gogrpc.GoGrpc{},
}

// New will return a scaffold based on the project type
func New(cfg *config.Config) (*scaffold.Scaffold, error) {

	projectType := cfg.GetString(config.CFG_PRJ_TYPE)
	if scfld, ok := scaffoldingTypes[projectType]; ok {
		dir := cfg.GetString(config.CFG_PRJ_DIR)
		scaffold, err := scfld.Init(cfg, dir, filesys.SetBasePath(dir))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return scaffold, nil
	}

	return nil, errors.New("project scaffold missing: " + projectType)

}
