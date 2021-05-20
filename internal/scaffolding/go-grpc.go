package scaffolding

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/caring/progenitor/internal/strings"
)

type goGrpcTemplateData struct {
	Team         string
	Name         string
	LocalPath    string
	UseDB        bool
	DBModel      string
	UseReporting bool
}

// Init sets the values for the goGrpcTemplateData struct
func (t goGrpcTemplateData) Init(config *config.Config) templateData {
	t.Team = config.GetString("projectTeam")
	t.Name = config.GetString("projectName")
	t.LocalPath = config.GetString("projectDir")
	t.UseDB = config.GetBool("dbRequired")
	t.DBModel = config.GetString("dbModel")
	t.UseReporting = config.GetBool("reportingRequired")
	return t
}

type goGrpc struct{}

// GetData fetched the template data needed to populate
// the go templates
func (g goGrpc) GetData(config *config.Config) templateData {
	td := goGrpcTemplateData{}
	return td.Init(config)
}

// Init will populate a Scaffold with all relevant data
// for the scaffolding to run for this service type
func (g goGrpc) Init(cfg *config.Config) (*Scaffold, error) {

	dir := cfg.GetString("projectDir")

	grpcProject := Scaffold{
		Source:        g,
		Config:        cfg,
		BaseDir:       Dir{Name: dir},
		TemplatePath:  "go-grpc",
		SkipTemplates: []string{},
		Fs:            SetBasePath(dir),
		ProcessHooks: map[string]func(*Scaffold) error{
			"postBuildFiles": postBuildFiles,
		},
	}

	if cfg.GetBool("dbRequired") == false {
		grpcProject.SkipTemplates = append(grpcProject.SkipTemplates, "terraform/rds.tf.tmpl")
	}

	apiDir := Dir{Name: "api"}
	apiDir.AddSubDirs(Dir{Name: "pb"})

	if cfg.GetBool("gqlRequired") == true {
		apiDir.AddSubDirs(Dir{Name: "graphql"})
	}

	cmdDir := Dir{Name: "cmd"}
	cmdDir.AddSubDirs(Dir{Name: "server"})

	internalDir := Dir{Name: "internal"}
	if cfg.GetBool("dbRequired") == true {
		dbDir := Dir{Name: "db"}
		dbDir.AddSubDirs(Dir{Name: "migrations"})

		internalDir.AddSubDirs(dbDir)
	}

	domainsDir := Dir{Name: "domain"}
	domainsDir.AddSubDirs(Dir{Name: "domain"})

	internalDir.AddSubDirs(domainsDir, Dir{Name: "handlers"}, Dir{Name: "service"})

	pkgDir := Dir{Name: "pkg"}
	pkgDir.AddSubDirs(Dir{Name: "client"})

	tfDir := Dir{Name: "terraform"}
	tfDir.AddSubDirs(Dir{Name: "templates"})

	grpcProject.BaseDir.AddSubDirs(apiDir, cmdDir, internalDir, pkgDir, Dir{Name: ".circleci"}, tfDir)

	return &grpcProject, nil
}

///// Scaffold Methods /////

// postBuildFiles will run the protoc to build the proto-generated
// go code
func postBuildFiles(s *Scaffold) error {

	if err := renameServiceFiles(s); err != nil {
		return err
	}

	if err := generateProto(s); err != nil {
		return err
	}

	return nil
}

// we have a few files named generically that we need to rename
func renameServiceFiles(s *Scaffold) error {

	base, err := os.Getwd()
	path := filepath.Join(base, s.Config.GetString("projectDir"))

	oldName := filepath.Join(path, "api/pb/service.proto")
	newName := filepath.Join(path, "api/pb", strings.ToPackage(s.Config.GetString("projectName"))+".proto")
	err = os.Rename(oldName, newName)
	if err != nil {
		log.Println(err)
	}

	if s.Config.GetBool("gqlRequired") == true {
		oldName = filepath.Join(path, "api/graphql/service.graphqls")
		newName = filepath.Join(path, "api/graphql", strings.ToPackage(s.Config.GetString("projectName"))+".graphqls")
		err = os.Rename(oldName, newName)
		if err != nil {
			log.Println(err)
		}
	}

	oldName = filepath.Join(path, "internal/handlers/handlers.go")
	newName = filepath.Join(path, "internal/handlers", s.Config.GetString("projectName")+".go")
	err = os.Rename(oldName, newName)
	if err != nil {
		log.Println(err)
	}

	if s.Config.GetBool("dbRequired") == true {
		oldName = filepath.Join(path, "internal/domain/domain/service_test.go")
		newName = filepath.Join(path, "internal/domain/domain", s.Config.GetString("dbModel")+"_test.go")
		err = os.Rename(oldName, newName)
		if err != nil {
			log.Println(err)
		}

		oldName = filepath.Join(path, "internal/domain/domain/")
		newName = filepath.Join(path, "internal/domain/", s.Config.GetString("dbModel"))
		err = os.Rename(oldName, newName)
		if err != nil {
			log.Println(err)
		}
	}

	return nil

}

// runs the protoc go transpiler
func generateProto(s *Scaffold) error {

	base, err := os.Getwd()
	if err != nil {
		log.Println("Relative path provided, unable to determine root.")
		os.Exit(1)
	}
	genProtoPath := "api/pb"
	path := filepath.Join(base, s.Config.GetString("projectDir"), genProtoPath)

	if err := os.Chmod(path, 0777); err != nil {
		log.Println("Unable to chmod filepath to generate proto, please manually run", err)
		return nil
	}

	executable := filepath.Join(path, "gen_proto.sh")
	if err := os.Chmod(executable, 0755); err != nil {
		log.Println("Unable to chmod filepath to generate proto, please manually run", err)
		return nil
	}

	cmd := &exec.Cmd{
		Path:   executable,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	if err := cmd.Run(); err != nil {
		log.Println("Unable to execute bashscript to generate proto, please manually run", err)
	}

	return nil

}
