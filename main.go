package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"wheel.smart26.com/generator"
	// "wheel.smart26.com/strcase"
)

func IsGoInstalled() bool {
	var out bytes.Buffer

	cmd := exec.Command("go", "version")
	cmd.Stdout = &out
	err := cmd.Run()

	return err == nil
}

func CheckDependences() {
	var out bytes.Buffer
	var hasDependence bool

	fmt.Println("Checking dependences...")

	cmd := exec.Command("go", "list", "...")
	cmd.Stdout = &out

	cmd.Run()
	// if err := cmd.Run(); err != nil {
	//   fmt.Println("error:", err)
	// }

	installedDependences := strings.Split(out.String(), "\n")

	requiredDependences := []string{"github.com/jinzhu/gorm", "gopkg.in/yaml.v2", "github.com/gorilla/mux", "github.com/dgrijalva/jwt-go", "github.com/satori/go.uuid", "github.com/lib/pq", "golang.org/x/crypto/bcrypt"}
	for _, requiredDependence := range requiredDependences {
		hasDependence = false
		for _, installedDependence := range installedDependences {
			hasDependence = (requiredDependence == installedDependence)
			if hasDependence {
				break
			}
		}

		if !hasDependence {
			fmt.Printf("package %s was not found, installing...\n", requiredDependence)
			cmd := exec.Command("go", "get", requiredDependence)
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("package %s was successfully installed\n", requiredDependence)
			}
		} else {
			fmt.Printf("package %s was found\n", requiredDependence)
		}
	}
}

func main() {
	var options = make(map[string]string)

	command := os.Args[1]
	fmt.Println("command:", command)

	if !IsGoInstalled() {
		log.Fatal("\"Go\" seems not installed")
	} else {
		fmt.Println("\"Go\" seems installed")
		CheckDependences()
	}

	// TODO: if command == new ... else if ... else ...
	// TODO: check options

	if command == "new" {
		options["app_name"] = os.Args[2]
		options["app_domain"] = os.Args[2] + "." + os.Args[3]

		generator.NewApp(options)
	} else {
		generator.Single(os.Args[2], []string{"title:string:001", "description:text:002"})
	}

}

/*

wheel new APP_NAME

- check whether go is installed, if don't alert and stop running
- check whether wheel's go packages are installed, for each if don't install it
- check "HOME_DIR -> go" folder exists, otherwise create it

*/

/*
  -h, [--help]        # Show this help message and quit
  -v, [--version]     # Show Wheel version number and quit
*/

/*
  -G, [--skip-git]    # Skip .gitignore file
*/
