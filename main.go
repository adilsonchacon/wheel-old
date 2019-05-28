package main

import (
	"flag"
	"net/http"
	// "wheel.smart26.com/app/models"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/config"
	"wheel.smart26.com/routes"
	// "wheel.smart26.com/app/session"
	// "wheel.smart26.com/app/user"
	"wheel.smart26.com/app/migration"
	"wheel.smart26.com/commons/db"
)

func main() {
	var mode string
	var port string
	var host string

	flag.StringVar(&mode, "mode", "server", "run mode (options: server/migrate)")
	flag.StringVar(&host, "host", "localhost", "http server host")
	flag.StringVar(&port, "port", "8081", "http server port")
	flag.Parse()

	log.Info.Println("starting app", config.AppName())

	db.Connect()

	// log.Info.Println(user.Find(1))
	// log.Info.Println(session.Find(1))
	//
	if mode == "migrate" {
		migration.Run()
	} else if mode == "s" || mode == "server" {
		log.Fatal.Println(http.ListenAndServe(host+":"+port, routes.Routes(host, port)))
	} else {
		log.Fatal.Println("invalid run mode, please, use \"--help\" for more details")
	}
}
