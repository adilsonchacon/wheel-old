package utils

import (
	"io"
	"log"
	"os"
)

var writer io.Writer

func LoggerInfo() *log.Logger {
	return log.New(writer, "[INFO] ", log.LstdFlags)
}

func LoggerWarn() *log.Logger {
	return log.New(writer, "[WARN] ", log.LstdFlags)
}

func LoggerError() *log.Logger {
	return log.New(writer, "[ERROR] ", log.LstdFlags)
}

func LoggerFatal() *log.Logger {
	return log.New(writer, "[FATAL] ", log.LstdFlags)
}

func init() {
	os.MkdirAll("./log", 0755)

	fileHandler, err := os.OpenFile("./log/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	writer = io.MultiWriter(os.Stdout, fileHandler)
}
