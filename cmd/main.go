package cmd

import "fmt"
import "github.com/spf13/cobra"
import "os"
import "github.com/manifoldco/prompui"

func init() {
	cobra.OnInitialize()
}


var rootCmd = &cobra.Command{
  Use:   "progenitor",
  Short: "Progenitor is an interactive service scaffolding system",
  Long: `
             .(/(             
           (///////           
 ,##%.     (///////,     /###    Hello, I am Progrenitor!!!
 ######(     /////    *###### 
    ######          *#####       Please answer my questions, and
       ####/       ####(         I will set up a nice set of 
        ,####    .####           boilerplate code, so that you 
          ####  /###/            do not need to do that awful
          .###  ####             copy and pasting you used to.
           ###%,###           
           .###.###           
           ###/ ###.          
          .###  ####          
          ###,   ###,          `,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}