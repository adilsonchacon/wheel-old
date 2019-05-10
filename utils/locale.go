package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type localeKeys struct {
	Welcome                        string
	Password_recovery_instructions string
}

var locales localeKeys

func LocaleWelcome() string {
	return locales.Welcome
}

func LocalePasswordRecoveryInstructions() string {
	return locales.Password_recovery_instructions
}

func LocaleLoad(language string) {
	err := yaml.Unmarshal(readLocaleFile(language), &locales)
	if err != nil {
		LoggerError().Fatalf("error: %v", err)
	}
}

func readLocaleFile(language string) []byte {
	data, err := ioutil.ReadFile("./config/locales/" + language + ".yml")
	if err != nil {
		LoggerError().Fatal(err)
	}

	return data
}
