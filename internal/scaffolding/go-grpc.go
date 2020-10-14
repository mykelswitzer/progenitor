package scaffolding

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)
import (
	_ "github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/caring/progenitor/internal/terraform"
)

type goGrpcTemplateData struct {
	Name         string
	LocalPath    string
	UseDB        bool
	DBModel      string
	UseReporting bool
}

// Init sets the values for the goGrpcTemplateData struct
func (t goGrpcTemplateData) Init(config *config.Config) templateData {
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

	cmdDir := Dir{Name: "cmd"}
	cmdClientDir := Dir{Name: "client"}
	cmdDir.SubDirs = append(cmdDir.SubDirs, cmdClientDir)
	cmdServerDir := Dir{Name: "server"}
	cmdDir.SubDirs = append(cmdDir.SubDirs, cmdServerDir)
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, cmdDir)

	internalDir := Dir{Name: "internal"}
	if cfg.GetBool("dbRequired") == true {
		dbDir := Dir{Name: "db"}
		dbMigrationsDir := Dir{Name: "migrations"}
		dbDir.SubDirs = append(dbDir.SubDirs, dbMigrationsDir)
		internalDir.SubDirs = append(internalDir.SubDirs, dbDir)
	}
	internalDir.SubDirs = append(internalDir.SubDirs, Dir{Name: "handlers"})
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, internalDir)

	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, Dir{Name: "pb"})
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, Dir{Name: "pkg"})
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, Dir{Name: ".circleci"})

	// Construct and append Terraform directory to project directory
	tfDir := Dir{Name: "terraform"}
	tfTmplDir := Dir{Name: "templates"}
	tfDir.SubDirs = append(tfDir.SubDirs, tfTmplDir)
	grpcProject.BaseDir.SubDirs = append(grpcProject.BaseDir.SubDirs, tfDir)

	return &grpcProject, nil

}

// postBuildFiles will run the protoc to build the proto-generated
// go code
func postBuildFiles(s *Scaffold) error {
	if err := generateProto(s); err != nil {
		return err
	}

	if err := renameServiceFiles(s); err != nil {
		return err
	}

	return nil
}

// we have a few files named generically that we need to rename
// pb/service.pb should be pb/{{config.projectName}}.pb
// internal/db/service*.go should be internal/db/{{config.dbObject}}.go
func renameServiceFiles(s *Scaffold) error {

	base, err := os.Getwd()
	path := filepath.Join(base, s.Config.GetString("projectDir"))

	oldName := filepath.Join(path, "pb/service.pb.go")
	newName := filepath.Join(path, "pb", s.Config.GetString("projectName")+".pb.go")
	err = os.Rename(oldName, newName)
	if err != nil {
		log.Println(err)
	}

	if s.Config.GetBool("dbRequired") == true {

		oldName = filepath.Join(path, "internal/db/service.go")
		newName = filepath.Join(path, "internal/db", s.Config.GetString("dbModel")+".go")
		err = os.Rename(oldName, newName)
		if err != nil {
			log.Println(err)
		}

		oldName = filepath.Join(path, "internal/db/service_test.go")
		newName = filepath.Join(path, "internal/db", s.Config.GetString("dbModel")+"_test.go")
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
	genProtoPath := "pb"
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
