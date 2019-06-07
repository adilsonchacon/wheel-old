package schema

var Path = []string{"db", "schema", "migrate.go"}

var Content = `package schema

import (
	"{{ .AppDomain }}/app/user"
	"{{ .AppDomain }}/commons/app/model"
	"{{ .AppDomain }}/commons/crypto"
	"{{ .AppDomain }}/db/entities"
)

func Migrate() {
	model.Db.AutoMigrate(&entities.User{})

	_, err := user.FindByEmail("{{ .UserEmail }}")
	if err != nil {
 		model.Db.Create(&entities.User{Name: "{{ .UserName }}", Email: "{{ .UserEmail }}", Password: crypto.SetPassword("{{ .UserPassword }}"), Locale: "en", Admin: true})
	}

	model.Db.AutoMigrate(&entities.Session{})
}`