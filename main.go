package progenitor

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/internal/terraform"

	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
	"github.com/mykelswitzer/progenitor/pkg/scaffold"

	"github.com/urfave/cli/v2"
)

func cliCommands(cfg *config.Config, scaffolds scaffold.Scaffolds) []*cli.Command {
  var commands []*cli.Command
  for _, s := range scaffolds {
    cmd := &cli.Command{
      Name:   s.GetName(),
      Usage:  s.GetDescription(),
      Action: func(c *cli.Context) error {
        return generate(cfg, s)
      },
    }
    commands = append(commands, cmd)
  }
  return commands
}

func Execute(cfg *config.Config, scaffolds scaffold.Scaffolds) {

	app := &cli.App{
		Name: "progenitor",
		Usage: `
      Hello, I am the Progenitor!!!
      Please answer my questions, and
      I will set up a nice set of
      boilerplate code, so that you
      do not need to do that awful
      copy pasta you used to do.
    `,
		Commands: cliCommands(cfg, scaffolds),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}

}

// default prompts we should ask for all generated projects
var defaultPrompts []prompt.PromptFunc = []prompt.PromptFunc{
  prompt.ProjectName,
  prompt.ProjectDir,
}

func buildPrompts(scaffoldPrompts []prompt.PromptFunc) []prompt.PromptFunc {
	prompts := defaultPrompts
	prompts = append(prompts, scaffoldPrompts...)
	return prompts
}

func generate(cfg *config.Config, s scaffold.ScaffoldDS) error {

	var err error
  
  cfg.Set("projectType", s.GetName())
	
  prompts := buildPrompts(s.GetPrompts())
	for _, p := range prompts {
		if err := p(cfg); err != nil {
			return handleError(err)
		}
	}

	err = setupRepo(cfg)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	dir := cfg.GetString(prompt.PRJ_DIR)
	tmplScfld, err := s.Generate(cfg, dir, filesys.SetBasePath(dir))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err = tmplScfld.BuildStructure(); err != nil {
		log.Println(err.Error())
		return err
	}

	if err = tmplScfld.BuildFiles(); err != nil {
		log.Println(err.Error())
		return err
	}

	if err = commitCodeToRepo(cfg, tmplScfld); err != nil {
		log.Println(err.Error())
		return err
	}

  // NOTE: need to move this into like a plugin
	// moved to here, note that this assumes that all projects will store
	// the terraform code in a /terraform folder in the root directory
	// of the project... while this appears true at this time, it may not
	// be in the future. This change enables committing code to the repo
	// independent of the success of terraform running... which was previously
	// breaking the code. There is probably a better long term fix, which we can
	// invest in if it continues to create issues
	if tmplScfld.Config.GetBool("runTerraform") {
		base, err := os.Getwd()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		tfDir := filepath.Join(base, tmplScfld.Config.GetString(prompt.PRJ_DIR), "terraform")

		if err := terraform.Run(tfDir); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

func handleError(err error) error {
	log.Println(err)
	os.Exit(1)
	return err
}
