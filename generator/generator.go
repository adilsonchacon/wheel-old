package generator

import (
	"github.com/adilsonchacon/wheel/generator/newapp"
	"github.com/adilsonchacon/wheel/generator/newcrud"
)

func NewApp(options map[string]interface{}) {
	newapp.Generate(options)
}

func NewCrud(entityName string, columns []string, options map[string]bool) {
	newcrud.Generate(entityName, columns, options)
}
