package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"wheel.smart26.com/generator"
	"wheel.smart26.com/help"
	"wheel.smart26.com/version"
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

func handleNewApp(args []string) {
	var options map[string]string

	options["app_name"] = os.Args[2]
	options["app_domain"] = os.Args[2] + "." + os.Args[3]

	generator.NewApp(options)
}

func buildGenerateOptions(args []string) (map[string]bool, error) {
	var options = make(map[string]bool)
	var subject string
	var err error

	subject = args[2]

	switch subject {
	case "scaffold":
		options["model"] = true
		options["entity"] = true
		options["view"] = true
		options["handler"] = true
		options["routes"] = true
		options["migrate"] = true
		if len(args) < 4 {
			err = errors.New("invalid scaffold name")
		}
	case "model":
		options["model"] = true
		options["entity"] = true
		options["view"] = false
		options["handler"] = false
		if len(args) < 4 {
			err = errors.New("invalid model name")
		}
	case "handler":
		options["model"] = false
		options["entity"] = false
		options["view"] = false
		options["handler"] = true
	case "view":
		options["model"] = false
		options["entity"] = false
		options["view"] = true
		options["handler"] = false
	case "entity":
		options["model"] = false
		options["entity"] = true
		options["view"] = false
		options["handler"] = false
		if len(args) < 4 {
			err = errors.New("invalid entity name")
		}
	default:
		err = errors.New("invalid generate subject")
	}

	return options, err
}

func handleGenerateNewCrud(args []string, options map[string]bool) {
	var columns []string

	for index, value := range args {
		if index <= 3 {
			continue
		} else {
			columns = append(columns, value)
		}
	}

	generator.NewCrud(args[3], columns, options)
}

func handleGenerate(args []string) {
	var options map[string]bool
	var err error

	options, err = buildGenerateOptions(args)
	if err == nil {
		handleGenerateNewCrud(args, options)
	} else {
		fmt.Println(err)
	}
}

func handleHelp() {
	fmt.Println(help.Content)
}

func handleVersion() {
	fmt.Println(version.Content)
}

func main() {
	command := os.Args[1]
	fmt.Println("command:", command)

	if !IsGoInstalled() {
		log.Fatal("\"Go\" seems not installed")
	} else {
		fmt.Println("\"Go\" seems installed")
		CheckDependences()
	}

	if command == "new" || command == "n" {
		handleNewApp(os.Args)
	} else if command == "generate" || command == "g" {
		handleGenerate(os.Args)
	} else if command == "--help" || command == "-h" {
		handleHelp()
	} else if command == "--version" || command == "-v" {
		handleVersion()
	}
}
