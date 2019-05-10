package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"wheel.smart26.com/app/views"
)

func ApiRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API root")
}

func Error404(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(views.SetNotFoundErrorMessage())
}

func normalizeUrlQueryParams(param string, mapParams map[string][]string) map[string]string {
	var criteria, value string

	query := make(map[string]string)

	filter, err1 := regexp.Compile(param + `\[[a-zA-Z0-9\-\_]+\](\[\]){0,1}`)
	if err1 != nil {
		fmt.Println(err1)
	}

	removePrefix, err2 := regexp.Compile(param + `\[`)
	if err2 != nil {
		fmt.Println(err2)
	}

	removeSufix, err3 := regexp.Compile(`\](\[\]){0,1}`)
	if err3 != nil {
		fmt.Println(err3)
	}

	for key := range mapParams {
		if filter.MatchString(key) {
			criteria = key
			criteria = removeSufix.ReplaceAllString(criteria, "")
			criteria = removePrefix.ReplaceAllString(criteria, "")
			value = strings.Join(mapParams[key], ",")

			query[criteria] = value
		}
	}

	return query
}
