package cmd

import (
  "log"
  "os"
  "path/filepath"

  "github.com/mykelswitzer/progenitor/internal/prompt"
  "github.com/mykelswitzer/progenitor/internal/scaffolding"
  "github.com/mykelswitzer/progenitor/internal/terraform"
  "github.com/mykelswitzer/progenitor/pkg/aws"
  "github.com/mykelswitzer/progenitor/pkg/config"
  "github.com/urfave/cli/v2"
)

var (
  awsClient *aws.Client
)

var prompts = map[string][]func(*config.Config) error{
  "go-grpc": {
    prompt.ProjectTeam,
    prompt.ProjectName,
    prompt.ProjectDir,
    prompt.SetupGraphql,
    prompt.UseDB,
    prompt.CoreDBObject,
    prompt.UseReporting,
    prompt.RunTerraform,
  },
}

func Execute(cfg *config.Config) {

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
    Commands: []*cli.Command{
      {
        Name:  "go-grpc",
        Usage: "scaffolds a gRPC service in Go",
        Action: func(c *cli.Context) error {
          cfg.Set("projectType", "go-grpc")
          return generate(cfg)
        },
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Println(err)
  }

}

func generate(cfg *config.Config) error {

  for _, p := range prompts[cfg.GetString(config.CFG_PRJ_TYPE)] {
    if err := p(cfg); err != nil {
      return handleError(err)
    }
  }

  err := setupRepo(cfg)
  if err != nil {
    log.Println(err.Error())
    return err
  }

  scaffold, err := scaffolding.New(cfg)
  if err != nil {
    log.Println(err.Error())
    return err
  }

  if err = scaffold.BuildStructure(); err != nil {
    log.Println(err.Error())
    return err
  }

  if err = scaffold.BuildFiles(); err != nil {
    log.Println(err.Error())
    return err
  }

  if err = commitCodeToRepo(cfg, scaffold); err != nil {
    log.Println(err.Error())
    return err
  }

  // moved to here, note that this assumes that all projects will store
  // the terraform code in a /terraform folder in the root directory
  // of the project... while this appears true at this time, it may not
  // be in the future. This change enables committing code to the repo
  // independent of the success of terraform running... which was previously
  // breaking the code. There is probably a better long term fix, which we can
  // invest in if it continues to create issues
  if scaffold.Config.GetBool("runTerraform") {
    base, err := os.Getwd()
    if err != nil {
      log.Println(err.Error())
      return err
    }
    tfDir := filepath.Join(base, scaffold.Config.GetString(config.CFG_PRJ_DIR), "terraform")

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
