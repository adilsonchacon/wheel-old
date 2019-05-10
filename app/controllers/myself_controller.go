package controllers

import (
	"encoding/json"
	"net/http"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/app/views"
	"wheel.smart26.com/utils"
)

func MyselfUpdate(w http.ResponseWriter, r *http.Request) {
	user := models.UserCurrent

	myselfSetParams(&user, r)

	utils.LoggerInfo().Println("controllers: MyselfUpdate")
	w.Header().Set("Content-Type", "application/json")

	if user.ID != 0 && models.UserSave(&user) {
		json.NewEncoder(w).Encode(views.SetSystemMessage("notice", "user was successfully updated"))
	} else {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "user was not updated", models.Errors))
	}
}

func MyselfUpdatePassword(w http.ResponseWriter, r *http.Request) {
	var errors []string
	user := models.UserCurrent

	utils.LoggerInfo().Println("controllers: MyselfChangePassword")
	w.Header().Set("Content-Type", "application/json")

	if !models.UserExists(user) {
		errors = append(errors, "invalid user")
	} else if r.FormValue("new_password") != r.FormValue("password_confirmation") {
		errors = append(errors, "password confirmation does not match new password")
	} else if !utils.SaferCheckPassword(r.FormValue("password"), user.Password) {
		errors = append(errors, "invalid password")
	} else if user.Password = r.FormValue("new_password"); models.UserSave(&user) {
		json.NewEncoder(w).Encode(views.SetSystemMessage("notice", "password was successfully changed"))
	} else {
		errors = models.Errors
	}

	if len(errors) > 0 {
		json.NewEncoder(w).Encode(views.SetErrorMessage("alert", "password could not be changed", errors))
	}
}

func MyselfDestroy(w http.ResponseWriter, r *http.Request) {
	user := models.UserCurrent

	utils.LoggerInfo().Println("Controller: MyselfDestroy")

	w.Header().Set("Content-Type", "application/json")

	if models.UserDestroy(&user) {
		json.NewEncoder(w).Encode(views.SetSystemMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(views.SetSystemMessage("alert", "user could not be destroyed"))
	}
}

func MyselfShow(w http.ResponseWriter, r *http.Request) {
	utils.LoggerInfo().Println("controllers: MyselfShow")

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(views.SetMyselfJson(models.UserCurrent))
}

func myselfSetParams(user *models.User, r *http.Request) {
	var allowedParams = []string{"name", "locale"}

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
