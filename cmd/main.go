package cmd

import (
	"log"
	"os"
)

import "github.com/urfave/cli/v2"
import "github.com/caring/progenitor/internal/config"
import "github.com/caring/progenitor/internal/scaffolding"
import "github.com/caring/progenitor/pkg/aws"

var (
	awsClient *aws.Client
	cfg       *config.Config
)

func Execute() {
	app := &cli.App{
		Name: "progenitor",
		Usage: `
             @@@@,             
           (@@@@@@@           
 ,##%.     (@@@@@@@,     *###    Hello, I am the Progenitor!!!
 ######*     @@@@@    *#####* 
    ######          *####*       Please answer my questions, and
       ####*       ####*         I will set up a nice set of 
        .####    .####           boilerplate code, so that you 
          ####  .###*            do not need to do that awful
          .###  ####             copy pasta you used to do.
           ###**###           
           .### ###           
           ###* ###.
          .###  ####          
          ###,   ###,          `,
		Action: func(c *cli.Context) error {

			awsClient = setupAwsClient()

			cfg = config.New()

			promptProjectName(cfg)
			promptProjectDir(cfg)

			token, err := awsClient.GetSecret("github_token")
			if err != nil {
				log.Println(err.Error())
				return err
			}

			createRepo(*token.SecretString, cfg)

			setupDb, err := promptDb()
			if err != nil {
				log.Println(err.Error())
				return err
			}
			log.Print(setupDb)

			cfg.Set("projectType", "go-grpc")

			scaffold, err := scaffolding.New(cfg)
			if err != nil {
				log.Println(err.Error())
				return err
			}

			scaffold.BuildStructure()

      scaffold.BuildFiles(*token.SecretString)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}

}

func setupAwsClient() *aws.Client {
	var region string = "us-east-1"
	var account_id string = "182565773517"
	var role string = "ops-mgmt-admin"

	awsClient := aws.New()
	awsClient.SetConfig(&region, &account_id, &role)

	return awsClient
}
