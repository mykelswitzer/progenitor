package cmd

import "fmt"
import "github.com/urfave/cli/v2"
import "os"

import "github.com/caring/progenitor/pkg/aws"

var (
  awsClient    aws.Client
)

func init() {

  var region string = "us-east-1"
  var account_id string = "182565773517"
  var role string = "ops-mgmt-admin"

  awsClient := aws.Client{}
  awsClient.SetConfig(&region, &account_id, &role)

}


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




   		reponame, err := promptReponame()
   		if err != nil {
   			fmt.Println(err.Error())
   			return nil
   		}

   		fmt.Println(reponame)


      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    fmt.Println(err)
  }
}