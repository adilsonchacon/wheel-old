package newapp

import (
	"wheel.smart26.com/generator/gencommon"
	"wheel.smart26.com/templates/templateapp"
	"wheel.smart26.com/templates/templateapp/app/handlers"
	"wheel.smart26.com/templates/templateapp/app/myself"
	"wheel.smart26.com/templates/templateapp/app/session"
	"wheel.smart26.com/templates/templateapp/app/session/sessionmailer"
	"wheel.smart26.com/templates/templateapp/app/usertemplate"
	"wheel.smart26.com/templates/templateapp/commons/app/handler"
	"wheel.smart26.com/templates/templateapp/commons/app/model"
	"wheel.smart26.com/templates/templateapp/commons/app/model/pagination"
	"wheel.smart26.com/templates/templateapp/commons/app/model/searchengine"
	"wheel.smart26.com/templates/templateapp/commons/app/view"
	"wheel.smart26.com/templates/templateapp/commons/conversor"
	"wheel.smart26.com/templates/templateapp/commons/crypto"
	"wheel.smart26.com/templates/templateapp/commons/locale"
	"wheel.smart26.com/templates/templateapp/commons/logtemplate"
	"wheel.smart26.com/templates/templateapp/commons/mailer"
	"wheel.smart26.com/templates/templateapp/config"
	"wheel.smart26.com/templates/templateapp/config/configkeys"
	"wheel.smart26.com/templates/templateapp/config/configlocales"
	"wheel.smart26.com/templates/templateapp/db/entities"
	"wheel.smart26.com/templates/templateapp/db/schema"
	"wheel.smart26.com/templates/templateapp/routes"
)

var templateVar gencommon.TemplateVar
var rootAppPath string

func prependRootAppPathToPath(path []string) []string {
	return append([]string{rootAppPath}, path...)
}

func Generate(options map[string]string) {
	// Main vars
	templateVar = gencommon.TemplateVar{
		AppName:   options["app_name"],
		AppDomain: options["app_domain"],
		SecretKey: gencommon.SecureRandom(128),
	}
	rootAppPath = gencommon.BuildRootAppPath(options["app_domain"])

	// APP Root path
	gencommon.CreateRootAppPath(rootAppPath)

	// APP Handlers
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.MyselfPath), handlers.MyselfContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.SessionPath), handlers.SessionContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.UserPath), handlers.UserContent, templateVar)

	// APP myself
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(myself.ViewPath), myself.ViewContent, templateVar)

	// APP session
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ModelPath), session.ModelContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ViewPath), session.ViewContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryEnPath), sessionmailer.PasswordRecoveryEnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryPtBrPath), sessionmailer.PasswordRecoveryPtBrContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpEnPath), sessionmailer.SignUpEnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpPtBrPath), sessionmailer.SignUpPtBrContent, templateVar)

	// APP user
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ModelPath), usertemplate.ModelContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ViewPath), usertemplate.ViewContent, templateVar)

	// COMMONS APPs
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handler.Path), handler.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.Path), model.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(pagination.Path), pagination.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(searchengine.Path), searchengine.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(view.Path), view.Content, templateVar)

	// COMMONS conversor
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(conversor.Path), conversor.Content, templateVar)

	// COMMONS crypto
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(crypto.Path), crypto.Content, templateVar)

	// COMMONS locale
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(locale.Path), locale.Content, templateVar)

	// COMMONS log
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(logtemplate.Path), logtemplate.Content, templateVar)

	// COMMONS mailer
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(mailer.Path), mailer.Content, templateVar)

	// config
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.Path), config.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.AppPath), config.AppContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.DatabasePath), config.DatabaseContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(config.EmailPath), config.EmailContent, templateVar)

	// config keys
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configkeys.RsaExamplePath), configkeys.RsaExampleContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configkeys.RsaPubExamplePath), configkeys.RsaPubExampleContent, templateVar)

	// config locales
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.EnPath), configlocales.EnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.PtBrPath), configlocales.PtBrContent, templateVar)

	// db
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.SessionPath), entities.SessionContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.UserPath), entities.UserContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(schema.Path), schema.Content, templateVar)

	// routes
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.AuthorizePath), routes.AuthorizeContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.MiddlewarePath), routes.MiddlewareContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.Path), routes.Content, templateVar)

	// main
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(templates.MainPath), templates.MainContent, templateVar)
}
