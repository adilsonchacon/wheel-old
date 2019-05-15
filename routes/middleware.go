package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"wheel.smart26.com/app/controllers"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/app/views"
	"wheel.smart26.com/commons/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info.Println(r.Method + ": " + filterUrlValues(r.URL.Path, r.URL.Query()) + " for " + r.RemoteAddr)
		r.ParseForm()
		log.Info.Println("Params: " + filterFormValues(r.Form))
		next.ServeHTTP(w, r)
	})
}

func authorizeMiddleware(next http.Handler) http.Handler {
	var adminUrl = regexp.MustCompile(`^\/user(s){0,1}`)
	var userUrl = regexp.MustCompile(`^\/myself`)
	var sessionRefreshUrl = regexp.MustCompile(`^\/sessions\/refresh`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		if adminUrl.MatchString(r.RequestURI) {
			if authorizeAdmin(token) {
				next.ServeHTTP(w, r)
			} else {
				json.NewEncoder(w).Encode(views.SetNotFoundErrorMessage())
			}
		} else if userUrl.MatchString(r.RequestURI) || sessionRefreshUrl.MatchString(r.RequestURI) {
			if authorizeUser(token) {
				next.ServeHTTP(w, r)
			} else {
				json.NewEncoder(w).Encode(views.SetNotFoundErrorMessage())
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func authorizeAdmin(token string) bool {
	log.Info.Println("Checking token...")
	id, err := controllers.SessionCheck(token)
	if err != nil {
		log.Info.Println("Token is not valid")
		return false
	} else {
		log.Info.Println("Token is valid")
		log.Info.Printf("Checking user %d...\n", id)
		user := models.UserFind(id)
		if models.UserExists(user) {
			if models.UserFind(id).Admin {
				log.Info.Println("User authorized")
				return true
			} else {
				log.Info.Println("User not authorized")
				return false
			}
		} else {
			log.Warn.Println("User does not exist")
			return false
		}
	}
}

func authorizeUser(token string) bool {
	log.Info.Println("Checking token...")
	id, err := controllers.SessionCheck(token)
	if err != nil {
		log.Info.Println("Token is not valid")
		return false
	} else {
		log.Info.Println("Token is valid")
		log.Info.Printf("Checking user %d...\n", id)
		user := models.UserFind(id)
		if models.UserExists(user) {
			models.UserSetCurrent(id)
			log.Info.Println("User authorized")
			return true
		} else {
			log.Warn.Println("User does not exist")
			return false
		}
	}
}

func filterParamsValues(queries map[string][]string) map[string][]string {
	var filter = regexp.MustCompile(`(?i)(password)|(token)`)
	queries_filtered := make(map[string][]string)

	for key := range queries {
		if filter.MatchString(key) {
			queries_filtered[key] = []string{"[FILTERED]"}
		} else {
			queries_filtered[key] = []string{}
			for _, element := range queries[key] {
				queries_filtered[key] = append(queries_filtered[key], element)
			}
		}

	}

	return queries_filtered
}

func filterUrlValues(path string, queries map[string][]string) string {
	var firstParam = true
	queries_filtered := filterParamsValues(queries)

	for key := range queries_filtered {
		if firstParam {
			path = path + "?"
			firstParam = false
		} else {
			path = path + "&"
		}

		path = path + key + "=" + strings.Join(queries_filtered[key], " ")
	}

	return path
}

func filterFormValues(queries map[string][]string) string {
	var buffer bytes.Buffer
	var index int
	queries_filtered := filterParamsValues(queries)

	index = 0
	buffer.WriteString("{ ")

	for key := range queries_filtered {
		buffer.WriteString("\"")
		buffer.WriteString(key)
		buffer.WriteString("\": \"")

		buffer.WriteString(strings.Join(queries_filtered[key], " "))
		buffer.WriteString("\"")

		if (index + 1) != len(queries_filtered) {
			buffer.WriteString(", ")
		}

		index++
	}

	buffer.WriteString(" }")

	return buffer.String()
}
