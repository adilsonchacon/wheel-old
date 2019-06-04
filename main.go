package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type TemplateVars struct {
	AppDomain string
	AppName   string
}

func GenerateHandlers(templatePath string, destinyPath string, templateVars TemplateVars) {
	var content bytes.Buffer

	fmt.Println("templatePath:", templatePath)
	fmt.Println("destinyPath:", destinyPath)
	fmt.Println("templateVars:", templateVars)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&content, &templateVars)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(destinyPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(content.String())
	if err != nil {
		panic(err)
	}

	f.Sync()
}

func main() {
	var out bytes.Buffer
	var hasDependence bool

	command := os.Args[1]
	fmt.Println("command:", command)
	appName := os.Args[2]
	domain := os.Args[3]
	appDomain := appName + "." + domain

	templateVars := TemplateVars{AppName: appName, AppDomain: appDomain}

	// TODO: if command == new ... else if ... else ...

	cmd := exec.Command("go", "version")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		log.Fatal("\"Go\" seems not installed")
	}

	cmd = exec.Command("go", "list", "...")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		log.Fatal("\"Go\" seems not installed")
	}
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

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ignore:", dir)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	rootAppPath := filepath.Join(usr.HomeDir, "go", "src", appName)
	if err := os.MkdirAll(rootAppPath, 0775); err != nil {
		log.Fatal(err)
	}

	// APP Handlers
	handlersDestinyBasePath := filepath.Join(rootAppPath, "app", "handlers")
	if err := os.MkdirAll(handlersDestinyBasePath, 0775); err != nil {
		log.Fatal(err)
	}
	handlersTemplateBasePath := filepath.Join(".", "templates", "baseapp", "app", "handlers")
	defaultHandlers := []string{"myself", "session", "user"}
	for _, handler := range defaultHandlers {
		GenerateHandlers(filepath.Join(handlersTemplateBasePath, handler+"_handler.template"), filepath.Join(handlersDestinyBasePath, handler+"_handler.go"), templateVars)
	}

	// APP myself
	if err := os.MkdirAll(filepath.Join(rootAppPath, "app", "myself"), 0775); err != nil {
		log.Fatal(err)
	}

	// APP session
	if err := os.MkdirAll(filepath.Join(rootAppPath, "app", "session"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "app", "session", "mailer"), 0775); err != nil {
		log.Fatal(err)
	}

	// APP user
	if err := os.MkdirAll(filepath.Join(rootAppPath, "app", "user"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "app", "handler"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "app", "model"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "app", "model", "pagination"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "app", "model", "searchengine"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "app", "view"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "conversor"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "crypto"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "locale"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "log"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "commons", "mailer"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "config"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "config", "keys"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "config", "locales"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "db"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "db", "entities"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "db", "schema"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "log"), 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootAppPath, "routes"), 0775); err != nil {
		log.Fatal(err)
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
