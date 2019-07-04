package generator

import (
	"wheel.smart26.com/generator/newapp"
	"wheel.smart26.com/generator/newcrud"
)

func NewApp(options map[string]string) {
	newapp.Generate(options)
}

func NewCrud(entityName string, columns []string, options map[string]bool) {
	newcrud.Generate(entityName, columns, options)
}
