package single

import (
	"bytes"
	"fmt"
	"log"
	"os"
	// "path/filepath"
	"regexp"
	"strings"
	"text/template"
	"wheel.smart26.com/inflection"
	"wheel.smart26.com/strcase"
	"wheel.smart26.com/templates/singles"
)

type EntityColumn struct {
	Name          string
	NameSnakeCase string
	Type          string
	Extras        string
}

type EntityNames struct {
	Name           string
	Plural         string
	LowerCamelCase string
	SnakeCase      string
}

type TemplateVar struct {
	AppDomain     string
	EntityNames   EntityNames
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

func gormSpecificationForString(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"type:varchar(255)" + index + "\"`"
}

func gormSpecificationForIntegers(extra string) string {
	var index string

	if extra == "index" {
		index = "`gorm:\"index\"`"
	} else {
		index = ""
	}

	return index
}

func gormSpecificationForDecimals(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"type:decimal\"" + index + "`"
}

func gormSpecificationForDatetime(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"default:null\"" + index + "`"
}

func gormSpecificationForBoolean(extra string) string {
	var index string

	if extra == "index" {
		index = "`gorm:\"default:null\";index`"
	} else if extra == "true" || extra == "t" {
		index = "`gorm:\"default:true\"`"
	} else if extra == "false" || extra == "f" {
		index = "`gorm:\"default:false\"`"
	} else {
		index = ""
	}

	return index
}

func getColumnType(columnType string, extra string) (string, string) {
	var regexpText = regexp.MustCompile(`text`)
	var regexpString = regexp.MustCompile(`string`)
	var regexpDecimal = regexp.MustCompile(`float|double|decimal`)
	var regexpInteger = regexp.MustCompile(`int|integer`)
	var regexpUnsignedInteger = regexp.MustCompile(`uint`)
	var regexpDatetime = regexp.MustCompile(`datetime`)
	var regexpBoolean = regexp.MustCompile(`bool`)

	// TODO: datetime
	if regexpText.MatchString(columnType) {
		columnType = "string"
		extra = "`gorm:\"type:text\"`"
	} else if regexpString.MatchString(columnType) || regexpText.MatchString(columnType) {
		columnType = "string"
		extra = gormSpecificationForString(extra)
	} else if regexpUnsignedInteger.MatchString(columnType) {
		columnType = "uint"
		extra = gormSpecificationForIntegers(extra)
	} else if regexpInteger.MatchString(columnType) {
		columnType = "int64"
		extra = gormSpecificationForIntegers(extra)
	} else if regexpDatetime.MatchString(columnType) {
		columnType = "*time.Time"
		extra = gormSpecificationForDatetime(extra)
	} else if regexpBoolean.MatchString(columnType) {
		columnType = "bool"
		extra = gormSpecificationForBoolean(extra)
	} else if regexpDecimal.MatchString(columnType) {
		columnType = "float64"
		extra = gormSpecificationForDecimals(extra)
	}

	return columnType, extra
}

func optionToEntityColumn(options string) EntityColumn {
	columnData := strings.Split(options, ":")
	if len(columnData) <= 2 {
		columnData = append(columnData, "")
	}
	columnType, extra := getColumnType(columnData[1], columnData[2])

	return EntityColumn{Name: strcase.ToCamel(columnData[0]), NameSnakeCase: strcase.ToSnake(columnData[0]), Type: columnType, Extras: extra}
}

func Generate(entityName string, options []string) {
	for _, option := range options {
		entityColumns = append(entityColumns, optionToEntityColumn(option))
	}

	entityNames := EntityNames{Name: strcase.ToCamel(entityName), Plural: inflection.Plural(strcase.ToCamel(entityName)), LowerCamelCase: strcase.ToLowerCamel(entityName), SnakeCase: strcase.ToSnake(entityName)}

	templateVars = TemplateVar{AppDomain: "test", EntityNames: entityNames, EntityColumns: entityColumns}

	GenerateFromTemplateString(singles.EntityContent, strcase.ToSnake(entityName)+"_entity.go")
	GenerateFromTemplateString(singles.ViewContent, strcase.ToSnake(entityName)+"_view.go")
}
