package locale

var Path = []string{"commons", "locale", "locale.go"}

var Content = `package locale

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"{{ .AppRepository }}/commons/log"
)

type Keys struct {
	Welcome                        string
	Password_recovery_instructions string
}

var locales Keys

func Welcome() string {
	return locales.Welcome
}

func PasswordRecoveryInstructions() string {
	return locales.Password_recovery_instructions
}

func Load(language string) {
	err := yaml.Unmarshal(readLocaleFile(language), &locales)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}
}

func readLocaleFile(language string) []byte {
	data, err := ioutil.ReadFile("./config/locales/" + language + ".yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}`
