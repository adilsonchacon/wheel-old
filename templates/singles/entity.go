package singles

var EntityContent = `package entities

import (
	"time"
)

type {{ .EntityNames.Name }} struct {
	ID uint ` + "`" + `gorm:"primary_key"` + "`" + `
  {{- range .EntityColumns }}
  {{ .Name }} {{ .Type }} {{ .Extras }}
  {{- end }}
	CreatedAt time.Time
	UpdatedAt time.Time
}`
