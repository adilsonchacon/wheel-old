package handler

var Path = []string{"commons", "app", "handler", "handler.go"}

var Content = `package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"{{ .AppDomain }}/commons/app/view"
	"{{ .AppDomain }}/commons/log"
)

func ApiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API root")
}

func Error404(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("handler: Error404")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
}

func NormalizeUrlQueryParams(param string, mapParams map[string][]string) map[string]string {
	var criteria, value string
	var checkParam, removePrefix, removeSufix *regexp.Regexp
	var err error

	query := make(map[string]string)

	checkParam, err = regexp.Compile(param + ` + "`" + `\[[a-zA-Z0-9\-\_]+\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removePrefix, err = regexp.Compile(param + ` + "`" + `\[` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removeSufix, err = regexp.Compile(` + "`" + `\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	for key := range mapParams {
		if checkParam.MatchString(key) {
			criteria = key
			criteria = removeSufix.ReplaceAllString(criteria, "")
			criteria = removePrefix.ReplaceAllString(criteria, "")
			value = strings.Join(mapParams[key], ",")

			query[criteria] = value
		}
	}

	return query
}`
