package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/app/views"
	"wheel.smart26.com/utils"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var user = models.User{}

	userSetParams(&user, r)

	utils.LoggerInfo().Println("controllers: Create")
	w.Header().Set("Content-Type", "application/json")

	if models.UserCreate(&user) {
		json.NewEncoder(w).Encode(views.UserSuccessfullySavedJson{SystemMessage: views.SetSystemMessage("notice", "user was successfully created"), User: views.SetUserJson(user)})
	} else {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "user was not created", models.Errors))
	}
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.UserFind(params["id"])

	userSetParams(&user, r)

	utils.LoggerInfo().Println("controllers: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	if models.UserUpdate(&user) {
		json.NewEncoder(w).Encode(views.UserSuccessfullySavedJson{SystemMessage: views.SetSystemMessage("notice", "user was successfully updated"), User: views.SetUserJson(user)})
	} else {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "user was not updated", models.Errors))
	}
}

func UserDestroy(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.UserFind(params["id"])

	utils.LoggerInfo().Println("Controller: UserDestroy")

	w.Header().Set("Content-Type", "application/json")

	if models.UserDestroy(&user) {
		json.NewEncoder(w).Encode(views.SetSystemMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(views.SetSystemMessage("alert", "user was not found"))
	}
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.UserFind(params["id"])

	utils.LoggerInfo().Println("controllers: UserShow")

	w.Header().Set("Content-Type", "application/json")

	if user.ID != 0 {
		json.NewEncoder(w).Encode(views.SetUserJson(user))
	} else {
		json.NewEncoder(w).Encode(views.SetSystemMessage("alert", "user was not found"))
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var userJsons []views.UserJson
	var users []models.User

	m := normalizeUrlQueryParams("search", r.URL.Query())

	utils.LoggerInfo().Println("controllers: UserList")

	users, page, pages, entries = models.UserPaginate(m, r.FormValue("page"), r.FormValue("perPage"))

	for i = 0; i < len(users); i++ {
		userJsons = append(userJsons, views.SetUserJson(users[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	pagination := views.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode(views.UserPaginationJson{Pagination: pagination, Users: userJsons})
}

func userSetParams(user *models.User, r *http.Request) {
	var allowedParams = []string{"name", "email", "password", "admin", "locale"}

	r.ParseForm()

	for key := range r.Form {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				models.SetColumnValue(user, allowedParam, r.FormValue(allowedParam))
				break
			}
		}
	}
}
