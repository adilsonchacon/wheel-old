package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"wheel.smart26.com/app/controllers"
	"wheel.smart26.com/utils"
)

func Routes(host string, port string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(controllers.Error404)

	// middlewares
	utils.LoggerInfo().Println("setting up middlewares")
	router.Use(loggingMiddleware)
	router.Use(authorizeMiddleware)

	utils.LoggerInfo().Println("setting up routes")
	// sessions
	router.HandleFunc("/sessions/sign_in", controllers.SessionSignIn).Methods("POST")
	router.HandleFunc("/sessions/sign_out", controllers.SessionSignOut).Methods("DELETE")
	router.HandleFunc("/sessions/sign_up", controllers.SessionSignUp).Methods("POST")
	router.HandleFunc("/sessions/password", controllers.SessionPassword).Methods("POST")
	router.HandleFunc("/sessions/password", controllers.SessionRecovery).Methods("PUT")
	router.HandleFunc("/sessions/refresh", controllers.SessionRefresh).Methods("POST")

	// user
	router.HandleFunc("/myself", controllers.MyselfShow).Methods("GET")
	router.HandleFunc("/myself", controllers.MyselfUpdate).Methods("PUT")
	router.HandleFunc("/myself/password", controllers.MyselfUpdatePassword).Methods("PUT")
	router.HandleFunc("/myself/destroy", controllers.MyselfDestroy).Methods("DELETE")

	// admin
	router.HandleFunc("/users", controllers.UserList).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.UserShow).Methods("GET")
	router.HandleFunc("/users", controllers.UserCreate).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.UserUpdate).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.UserDestroy).Methods("DELETE")

	utils.LoggerInfo().Println("listening on " + host + ":" + port + ", CTRL+C to stop")

	return router
}
