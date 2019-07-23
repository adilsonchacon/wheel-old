package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/adilsonchacon/wheel/commons/notify"
	"github.com/adilsonchacon/wheel/generator"
	"github.com/adilsonchacon/wheel/help"
	"github.com/adilsonchacon/wheel/version"
	"os"
	"os/exec"
	"strings"
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

	notify.Simpleln("Checking dependences...")

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
			notify.Simple(fmt.Sprintf("package %s was not found, installing...", requiredDependence))
			cmd := exec.Command("go", "get", requiredDependence)
			cmd.Stdout = &out
			err := cmd.Run()
			notify.FatalIfError(err)

			notify.Simpleln(fmt.Sprintf("package %s was successfully installed", requiredDependence))
		} else {
			notify.Simpleln(fmt.Sprintf("package %s was found", requiredDependence))
		}
	}
}

func handleNewApp(args []string) {
	var options = make(map[string]string)

	preOptions := strings.Split(os.Args[2], "/")

	options["app_name"] = preOptions[len(preOptions)-1]
	options["app_repository"] = os.Args[2]

	generator.NewApp(options)
}

func buildGenerateOptions(args []string) (map[string]bool, error) {
	var options = make(map[string]bool)
	var subject string
	var err error

	subject = args[2]

	options["model"] = false
	options["entity"] = false
	options["view"] = false
	options["handler"] = false
	options["routes"] = false
	options["migrate"] = false
	options["authorize"] = false

	switch subject {
	case "scaffold":
		options["model"] = true
		options["entity"] = true
		options["view"] = true
		options["handler"] = true
		options["routes"] = true
		options["migrate"] = true
		options["authorize"] = true
		if len(args) < 4 {
			err = errors.New("invalid scaffold name")
		}
	case "model":
		options["model"] = true
		options["entity"] = true
		if len(args) < 4 {
			err = errors.New("invalid model name")
		}
	case "handler":
		options["handler"] = true
	case "view":
		options["view"] = true
	case "entity":
		options["entity"] = true
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
	notify.FatalIfError(err)

	handleGenerateNewCrud(args, options)
}

func handleHelp() {
	notify.Simpleln(help.Content)
}

func handleVersion() {
	notify.Simpleln(version.Content)
}

func main() {
	command := os.Args[1]

	if !IsGoInstalled() {
		notify.FatalIfError(errors.New("\"Go\" seems not installed"))
	} else {
		notify.Simpleln("\"Go\" seems installed")
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
