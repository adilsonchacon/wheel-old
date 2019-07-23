package newcrud

import (
	"github.com/adilsonchacon/wheel/commons/notify"
	"github.com/adilsonchacon/wheel/generator/gencommon"
	"github.com/adilsonchacon/wheel/generator/newauthorize"
	"github.com/adilsonchacon/wheel/generator/newmigrate"
	"github.com/adilsonchacon/wheel/generator/newroutes"
	"github.com/adilsonchacon/wheel/templates/templatecrud"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"path/filepath"
	"strings"
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
		Name:                strcase.ToCamel(columnName),
		NameSnakeCase:       strcase.ToSnake(columnName),
		NameSnakeCasePlural: inflection.Plural(strcase.ToSnake(columnName)),
		Type:                columnType,
		Extras:              extra,
		IsReference:         isReference,
		IsForeignKey:        isForeignKey,
	}
}

func setMigrate() {
	newCode := gencommon.GenerateMigrateNewCode(templatecrud.MigrateContent, templateVar)
	currentFullCode := gencommon.ReadTextFile(filepath.Join(".", "db", "schema"), "migrate.go")
	newFullCode, err := newmigrate.AppendNewCode(newCode, currentFullCode)

	if err != nil {
		notify.WarnAppendToMigrate(err, newCode)
	} else if newFullCode == "" {
		notify.Identical(filepath.Join(".", "db", "schema", "migrate.go"))
	} else {
		gencommon.UpdateTextFile(newFullCode, filepath.Join(".", "db", "schema"), "migrate.go")
	}
}

func setRoutes() {
	newCode := gencommon.GenerateRoutesNewCode(templatecrud.RoutesContent, templateVar)
	currentFullCode := gencommon.ReadTextFile(filepath.Join(".", "routes"), "routes.go")
	newFullCode, err := newroutes.AppendNewCode(newCode, currentFullCode)

	if err != nil {
		notify.WarnAppendToRoutes(err, newCode)
	} else if newFullCode == "" {
		notify.Identical(filepath.Join(".", "routes", "routes.go"))
	} else {
		gencommon.UpdateTextFile(newFullCode, filepath.Join(".", "routes"), "routes.go")
	}
}

func setAuthorize() {
	newCode := gencommon.GenerateAuthorizeNewCode(templatecrud.AuthorizeContent, templateVar)
	currentFullCode := gencommon.ReadTextFile(filepath.Join(".", "routes"), "authorize.go")
	newFullCode, err := newauthorize.AppendNewCode(newCode, currentFullCode)

	if err != nil {
		notify.WarnAppendToAuthorize(err, newCode)
	} else if newFullCode == "" {
		notify.Identical(filepath.Join(".", "routes", "authorize.go"))
	} else {
		gencommon.UpdateTextFile(newFullCode, filepath.Join(".", "routes"), "authorize.go")
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
		SnakeCasePlural:      inflection.Plural(strcase.ToSnake(entityName)),
		LowerCase:            strings.ToLower(strcase.ToCamel(entityName)),
	}

	templateVar = gencommon.TemplateVar{AppRepository: gencommon.GetAppConfig().AppRepository, EntityName: tEntityName, EntityColumns: entityColumns}

	if options["model"] {
		path = []string{".", "app", tEntityName.LowerCase, tEntityName.SnakeCase + "_model.go"}
		gencommon.GeneratePathAndFileFromTemplateString(path, templatecrud.ModelContent, templateVar)
	}

	if options["view"] {
		path = []string{".", "app", tEntityName.LowerCase, tEntityName.SnakeCase + "_view.go"}
		gencommon.GeneratePathAndFileFromTemplateString(path, templatecrud.ViewContent, templateVar)
	}

	if options["entity"] {
		path = []string{".", "db", "entities", tEntityName.SnakeCase + "_entity.go"}
		gencommon.GeneratePathAndFileFromTemplateString(path, templatecrud.EntityContent, templateVar)
	}

	if options["handler"] {
		path = []string{".", "app", "handlers", tEntityName.SnakeCase + "_handler.go"}
		gencommon.GeneratePathAndFileFromTemplateString(path, templatecrud.HandlerContent, templateVar)
	}

	if options["routes"] {
		setRoutes()
	}

	if options["migrate"] {
		setMigrate()
	}

	if options["authorize"] {
		setAuthorize()
	}

}
