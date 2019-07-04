package newcrud

import (
	// "fmt"
	// "log"
	// "path/filepath"
	"strings"
	"wheel.smart26.com/commons/inflection"
	"wheel.smart26.com/commons/strcase"
	"wheel.smart26.com/generator/gencommon"
	"wheel.smart26.com/templates/templatecrud"
)

var entityColumns []gencommon.EntityColumn
var templateVar gencommon.TemplateVar

func optionToEntityColumn(options string, isForeignKey bool) gencommon.EntityColumn {
	var columnName, columnType, extra string
	var isReference bool

	columnData := strings.Split(options, ":")
	if len(columnData) == 1 {
		columnData = append(columnData, "string")
		columnData = append(columnData, "")
	} else if len(columnData) == 2 {
		columnData = append(columnData, "")
	}

	if isForeignKey {
		columnName = columnData[0]
		columnType = strcase.ToCamel(columnData[0])
		extra = ""
		isReference = false
	} else {
		columnName, columnType, extra, isReference = gencommon.GetColumnInfo(columnData[0], columnData[1], columnData[2])
	}

	return gencommon.EntityColumn{
		Name:          strcase.ToCamel(columnName),
		NameSnakeCase: strcase.ToSnake(columnName),
		Type:          columnType,
		Extras:        extra,
		IsReference:   isReference,
		IsForeignKey:  isForeignKey,
	}
}

func Generate(entityName string, columns []string, options map[string]bool) {
	var path []string

	for _, column := range columns {
		entityColumns = append(entityColumns, optionToEntityColumn(column, false))

		if entityColumns[len(entityColumns)-1].IsReference {
			entityColumns = append(entityColumns, optionToEntityColumn(column, true))
		}
	}

	tEntityName := gencommon.EntityName{
		CamelCase:            strcase.ToCamel(entityName),
		CamelCasePlural:      inflection.Plural(strcase.ToCamel(entityName)),
		LowerCamelCase:       strcase.ToLowerCamel(entityName),
		LowerCamelCasePlural: inflection.Plural(strcase.ToLowerCamel(entityName)),
		SnakeCase:            strcase.ToSnake(entityName),
		LowerCase:            strings.ToLower(strcase.ToCamel(entityName)),
	}

	templateVar = gencommon.TemplateVar{AppDomain: "test", EntityName: tEntityName, EntityColumns: entityColumns}

	if options["model"] {
		path = []string{"app", tEntityName.LowerCase, tEntityName.SnakeCase + "_model.go"}
		gencommon.GenerateFromTemplateFullPath(".", path, templatecrud.ModelContent, templateVar, false)
	}

	if options["view"] {
		path = []string{"app", tEntityName.LowerCase, tEntityName.SnakeCase + "_view.go"}
		gencommon.GenerateFromTemplateFullPath(".", path, templatecrud.ViewContent, templateVar, false)
	}

	if options["entity"] {
		path = []string{"db", "entities", tEntityName.SnakeCase + "_entity.go"}
		gencommon.GenerateFromTemplateFullPath(".", path, templatecrud.EntityContent, templateVar, false)
	}

	if options["handler"] {
		path = []string{"app", "handlers", tEntityName.SnakeCase + "_handler.go"}
		gencommon.GenerateFromTemplateFullPath(".", path, templatecrud.HandlerContent, templateVar, false)
	}
}
