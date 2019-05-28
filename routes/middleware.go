package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"wheel.smart26.com/app/handler"
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/app/view"
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
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
			}
		} else if userUrl.MatchString(r.RequestURI) || sessionRefreshUrl.MatchString(r.RequestURI) {
			if authorizeUser(token) {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func authorizeAdmin(token string) bool {
	var id uint

	id, ok := checkToken(token)
	if !ok {
		return false
	}

	log.Info.Printf("checking admin for user id: %d...\n", id)

	authUser, err := user.Find(id)
	if err == nil && authUser.Admin {
		log.Info.Println("access granted to user")
		return true
	} else {
		log.Info.Println("access denied to user")
		return false
	}
}

func authorizeUser(token string) bool {
	var id uint

	id, ok := checkToken(token)
	if !ok {
		return false
	}

	log.Info.Printf("checking user id: %d...\n", id)

	_, err := user.Find(id)
	if err == nil {
		user.SetCurrent(id)
		log.Info.Println("access granted to user")
		return true
	} else {
		log.Info.Println("access denied to user")
		return false
	}

}

func checkToken(token string) (uint, bool) {
	log.Info.Println("checking token...")

	id, err := handler.SessionCheck(token)
	if err != nil {
		log.Info.Println("invalid token")
		return 0, false
	} else {
		log.Info.Println("token is valid")
		return id, true
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
