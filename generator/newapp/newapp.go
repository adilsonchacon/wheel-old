package newapp

import (
	"path/filepath"
	"wheel.smart26.com/generator/common"
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

var templateVar common.TemplateVar
var rootAppPath string
var appTemplatesPath string

func Generate(options map[string]string) {
	// Main vars
	templateVar = common.TemplateVar{
		AppName:   options["app_name"],
		AppDomain: options["app_domain"],
		SecretKey: common.SecureRandom(128),
	}
	rootAppPath = common.BuildRootAppPath(options["app_domain"])
	appTemplatesPath = filepath.Join("templates", "baseapp")

	// APP Root path
	common.CreateRootAppPath(rootAppPath)

	// APP Handlers
	common.GeneratePathAndFileFromPackage(rootAppPath, handlers.MyselfPath, handlers.MyselfContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, handlers.SessionPath, handlers.SessionContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, handlers.UserPath, handlers.UserContent, templateVar, false)

	// APP myself
	common.GeneratePathAndFileFromPackage(rootAppPath, myself.ViewPath, myself.ViewContent, templateVar, false)

	// APP session
	common.GeneratePathAndFileFromPackage(rootAppPath, session.ModelPath, session.ModelContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, session.ViewPath, session.ViewContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, sessionmailer.PasswordRecoveryEnPath, sessionmailer.PasswordRecoveryEnContent, templateVar, true)
	common.GeneratePathAndFileFromPackage(rootAppPath, sessionmailer.PasswordRecoveryPtBrPath, sessionmailer.PasswordRecoveryPtBrContent, templateVar, true)
	common.GeneratePathAndFileFromPackage(rootAppPath, sessionmailer.SignUpEnPath, sessionmailer.SignUpEnContent, templateVar, true)
	common.GeneratePathAndFileFromPackage(rootAppPath, sessionmailer.SignUpPtBrPath, sessionmailer.SignUpPtBrContent, templateVar, true)

	// APP user
	common.GeneratePathAndFileFromPackage(rootAppPath, usertemplate.ModelPath, usertemplate.ModelContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, usertemplate.ViewPath, usertemplate.ViewContent, templateVar, false)

	// COMMONS APPs
	common.GeneratePathAndFileFromPackage(rootAppPath, handler.Path, handler.Content, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, model.Path, model.Content, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, pagination.Path, pagination.Content, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, searchengine.Path, searchengine.Content, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, view.Path, view.Content, templateVar, false)

	// COMMONS conversor
	common.GeneratePathAndFileFromPackage(rootAppPath, conversor.Path, conversor.Content, templateVar, false)

	// COMMONS crypto
	common.GeneratePathAndFileFromPackage(rootAppPath, crypto.Path, crypto.Content, templateVar, false)

	// COMMONS locale
	common.GeneratePathAndFileFromPackage(rootAppPath, locale.Path, locale.Content, templateVar, false)

	// COMMONS log
	common.GeneratePathAndFileFromPackage(rootAppPath, logtemplate.Path, logtemplate.Content, templateVar, false)

	// COMMONS mailer
	common.GeneratePathAndFileFromPackage(rootAppPath, mailer.Path, mailer.Content, templateVar, false)

	// config
	common.GeneratePathAndFileFromPackage(rootAppPath, config.Path, config.Content, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, config.AppPath, config.AppContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, config.DatabasePath, config.DatabaseContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, config.EmailPath, config.EmailContent, templateVar, true)

	// config keys
	common.GeneratePathAndFileFromPackage(rootAppPath, configkeys.RsaExamplePath, configkeys.RsaExampleContent, templateVar, true)
	common.GeneratePathAndFileFromPackage(rootAppPath, configkeys.RsaPubExamplePath, configkeys.RsaPubExampleContent, templateVar, true)

	// config locales
	common.GeneratePathAndFileFromPackage(rootAppPath, configlocales.EnPath, configlocales.EnContent, templateVar, true)
	common.GeneratePathAndFileFromPackage(rootAppPath, configlocales.PtBrPath, configlocales.PtBrContent, templateVar, true)

	// db
	common.GeneratePathAndFileFromPackage(rootAppPath, entities.SessionPath, entities.SessionContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, entities.UserPath, entities.UserContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, schema.Path, schema.Content, templateVar, false)

	// routes
	common.GeneratePathAndFileFromPackage(rootAppPath, routes.AuthorizePath, routes.AuthorizeContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, routes.MiddlewarePath, routes.MiddlewareContent, templateVar, false)
	common.GeneratePathAndFileFromPackage(rootAppPath, routes.Path, routes.Content, templateVar, false)

	// main
	common.GeneratePathAndFileFromPackage(rootAppPath, templates.MainPath, templates.MainContent, templateVar, false)
}
