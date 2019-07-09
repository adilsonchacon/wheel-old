package config

var Path = []string{"config", "config.go"}

var Content = `package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"{{ .AppRepository }}/commons/log"
)

type AppConfig struct {
	App_name                          string
	App_repository                        string
	Secret_key                        string
	Reset_password_expiration_seconds int
	Reset_password_url                string
	Token_expiration_seconds          int
	Locales                           []string
}

var appConfig AppConfig

func AppName() string {
	return appConfig.App_name
}

func AppRepository() string {
	return appConfig.App_repository
}

func SecretKey() string {
	return appConfig.Secret_key
}

func ResetPasswordExpirationSeconds() int {
	return appConfig.Reset_password_expiration_seconds
}

func ResetPasswordUrl() string {
	return appConfig.Reset_password_url
}

func TokenExpirationSeconds() int {
	return appConfig.Token_expiration_seconds
}

func Locales() []string {
	return appConfig.Locales
}

func readAppConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/app.yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}

func init() {
	err := yaml.Unmarshal(readAppConfigFile(), &appConfig)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}
}`
