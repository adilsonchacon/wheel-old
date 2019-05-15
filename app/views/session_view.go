package views

import (
	"bytes"
	"html/template"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/config"
)

type SessionSignInSuccess struct {
	Message SystemMessage `json:"system_message"`
	Token   string        `json:"token"`
	Expires int           `json:"expires"`
}

type SessionSignOutSuccess struct {
	Message SystemMessage `json:"system_message"`
}

type SessionSignUpSuccess struct {
	UserFirstName string
	AppName       string
}

type SessionPasswordRecoveryInstructions struct {
	UserFirstName          string
	LinkToPasswordRecovery string
}

func SessionSignInSuccessMessage(mType string, content string, token string) SessionSignInSuccess {
	return SessionSignInSuccess{Message: SystemMessage{mType, content}, Token: token, Expires: config.TokenExpirationSeconds()}
}

func SessionSignOutSuccessMessage(mType string, content string) SessionSignOutSuccess {
	return SessionSignOutSuccess{Message: SystemMessage{mType, content}}
}

func SessionRefreshSuccessMessage(mType string, content string, token string) SessionSignInSuccess {
	return SessionSignInSuccessMessage(mType, content, token)
}

func SessionSignUpSuccessMessage(mType string, content string, token string) SessionSignInSuccess {
	return SessionSignInSuccessMessage(mType, content, token)
}

func SessionSignUpMailer(user *models.User) string {
	var content bytes.Buffer

	data := SessionSignUpSuccess{UserFirstName: models.UserFirstName(user), AppName: config.AppName()}

	tmpl, err := template.ParseFiles("./app/views/mailer/sessions/sign_up." + user.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}

func SessionPasswordRecoveryInstructionsMailer(user *models.User, token string) string {
	var content bytes.Buffer

	data := SessionPasswordRecoveryInstructions{UserFirstName: models.UserFirstName(user), LinkToPasswordRecovery: config.ResetPasswordUrl() + "?token=" + token}

	tmpl, err := template.ParseFiles("./app/views/mailer/sessions/password_recovery." + user.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}
