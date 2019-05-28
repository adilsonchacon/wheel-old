package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"wheel.smart26.com/app/handler"
	"wheel.smart26.com/commons/controller"
	"wheel.smart26.com/commons/log"
)

func Routes(host string, port string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(controller.Error404)

	// middlewares
	log.Info.Println("setting up middlewares")
	router.Use(loggingMiddleware)
	router.Use(authorizeMiddleware)

	log.Info.Println("setting up routes")
	// sessions
	router.HandleFunc("/sessions/sign_in", handler.SessionSignIn).Methods("POST")
	router.HandleFunc("/sessions/sign_out", handler.SessionSignOut).Methods("DELETE")
	router.HandleFunc("/sessions/sign_up", handler.SessionSignUp).Methods("POST")
	router.HandleFunc("/sessions/password", handler.SessionPassword).Methods("POST")
	router.HandleFunc("/sessions/password", handler.SessionRecovery).Methods("PUT")
	router.HandleFunc("/sessions/refresh", handler.SessionRefresh).Methods("POST")

	// user
	router.HandleFunc("/myself", handler.MyselfShow).Methods("GET")
	router.HandleFunc("/myself", handler.MyselfUpdate).Methods("PUT")
	router.HandleFunc("/myself/password", handler.MyselfUpdatePassword).Methods("PUT")
	router.HandleFunc("/myself/destroy", handler.MyselfDestroy).Methods("DELETE")

	// admin
	router.HandleFunc("/users", handler.UserList).Methods("GET")
	router.HandleFunc("/users/{id}", handler.UserShow).Methods("GET")
	router.HandleFunc("/users", handler.UserCreate).Methods("POST")
	router.HandleFunc("/users/{id}", handler.UserUpdate).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.UserDestroy).Methods("DELETE")

	log.Info.Println("listening on " + host + ":" + port + ", CTRL+C to stop")

	return router
}
