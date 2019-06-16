package generator

import (
	"wheel.smart26.com/generator/newapp"
	"wheel.smart26.com/generator/single"
)

func NewApp(options map[string]string) {
	newapp.Generate(options)
}

func Single(entityName string, options []string) {
	single.Generate(entityName, options)
}
