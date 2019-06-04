package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"wheel.smart26.com/app/handlers"
	"wheel.smart26.com/commons/app/handler"
	"wheel.smart26.com/commons/log"
)

func Routes(host string, port string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(handler.Error404)

	// middlewares
	log.Info.Println("setting up middlewares")
	router.Use(loggingMiddleware)
	router.Use(authorizeMiddleware)

	log.Info.Println("setting up routes")
	// sessions
	router.HandleFunc("/sessions/sign_in", handlers.SessionSignIn).Methods("POST")
	router.HandleFunc("/sessions/sign_out", handlers.SessionSignOut).Methods("DELETE")
	router.HandleFunc("/sessions/sign_up", handlers.SessionSignUp).Methods("POST")
	router.HandleFunc("/sessions/password", handlers.SessionPassword).Methods("POST")
	router.HandleFunc("/sessions/password", handlers.SessionRecovery).Methods("PUT")
	router.HandleFunc("/sessions/refresh", handlers.SessionRefresh).Methods("POST")

	// user
	router.HandleFunc("/myself", handlers.MyselfShow).Methods("GET")
	router.HandleFunc("/myself", handlers.MyselfUpdate).Methods("PUT")
	router.HandleFunc("/myself/password", handlers.MyselfUpdatePassword).Methods("PUT")
	router.HandleFunc("/myself/destroy", handlers.MyselfDestroy).Methods("DELETE")

	// admin
	router.HandleFunc("/users", handlers.UserList).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.UserShow).Methods("GET")
	router.HandleFunc("/users", handlers.UserCreate).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.UserUpdate).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.UserDestroy).Methods("DELETE")

	log.Info.Println("listening on " + host + ":" + port + ", CTRL+C to stop")

	return router
}
