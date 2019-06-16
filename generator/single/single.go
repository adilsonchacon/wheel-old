package single

import (
	"bytes"
	"fmt"
	"log"
	"os"
	// "path/filepath"
	"strings"
	"text/template"
	"wheel.smart26.com/templates/singles"
)

type EntityColumn struct {
	Name   string
	Type   string
	Extras string
}

type TemplateVar struct {
	EntityName    string
	EntityColumns []EntityColumn
}

var entityColumns []EntityColumn
var templateVars TemplateVar

func saveTextFile(content string, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}

	f.Sync()

}

func GenerateFromTemplateString(content string, destinyPath string) {
	var buffContent bytes.Buffer

	fmt.Println("created:", destinyPath)

	tmpl, err := template.New("T").Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&buffContent, templateVars)
	if err != nil {
		log.Fatal(err)
	}

	saveTextFile(buffContent.String(), destinyPath)
}

func Generate(entityName string, options []string) {
	var columns []string

	for _, value := range options {
		columns = strings.Split(value, ":")
		entityColumns = append(entityColumns, EntityColumn{Name: columns[0], Type: columns[1], Extras: columns[2]})
	}

	templateVars = TemplateVar{EntityName: entityName, EntityColumns: entityColumns}

	GenerateFromTemplateString(singles.EntityContent, "test")
}
