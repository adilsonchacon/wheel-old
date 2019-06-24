package singles

var ViewContent = `import (
	"time"
	"{{ .AppDomain }}/commons/app/view"
	"{{ .AppDomain }}/db/entities"
)

type PaginationJson struct {
	Pagination view.MainPagination ` + "`" + `json:"pagination"` + "`" + `
	{{ .EntityNames.Plural }} []Json ` + "`" + `json:"{{ .EntityNames.SnakeCase }}"` + "`" + `
}

type SuccessfullySavedJson struct {
	SystemMessage view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	{{ .EntityNames.Name }} Json ` + "`" + `json:"{{ .EntityNames.SnakeCase }}"` + "`" + `
}

type Json struct {
	ID uint ` + "`" + `json:"id"` + "`" + `
  {{- range .EntityColumns }}
  {{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .NameSnakeCase }}"` + "`" + `
  {{- end }}
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func SetJson({{ .EntityNames.LowerCamelCase }} entities.{{ .EntityNames.Name }}) Json {
	return Json{
		ID: {{ .EntityNames.LowerCamelCase }}.ID,
    {{- range .EntityColumns }}
    {{ .Name }}: {{ $.EntityNames.LowerCamelCase }}.{{ .Name }},
    {{- end }}
		CreatedAt: {{ .EntityNames.LowerCamelCase }}.CreatedAt,
		UpdatedAt: {{ .EntityNames.LowerCamelCase }}.UpdatedAt,
	}
}`
