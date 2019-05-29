package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"wheel.smart26.com/app/entity"
	"wheel.smart26.com/app/handler"
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/app/view"
	"wheel.smart26.com/commons/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info.Println(r.Method + ": " + filterUrlValues(r.URL.Path, r.URL.Query()) + " for " + r.RemoteAddr)
		r.ParseMultipartForm(100 * 1024)
		log.Info.Println("Params: " + filterFormValues(r.Form))
		next.ServeHTTP(w, r)
	})
}

func authorizeMiddleware(next http.Handler) http.Handler {
	var userId uint
	var err error
	var userRole string
	var signedInUser entity.User

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId = 0
		err = nil
		userRole = "public"

		userId, err = checkToken(r.Header.Get("token"))
		if err == nil {
			if signedInUser, err = checkSignedInUser(userId); err == nil {
				userRole = "signed_in"
			}

			if userRole == "signed_in" && checkAdminUser(user.Current) {
				userRole = "admin"
			}
		}

		if GrantPermission(r.RequestURI, r.Method, userRole) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
		}
	})
}

func checkAdminUser(signedInUser entity.User) bool {
	if signedInUser.Admin {
		log.Info.Println("admin access granted to user")
		return true
	} else {
		log.Info.Println("admin access denied to user")
		return false
	}
}

func checkSignedInUser(userId uint) (entity.User, error) {
	log.Info.Printf("checking user id: %d...\n", userId)

	if signedInUser, err := user.Find(userId); err == nil {
		user.SetCurrent(userId)
		log.Info.Println("user was found")
		return signedInUser, nil
	} else {
		log.Info.Println("user was not found")
		return signedInUser, errors.New("user was not found")
	}
}

func checkToken(token string) (uint, error) {
	log.Info.Println("checking token...")

	userId, err := handler.SessionCheck(token)

	if err == nil {
		log.Info.Println("token is valid")
		return userId, nil
	} else {
		log.Info.Println("invalid token")
		return 0, errors.New("invalid token")
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
