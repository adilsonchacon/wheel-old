package newapp

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"text/template"
	"time"
	"wheel.smart26.com/templates/baseapp"
	"wheel.smart26.com/templates/baseapp/app/handlers"
	"wheel.smart26.com/templates/baseapp/app/myself"
	"wheel.smart26.com/templates/baseapp/app/session"
	"wheel.smart26.com/templates/baseapp/app/session/sessionmailer"
	"wheel.smart26.com/templates/baseapp/app/usertemplate"
	"wheel.smart26.com/templates/baseapp/commons/app/handler"
	"wheel.smart26.com/templates/baseapp/commons/app/model"
	"wheel.smart26.com/templates/baseapp/commons/app/model/pagination"
	"wheel.smart26.com/templates/baseapp/commons/app/model/searchengine"
	"wheel.smart26.com/templates/baseapp/commons/app/view"
	"wheel.smart26.com/templates/baseapp/commons/conversor"
	"wheel.smart26.com/templates/baseapp/commons/crypto"
	"wheel.smart26.com/templates/baseapp/commons/locale"
	"wheel.smart26.com/templates/baseapp/commons/logtemplate"
	"wheel.smart26.com/templates/baseapp/commons/mailer"
	"wheel.smart26.com/templates/baseapp/config"
	"wheel.smart26.com/templates/baseapp/config/configkeys"
	"wheel.smart26.com/templates/baseapp/config/configlocales"
	"wheel.smart26.com/templates/baseapp/db/entities"
	"wheel.smart26.com/templates/baseapp/db/schema"
	"wheel.smart26.com/templates/baseapp/routes"
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

func GeneratePathAndFileFromPackage(path []string, content string, skipTemplate bool) {
	basePath, fileName := HandlePackagesPathInfo(path)

	if err := os.MkdirAll(filepath.Join(rootAppPath, basePath), 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", basePath)
	}

	if skipTemplate {
		saveTextFile(content, filepath.Join(rootAppPath, basePath, fileName))
	} else {
		GenerateFromTemplateString(content, filepath.Join(rootAppPath, basePath, fileName))
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

func Generate(options map[string]string) {
	// Main vars
	templateVars = TemplateVars{AppName: options["app_name"], AppDomain: options["app_domain"], SecretKey: SecureRandom(64)}
	rootAppPath = BuildRootAppPath(options["app_domain"])
	appTemplatesPath = filepath.Join("templates", "baseapp")

	if err := os.MkdirAll(rootAppPath, 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", rootAppPath)
	}

	// APP Handlers
	GeneratePathAndFileFromPackage(handlers.MyselfPath, handlers.MyselfContent, false)
	GeneratePathAndFileFromPackage(handlers.SessionPath, handlers.SessionContent, false)
	GeneratePathAndFileFromPackage(handlers.UserPath, handlers.UserContent, false)

	// APP myself
	GeneratePathAndFileFromPackage(myself.ViewPath, myself.ViewContent, false)

	// APP session
	GeneratePathAndFileFromPackage(session.ModelPath, session.ModelContent, false)
	GeneratePathAndFileFromPackage(session.ViewPath, session.ViewContent, false)
	GeneratePathAndFileFromPackage(sessionmailer.PasswordRecoveryEnPath, sessionmailer.PasswordRecoveryEnContent, true)
	GeneratePathAndFileFromPackage(sessionmailer.PasswordRecoveryPtBrPath, sessionmailer.PasswordRecoveryPtBrContent, true)
	GeneratePathAndFileFromPackage(sessionmailer.SignUpEnPath, sessionmailer.SignUpEnContent, true)
	GeneratePathAndFileFromPackage(sessionmailer.SignUpPtBrPath, sessionmailer.SignUpPtBrContent, true)

	// APP user
	GeneratePathAndFileFromPackage(usertemplate.ModelPath, usertemplate.ModelContent, false)
	GeneratePathAndFileFromPackage(usertemplate.ViewPath, usertemplate.ViewContent, false)

	// COMMONS APPs
	GeneratePathAndFileFromPackage(handler.Path, handler.Content, false)
	GeneratePathAndFileFromPackage(model.Path, model.Content, false)
	GeneratePathAndFileFromPackage(pagination.Path, pagination.Content, false)
	GeneratePathAndFileFromPackage(searchengine.Path, searchengine.Content, false)
	GeneratePathAndFileFromPackage(view.Path, view.Content, false)

	// COMMONS conversor
	GeneratePathAndFileFromPackage(conversor.Path, conversor.Content, false)

	// COMMONS crypto
	GeneratePathAndFileFromPackage(crypto.Path, crypto.Content, false)

	// COMMONS locale
	GeneratePathAndFileFromPackage(locale.Path, locale.Content, false)

	// COMMONS log
	GeneratePathAndFileFromPackage(logtemplate.Path, logtemplate.Content, false)

	// COMMONS mailer
	GeneratePathAndFileFromPackage(mailer.Path, mailer.Content, false)

	// config
	GeneratePathAndFileFromPackage(config.Path, config.Content, false)
	GeneratePathAndFileFromPackage(config.AppPath, config.AppContent, false)
	GeneratePathAndFileFromPackage(config.DatabasePath, config.DatabaseContent, false)
	GeneratePathAndFileFromPackage(config.EmailPath, config.EmailContent, true)

	// config keys
	GeneratePathAndFileFromPackage(configkeys.RsaExamplePath, configkeys.RsaExampleContent, true)
	GeneratePathAndFileFromPackage(configkeys.RsaPubExamplePath, configkeys.RsaPubExampleContent, true)

	// config locales
	GeneratePathAndFileFromPackage(configlocales.EnPath, configlocales.EnContent, true)
	GeneratePathAndFileFromPackage(configlocales.PtBrPath, configlocales.PtBrContent, true)

	// db
	GeneratePathAndFileFromPackage(entities.SessionPath, entities.SessionContent, false)
	GeneratePathAndFileFromPackage(entities.UserPath, entities.UserContent, false)
	GeneratePathAndFileFromPackage(schema.Path, schema.Content, false)

	// routes
	GeneratePathAndFileFromPackage(routes.AuthorizePath, routes.AuthorizeContent, false)
	GeneratePathAndFileFromPackage(routes.MiddlewarePath, routes.MiddlewareContent, false)
	GeneratePathAndFileFromPackage(routes.Path, routes.Content, false)

	// main
	GeneratePathAndFileFromPackage(templates.MainPath, templates.MainContent, false)
}
