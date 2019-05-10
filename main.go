package main

import (
	"flag"
	"net/http"
	"wheel.smart26.com/app/models"
	"wheel.smart26.com/config"
	"wheel.smart26.com/routes"
	"wheel.smart26.com/utils"
)

func main() {
	var mode string
	var port string
	var host string

	flag.StringVar(&mode, "mode", "server", "run mode (options: server/migrate)")
	flag.StringVar(&host, "host", "localhost", "http server host")
	flag.StringVar(&port, "port", "8081", "http server port")
	flag.Parse()

	utils.LoggerInfo().Println("starting app", config.AppName())

	models.DbConnect()

	if mode == "migrate" {
		models.Migrate()
	} else if mode == "s" || mode == "server" {
		utils.LoggerFatal().Fatal(http.ListenAndServe(host+":"+port, routes.Routes(host, port)))
	} else {
		utils.LoggerFatal().Fatal("invalid run mode, please, use \"--help\" for more details")
	}
}
