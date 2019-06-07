package handlers

var UserPath = []string{"app", "handlers", "user_handler.go"}

var UserContent = `package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"{{ .AppDomain }}/app/user"
	"{{ .AppDomain }}/commons/app/apphandler"
	"{{ .AppDomain }}/commons/app/model"
	"{{ .AppDomain }}/commons/app/view"
	"{{ .AppDomain }}/commons/log"
	"{{ .AppDomain }}/db/entities"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var newUser = entities.User{}

	log.Info.Println("Handler: UserCreate")
	w.Header().Set("Content-Type", "application/json")

	userSetParams(&newUser, r)

	valid, errs := user.Create(&newUser)
	log.Debug.Println(errs)

	if valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully created"), User: user.SetJson(newUser)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not created", errs))
	}
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userCurrent, err := user.Find(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", []error{err}))
		return
	}

	userSetParams(&userCurrent, r)

	if valid, errs := user.Update(&userCurrent); valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully updated"), User: user.SetJson(userCurrent)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func UserDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserDestroy")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil && user.Destroy(&userCurrent) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "user was not found"))
	}
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserShow")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil {
		json.NewEncoder(w).Encode(user.SetJson(userCurrent))
	} else {
		json.NewEncoder(w).Encode(view.SetSystemMessage("alert", "user was not found"))
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var userJsons []user.Json
	var userList []entities.User

	log.Info.Println("Handler: UserList")
	w.Header().Set("Content-Type", "application/json")

	normalizedParams := handler.NormalizeUrlQueryParams("search", r.URL.Query())

	userList, page, pages, entries = user.Paginate(normalizedParams, r.FormValue("page"), r.FormValue("per_page"))

	for i = 0; i < len(userList); i++ {
		userJsons = append(userJsons, user.SetJson(userList[i]))
	}

	pagination := view.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode(user.PaginationJson{Pagination: pagination, Users: userJsons})
}

func userSetParams(userSet *entities.User, r *http.Request) {
	var allowedParams = []string{"name", "email", "password", "admin", "locale"}

	r.ParseForm()

	for key := range r.Form {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				model.SetColumnValue(userSet, allowedParam, r.FormValue(allowedParam))
				break
			}
		}
	}
}`