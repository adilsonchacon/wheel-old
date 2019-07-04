package templatecrud

var HandlerContent = `package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"{{ .AppDomain }}/app/{{ .EntityName.LowerCase }}"
	"{{ .AppDomain }}/commons/app/handler"
	"{{ .AppDomain }}/commons/app/model"
	"{{ .AppDomain }}/commons/app/view"
	"{{ .AppDomain }}/commons/log"
	"{{ .AppDomain }}/db/entities"
)

func {{ .EntityName.CamelCase }}Create(w http.ResponseWriter, r *http.Request) {
	var new{{ .EntityName.CamelCase }} = entities.{{ .EntityName.CamelCase }}{}

	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Create")
	w.Header().Set("Content-Type", "application/json")

	{{ .EntityName.LowerCamelCase }}SetParams(&new{{ .EntityName.CamelCase }}, r)

	valid, errs := {{ .EntityName.LowerCase }}.Create(&new{{ .EntityName.CamelCase }})
	log.Debug.Println(errs)

	if valid {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCamelCase }}.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "{{ .EntityName.SnakeCase }} was successfully created"), {{ .EntityName.CamelCase }}: {{ .EntityName.LowerCamelCase }}.SetJson(new{{ .EntityName.CamelCase }})})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not created", errs))
	}
}

func {{ .EntityName.CamelCase }}Update(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Update")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not updated", []error{err}))
		return
	}

	{{ .EntityName.CamelCase }}SetParams(&{{ .EntityName.LowerCamelCase }}Current, r)

	if valid, errs := {{ .EntityName.LowerCase }}.Update(&{{ .EntityName.LowerCamelCase }}Current); valid {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "{{ .EntityName.SnakeCase }} was successfully updated"), {{ .EntityName.CamelCase }}: {{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}Current)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not updated", errs))
	}
}

func {{ .EntityName.CamelCase }}Destroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Destroy")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])

	if err == nil && {{ .EntityName.LowerCase }}.Destroy(&{{ .EntityName.LowerCamelCase }}Current) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "{{ .EntityName.SnakeCase }} was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "{{ .EntityName.SnakeCase }} was not found"))
	}
}

func {{ .EntityName.CamelCase }}Show(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Show")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])

	if err == nil {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}Current))
	} else {
		json.NewEncoder(w).Encode(view.SetSystemMessage("alert", "{{ .EntityName.SnakeCase }} was not found"))
	}
}

func {{ .EntityName.CamelCase }}List(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var {{ .EntityName.LowerCamelCase }}Jsons []{{ .EntityName.LowerCase }}.Json
	var {{ .EntityName.LowerCamelCase }}List []entities.{{ .EntityName.CamelCase }}

	log.Info.Println("Handler: {{ .EntityName.CamelCase }}List")
	w.Header().Set("Content-Type", "application/json")

	normalizedParams := handler.NormalizeUrlQueryParams("search", r.URL.Query())

	{{ .EntityName.LowerCamelCase }}List, page, pages, entries = {{ .EntityName.LowerCase }}.Paginate(normalizedParams, r.FormValue("page"), r.FormValue("per_page"))

	for i = 0; i < len({{ .EntityName.LowerCamelCase }}List); i++ {
		{{ .EntityName.LowerCamelCase }}Jsons = append({{ .EntityName.LowerCamelCase }}Jsons, {{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}List[i]))
	}

	pagination := view.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.PaginationJson{Pagination: pagination, {{ .EntityName.CamelCasePlural }}: {{ .EntityName.LowerCamelCase }}Jsons})
}

func {{ .EntityName.LowerCamelCase }}SetParams({{ .EntityName.LowerCamelCase }}Set *entities.{{ .EntityName.CamelCase }}, r *http.Request) {
  {{- $filteredEntityColumns := filterEntityColumnsNotForeignKeys .EntityColumns }}
	var allowedParams = []string{ {{- range $index, $element := $filteredEntityColumns }} "{{ $element.NameSnakeCase }}" {{- if isNotLastIndex $index $filteredEntityColumns }}, {{- end }} {{- end }} }

	r.ParseMultipartForm()

	for key := range r.Form {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				model.SetColumnValue({{ .EntityName.LowerCamelCase }}Set, allowedParam, r.FormValue(allowedParam))
				break
			}
		}
	}
}`