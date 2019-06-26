package single

import (
	// "fmt"
	// "log"
	// "path/filepath"
	"strings"
	"wheel.smart26.com/generator/common"
	"wheel.smart26.com/inflection"
	"wheel.smart26.com/strcase"
	"wheel.smart26.com/templates/singles"
)

var entityColumns []common.EntityColumn
var templateVar common.TemplateVar

func optionToEntityColumn(options string, isForeignKey bool) common.EntityColumn {
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
		columnName, columnType, extra, isReference = common.GetColumnInfo(columnData[0], columnData[1], columnData[2])
	}

	return common.EntityColumn{
		Name:          strcase.ToCamel(columnName),
		NameSnakeCase: strcase.ToSnake(columnName),
		Type:          columnType,
		Extras:        extra,
		IsReference:   isReference,
		IsForeignKey:  isForeignKey,
	}
}

func Generate(entityName string, options []string) {
	for _, option := range options {
		entityColumns = append(entityColumns, optionToEntityColumn(option, false))

		if entityColumns[len(entityColumns)-1].IsReference {
			entityColumns = append(entityColumns, optionToEntityColumn(option, true))
		}
	}

	tEntityName := common.EntityName{
		CamelCase:            strcase.ToCamel(entityName),
		CamelCasePlural:      inflection.Plural(strcase.ToCamel(entityName)),
		LowerCamelCase:       strcase.ToLowerCamel(entityName),
		LowerCamelCasePlural: inflection.Plural(strcase.ToLowerCamel(entityName)),
		SnakeCase:            strcase.ToSnake(entityName),
		LowerCase:            strings.ToLower(strcase.ToCamel(entityName)),
	}

	templateVar = common.TemplateVar{AppDomain: "test", EntityName: tEntityName, EntityColumns: entityColumns}

	common.GenerateFromTemplateString(singles.EntityContent, strcase.ToSnake(entityName)+"_entity.go", templateVar)
	common.GenerateFromTemplateString(singles.ViewContent, strcase.ToSnake(entityName)+"_view.go", templateVar)
	common.GenerateFromTemplateString(singles.ModelContent, strcase.ToSnake(entityName)+"_model.go", templateVar)
	common.GenerateFromTemplateString(singles.HandlerContent, strcase.ToSnake(entityName)+"_handler.go", templateVar)
}
