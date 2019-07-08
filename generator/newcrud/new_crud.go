package newcrud

import (
	"fmt"
	"path/filepath"
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
	var newCode, currentFullCode, newFullCode string
	var err error

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

	templateVar = gencommon.TemplateVar{AppDomain: "test", EntityName: tEntityName, EntityColumns: entityColumns}

	if options["model"] {
		path = []string{".", "app", tEntityName.SnakeCase + "_model.go"}
		gencommon.GeneratePathAndFileFromTemplateString(path, templatecrud.ModelContent, templateVar)
	}

	if options["view"] {
		path = []string{".", "app", tEntityName.SnakeCase + "_view.go"}
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
		newCode = gencommon.GenerateRoutesNewCode(templatecrud.RoutesContent, templateVar)
		currentFullCode = gencommon.ReadTextFile(filepath.Join(".", "routes"), "routes.go")
		newFullCode, err = AppendNewCodeToRoutes(newCode, currentFullCode)
		if err == nil {
			gencommon.SaveTextFile(newFullCode, filepath.Join(".", "routes"), "routes.go")
		} else {
			fmt.Println("--------------------------------------")
			fmt.Println("")
			fmt.Printf("\033[93mwarn:\033[39m ")
			fmt.Println(err)
			fmt.Println("Edit file \"routes/rotes.go\" and add this new line bellow at \"func Routes()\"")
			fmt.Println("")
			fmt.Println(newCode)
			fmt.Println("")
			fmt.Println("--------------------------------------\n")
		}
	}

	if options["migrate"] {
		newCode = gencommon.GenerateMigrateNewCode(templatecrud.MigrateContent, templateVar)
		currentFullCode = gencommon.ReadTextFile(filepath.Join(".", "db", "schema"), "migrate.go")
		newFullCode, err = AppendNewCodeToMigrate(newCode, currentFullCode)
		if err == nil {
			gencommon.SaveTextFile(newFullCode, filepath.Join(".", "db", "schema"), "migrate.go")
		} else {
			fmt.Println("--------------------------------------")
			fmt.Println("")
			fmt.Printf("\033[93mwarn:\033[39m ")
			fmt.Println(err)
			fmt.Println("Edit file \"db/schema/migrate.go\" and add this new line bellow at \"func Migrate()\"")
			fmt.Println("")
			fmt.Println(newCode)
			fmt.Println("")
			fmt.Println("--------------------------------------\n")
		}
	}
}
