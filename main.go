package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"wheel.smart26.com/copy"
	"wheel.smart26.com/templates"
	"wheel.smart26.com/templates/app/handlers"
	"wheel.smart26.com/templates/app/myself"
	"wheel.smart26.com/templates/app/session"
	"wheel.smart26.com/templates/app/usertemplate"
	"wheel.smart26.com/templates/commons/app/handler"
	"wheel.smart26.com/templates/commons/app/model"
	"wheel.smart26.com/templates/commons/app/model/pagination"
	"wheel.smart26.com/templates/commons/app/model/searchengine"
	"wheel.smart26.com/templates/commons/app/view"
	"wheel.smart26.com/templates/commons/conversor"
	"wheel.smart26.com/templates/commons/crypto"
	"wheel.smart26.com/templates/commons/locale"
	"wheel.smart26.com/templates/commons/logtemplate"
	"wheel.smart26.com/templates/commons/mailer"
	"wheel.smart26.com/templates/config"
	"wheel.smart26.com/templates/db/entities"
	"wheel.smart26.com/templates/db/schema"
	"wheel.smart26.com/templates/routes"
)

type TemplateVars struct {
	AppDomain string
	AppName   string
	SecretKey string
}

var templateVars TemplateVars
var rootAppPath string
var appTemplatesPath string

func saveTextFile(content string, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}

	f.Sync()

}

func GenerateFromTemplateFile(templatePath string, destinyPath string) {
	var content bytes.Buffer

	fmt.Println("created:", destinyPath)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&content, &templateVars)
	if err != nil {
		log.Fatal(err)
	}

	saveTextFile(content.String(), destinyPath)
}

func GenerateFromTemplateString(content string, destinyPath string) {
	var buffContent bytes.Buffer

	fmt.Println("created:", destinyPath)

	tmpl, err := template.New("T").Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&buffContent, templateVars)
	if err != nil {
		log.Fatal(err)
	}

	saveTextFile(buffContent.String(), destinyPath)
}

func HandlePackagesPathInfo(path []string) (string, string) {
	var basePath, fileName string

	for index, value := range path {
		if index+1 != len(path) {
			basePath = filepath.Join(basePath, value)
		} else {
			fileName = value
		}
	}

	return basePath, fileName
}

func GeneratePathAndFileFromPackage(path []string, content string) {
	basePath, fileName := HandlePackagesPathInfo(path)

	if err := os.MkdirAll(filepath.Join(rootAppPath, basePath), 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", basePath)
	}

	GenerateFromTemplateString(content, filepath.Join(rootAppPath, basePath, fileName))
}

func CopyFiles(sourcePath string, files []string, destinyPath string) {
	if err := os.MkdirAll(destinyPath, 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", destinyPath)
	}

	for _, file := range files {
		if err := copy.File(filepath.Join(sourcePath, file), filepath.Join(destinyPath, file)); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("created:", filepath.Join(destinyPath, file))
		}
	}
}

func SecureRandom(size int) string {
	var letters = []rune("0123456789abcdefABCDEF")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

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

	cmd := exec.Command("go", "list", "...")
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
}

func BuildRootAppPath(appDomain string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ignore:", dir)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(usr.HomeDir, "go", "src", appDomain)
}

