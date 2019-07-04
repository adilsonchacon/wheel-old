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
	gencommon.GenerateFromTemplateFullPath(rootAppPath, handlers.MyselfPath, handlers.MyselfContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, handlers.SessionPath, handlers.SessionContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, handlers.UserPath, handlers.UserContent, templateVar, false)

	// APP myself
	gencommon.GenerateFromTemplateFullPath(rootAppPath, myself.ViewPath, myself.ViewContent, templateVar, false)

	// APP session
	gencommon.GenerateFromTemplateFullPath(rootAppPath, session.ModelPath, session.ModelContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, session.ViewPath, session.ViewContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, sessionmailer.PasswordRecoveryEnPath, sessionmailer.PasswordRecoveryEnContent, templateVar, true)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, sessionmailer.PasswordRecoveryPtBrPath, sessionmailer.PasswordRecoveryPtBrContent, templateVar, true)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, sessionmailer.SignUpEnPath, sessionmailer.SignUpEnContent, templateVar, true)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, sessionmailer.SignUpPtBrPath, sessionmailer.SignUpPtBrContent, templateVar, true)

	// APP user
	gencommon.GenerateFromTemplateFullPath(rootAppPath, usertemplate.ModelPath, usertemplate.ModelContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, usertemplate.ViewPath, usertemplate.ViewContent, templateVar, false)

	// COMMONS APPs
	gencommon.GenerateFromTemplateFullPath(rootAppPath, handler.Path, handler.Content, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, model.Path, model.Content, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, pagination.Path, pagination.Content, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, searchengine.Path, searchengine.Content, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, view.Path, view.Content, templateVar, false)

	// COMMONS conversor
	gencommon.GenerateFromTemplateFullPath(rootAppPath, conversor.Path, conversor.Content, templateVar, false)

	// COMMONS crypto
	gencommon.GenerateFromTemplateFullPath(rootAppPath, crypto.Path, crypto.Content, templateVar, false)

	// COMMONS locale
	gencommon.GenerateFromTemplateFullPath(rootAppPath, locale.Path, locale.Content, templateVar, false)

	// COMMONS log
	gencommon.GenerateFromTemplateFullPath(rootAppPath, logtemplate.Path, logtemplate.Content, templateVar, false)

	// COMMONS mailer
	gencommon.GenerateFromTemplateFullPath(rootAppPath, mailer.Path, mailer.Content, templateVar, false)

	// config
	gencommon.GenerateFromTemplateFullPath(rootAppPath, config.Path, config.Content, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, config.AppPath, config.AppContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, config.DatabasePath, config.DatabaseContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, config.EmailPath, config.EmailContent, templateVar, true)

	// config keys
	gencommon.GenerateFromTemplateFullPath(rootAppPath, configkeys.RsaExamplePath, configkeys.RsaExampleContent, templateVar, true)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, configkeys.RsaPubExamplePath, configkeys.RsaPubExampleContent, templateVar, true)

	// config locales
	gencommon.GenerateFromTemplateFullPath(rootAppPath, configlocales.EnPath, configlocales.EnContent, templateVar, true)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, configlocales.PtBrPath, configlocales.PtBrContent, templateVar, true)

	// db
	gencommon.GenerateFromTemplateFullPath(rootAppPath, entities.SessionPath, entities.SessionContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, entities.UserPath, entities.UserContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, schema.Path, schema.Content, templateVar, false)

	// routes
	gencommon.GenerateFromTemplateFullPath(rootAppPath, routes.AuthorizePath, routes.AuthorizeContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, routes.MiddlewarePath, routes.MiddlewareContent, templateVar, false)
	gencommon.GenerateFromTemplateFullPath(rootAppPath, routes.Path, routes.Content, templateVar, false)

	// main
	gencommon.GenerateFromTemplateFullPath(rootAppPath, templates.MainPath, templates.MainContent, templateVar, false)
}
