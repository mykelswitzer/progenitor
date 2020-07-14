package cmd


import (
	_ "errors"
	_ "io/ioutil"
	_ "os"
)
import "github.com/manifoldco/promptui"

func promptDb() (bool, error) {

	output := map[string]bool {"Yes":true,"No":false}
	var keys []string
	for k := range output {
	    keys = append(keys, k)
	}
  
	prompt := promptui.Select{
		Label: "Do you need a database?",
		Items: keys,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return output[result], err

}

// func setupDbMigrateDir(path string) {

// 	  // Check if file already exists
//   if _, err := os.Stat(path); err == nil {
//     return true
//   }

//   // Attempt to create it
//   var d []byte
//   if err := ioutil.WriteFile(fp, d, 0644); err == nil {
//     os.Remove(fp) // And delete it
//     return true
//   }

//   return false
// }