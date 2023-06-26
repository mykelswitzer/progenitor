package scaffolding

import (
	"log"

	_ "github.com/pkg/errors"
	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/scaffold"
)

// New will return a scaffold based on the project type
func New(cfg *config.Config, scaffoldingTypes scaffold.Scaffolds) (*scaffold.Scaffold, error) {

	projectType := cfg.GetString(config.CFG_PRJ_TYPE)
	scfld, err := scaffoldingTypes.Get(projectType);
	if err != nil {
		return  nil, err
	}

	dir := cfg.GetString(config.CFG_PRJ_DIR)
	scaffold, err := scfld.Init(cfg, dir, filesys.SetBasePath(dir))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return scaffold, nil


}
