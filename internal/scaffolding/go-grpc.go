package scaffolding

import _ "log"
import "github.com/caring/progenitor/internal/config"

type goGrpc struct{}

func (g goGrpc) Init(cfg *config.Config) (*Scaffold, error) {

	dir := cfg.GetString("projectDir")

	grpcProject := Scaffold{
		Config:       cfg,
		BaseDir:      Dir{Name: dir},
		TemplatePath: "go-grpc",
		Fs:						SetBasePath(dir),
	}

	cmdDir := Dir{Name: "cmd"}
	cmdServerDir := Dir{Name: "server"}
	cmdDir.SubDirs = append(cmdDir.SubDirs, cmdServerDir)
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, cmdDir)

	internalDir := Dir{Name: "internal"}
	if cfg.GetBool("setupDb") == true {
		dbDir := Dir{Name: "db"}
		dbMigrationsDir := Dir{Name: "migrations"}
		dbDir.SubDirs = append(dbDir.SubDirs, dbMigrationsDir)
		internalDir.SubDirs = append(internalDir.SubDirs, dbDir)
	}
	internalDir.SubDirs = append(internalDir.SubDirs, Dir{Name: "handlers"})
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, internalDir)

	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, Dir{Name: "pkg"})

	return &grpcProject, nil

}
