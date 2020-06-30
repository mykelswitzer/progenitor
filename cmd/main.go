package cmd

import "fmt"
import "github.com/urfave/cli/v2"
import "os"


func Execute() {
	app := &cli.App{
    Name: "progenitor",
    Usage: `
             @@@@,             
           (@@@@@@@           
 ,##%.     (@@@@@@@,     /###    Hello, I am the Progenitor!!!
 ######(     @@@@@    *###### 
    ######          *#####       Please answer my questions, and
       ####/       ####(         I will set up a nice set of 
        ,####    .####           boilerplate code, so that you 
          ####  /###/            do not need to do that awful
          .###  ####             copy pasta you used to do.
           ###%,###           
           .###.###           
           ###/ ###.          
          .###  ####          
          ###,   ###,          `,
   Action: func(c *cli.Context) error {

   		reponame, err := promptReponame()
   		if err != nil {
   			fmt.Println(err.Error())
   			return nil
   		}

   		createRepo(reponame)


      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    fmt.Println(err)
  }
}
