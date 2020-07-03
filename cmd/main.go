package cmd

import (
  "os"
  "log"
)

import "github.com/urfave/cli/v2"
import "github.com/caring/progenitor/pkg/aws"

var (
  awsClient *aws.Client
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



      var region string = "us-east-1"
      var account_id string = "182565773517"
      var role string = "ops-mgmt-admin"

      awsClient := aws.New()
      awsClient.SetConfig(&region, &account_id, &role)

   		reponame, err := promptReponame()
   		if err != nil {
   			log.Println(err.Error())
        os.Exit(1)
   		}

      token, err := awsClient.GetSecret("github_token")
      if err != nil {
        log.Println(err.Error())
        os.Exit(1)
      }

   		createRepo(*token.SecretString, reponame)


      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Println(err)
  }

}