func main() {
	command := os.Args[1]
	fmt.Println("command:", command)
	appName := os.Args[2]
	domain := os.Args[3]
	appDomain := appName + "." + domain

	// TODO: if command == new ... else if ... else ...

	if !IsGoInstalled() {
		log.Fatal("\"Go\" seems not installed")
	} else {
		fmt.Println("\"Go\" seems installed")
		fmt.Println("Checking dependences...")
		CheckDependences()
	}

	// Main vars
	templateVars = TemplateVars{AppName: appName, AppDomain: appDomain, SecretKey: SecureRandom(32)}
	rootAppPath = BuildRootAppPath(appDomain)
	appTemplatesPath = filepath.Join("templates")

	if err := os.MkdirAll(rootAppPath, 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", rootAppPath)
	}

	// APP Handlers
	GeneratePathAndFileFromPackage(handlers.MyselfPath, handlers.MyselfContent)
	GeneratePathAndFileFromPackage(handlers.SessionPath, handlers.SessionContent)
	GeneratePathAndFileFromPackage(handlers.UserPath, handlers.UserContent)

	// APP packages path
	GeneratePathAndFileFromPackage(myself.ViewPath, myself.ViewContent)
	GeneratePathAndFileFromPackage(session.ModelPath, session.ModelContent)
	GeneratePathAndFileFromPackage(session.ViewPath, session.ViewContent)
	GeneratePathAndFileFromPackage(usertemplate.ModelPath, usertemplate.ModelContent)
	GeneratePathAndFileFromPackage(usertemplate.ViewPath, usertemplate.ViewContent)

	// APP session
	sessionMailerSourcePath := filepath.Join(appTemplatesPath, "app", "session", "mailer")
	sessionMailerFiles := []string{"password_recovery.en.html", "password_recovery.pt-BR.html", "sign_up.en.html", "sign_up.pt-BR.html"}
	sessionMailerDestinyPath := filepath.Join(rootAppPath, "app", "session", "mailer")
	CopyFiles(sessionMailerSourcePath, sessionMailerFiles, sessionMailerDestinyPath)

	// COMMONS APPs
	GeneratePathAndFileFromPackage(handler.Path, handler.Content)
	GeneratePathAndFileFromPackage(model.Path, model.Content)
	GeneratePathAndFileFromPackage(pagination.Path, pagination.Content)
	GeneratePathAndFileFromPackage(searchengine.Path, searchengine.Content)
	GeneratePathAndFileFromPackage(view.Path, view.Content)

	// COMMONS conversor
	GeneratePathAndFileFromPackage(conversor.Path, conversor.Content)

	// COMMONS crypto
	GeneratePathAndFileFromPackage(crypto.Path, crypto.Content)

	// COMMONS locale
	GeneratePathAndFileFromPackage(locale.Path, locale.Content)

	// COMMONS log
	GeneratePathAndFileFromPackage(logtemplate.Path, logtemplate.Content)

	// COMMONS mailer
	GeneratePathAndFileFromPackage(mailer.Path, mailer.Content)

	// config
	GeneratePathAndFileFromPackage(config.Path, config.Content)

	configKeysSourcePath := filepath.Join(appTemplatesPath, "config", "keys")
	configKeysFiles := []string{"app.key.rsa.example", "app.key.rsa.pub.example"}
	configKeysDestinyPath := filepath.Join(rootAppPath, "config", "keys")
	CopyFiles(configKeysSourcePath, configKeysFiles, configKeysDestinyPath)

	localeSourcePath := filepath.Join(appTemplatesPath, "config", "locales")
	localeFiles := []string{"en.yml", "pt-BR.yml"}
	localeDestinyPath := filepath.Join(rootAppPath, "config", "locales")
	CopyFiles(localeSourcePath, localeFiles, localeDestinyPath)

	confingSourcePath := filepath.Join(appTemplatesPath, "config")
	confingFiles := []string{"database.example.yml", "email.example.yml"}
	confingDestinyPath := filepath.Join(rootAppPath, "config")
	CopyFiles(confingSourcePath, confingFiles, confingDestinyPath)

	appConfigSourcePath := filepath.Join(appTemplatesPath, "config", "app.example.yml")
	appConfigDestinyPath := filepath.Join(rootAppPath, "config", "app.example.yml")
	GenerateFromTemplateFile(appConfigSourcePath, appConfigDestinyPath)

	// db
	GeneratePathAndFileFromPackage(entities.SessionPath, entities.SessionContent)
	GeneratePathAndFileFromPackage(entities.UserPath, entities.UserContent)
	GeneratePathAndFileFromPackage(schema.Path, schema.Content)

	// routes
	GeneratePathAndFileFromPackage(routes.AuthorizePath, routes.AuthorizeContent)
	GeneratePathAndFileFromPackage(routes.MiddlewarePath, routes.MiddlewareContent)
	GeneratePathAndFileFromPackage(routes.Path, routes.Content)

	// main
	GeneratePathAndFileFromPackage(templates.MainPath, templates.MainContent)
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
