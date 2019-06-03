package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"wheel.smart26.com/app/myself"
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/app/model"
	"wheel.smart26.com/commons/app/view"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/db/entities"
)

func MyselfUpdate(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfUpdate")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	myselfSetParams(&userMyself, r)

	if valid, errs := user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully updated"))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func MyselfUpdatePassword(w http.ResponseWriter, r *http.Request) {
	var errs []error
	var valid bool

	log.Info.Println("Handler: MyselfChangePassword")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current
	userMyself.Password = r.FormValue("new_password")

	if !user.Exists(&userMyself) {
		errs = append(errs, errors.New("invalid user"))
	} else if r.FormValue("new_password") != r.FormValue("password_confirmation") {
		errs = append(errs, errors.New("password confirmation does not match new password"))
	} else if valid, errs = user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "password was successfully changed"))
	}

	if len(errs) > 0 {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "password could not be changed", errs))
	}
}

func MyselfDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfDestroy")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	if user.Destroy(&userMyself) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "user could not be destroyed"))
	}
}

func MyselfShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfShow")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(myself.SetJson(user.Current))
}

func myselfSetParams(userMyself *entities.User, r *http.Request) {
	var allowedParams = []string{"name", "locale"}

	r.ParseForm()

	for key := range r.Form {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				model.SetColumnValue(userMyself, allowedParam, r.FormValue(allowedParam))
				break
			}
		}
	}
}
