package gencommon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"text/template"
	"time"
)

type EntityColumn struct {
	Name          string
	NameSnakeCase string
	Type          string
	Extras        string
	IsReference   bool
	IsForeignKey  bool
}

type EntityName struct {
	CamelCase            string
	CamelCasePlural      string
	LowerCamelCase       string
	LowerCamelCasePlural string
	SnakeCase            string
	SnakeCasePlural      string
	LowerCase            string
}

type TemplateVar struct {
	AppDomain     string
	AppName       string
	SecretKey     string
	EntityName    EntityName
	EntityColumns []EntityColumn
}

func SaveTextFile(content string, filePath string, fileName string) {
	if err := os.MkdirAll(filePath, 0775); err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(filePath, fileName)

	f, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}

	fmt.Println("created:", fullPath)

	f.Sync()
}

func ReadTextFile(filePath string, fileName string) string {
	file, err := os.Open(filepath.Join(filePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)

	return string(b)
}

func GenerateFromTemplateString(content string, templateVar TemplateVar) string {
	var buffContent bytes.Buffer

	FuncMap := template.FuncMap{
		"isLastIndex": func(index int, tSlice interface{}) bool {
			return index == reflect.ValueOf(tSlice).Len()-1
		},
		"isNotLastIndex": func(index int, tSlice interface{}) bool {
			return index != reflect.ValueOf(tSlice).Len()-1
		},
		"filterEntityColumnsNotForeignKeys": func(tEntityColumns []EntityColumn) []EntityColumn {
			var notForeignKeys []EntityColumn
			for _, element := range tEntityColumns {
				if !element.IsForeignKey {
					notForeignKeys = append(notForeignKeys, element)
				}
			}
			return notForeignKeys
		},
	}

	tmpl, err := template.New("T").Funcs(FuncMap).Parse(content)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&buffContent, templateVar)
	if err != nil {
		log.Fatal(err)
	}

	return buffContent.String()
}

func GenerateFromTemplateFile(templatePath string, templateVar TemplateVar) string {
	var content bytes.Buffer

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(&content, &templateVar)
	if err != nil {
		log.Fatal(err)
	}

	return content.String()
}

func GeneratePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) {
	fileName, filePath := path[len(path)-1], path[:len(path)-1]
	content = GenerateFromTemplateString(content, templateVar)
	SaveTextFile(content, sliceToPath(filePath), fileName)
}

func CreatePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) {
	fileName, filePath := path[len(path)-1], path[:len(path)-1]
	SaveTextFile(content, sliceToPath(filePath), fileName)
}

func GenerateRoutesNewCode(content string, templateVar TemplateVar) string {
	return GenerateFromTemplateString(content, templateVar)
}

func GenerateMigrateNewCode(content string, templateVar TemplateVar) string {
	return GenerateFromTemplateString(content, templateVar)
}

func HandlePathInfo(path []string) (string, string) {
	var basePath, fileName string

	for index, value := range path {
		if index+1 != len(path) {
			basePath = filepath.Join(basePath, value)
		} else {
			fileName = value
		}
	}

	return basePath, fileName
}

func SecureRandom(size int) string {
	var letters = []rune("0123456789abcdefABCDEF")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func BuildRootAppPath(appDomain string) string {
	_, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(usr.HomeDir, "go", "src", appDomain)
}

func CreateRootAppPath(rootAppPath string) {
	if err := os.MkdirAll(rootAppPath, 0775); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("created:", rootAppPath)
	}
}

func GetColumnInfo(columnName string, columnType string, extra string) (string, string, string, bool) {
	var regexpText = regexp.MustCompile(`text`)
	var regexpString = regexp.MustCompile(`string`)
	var regexpDecimal = regexp.MustCompile(`float|double|decimal`)
	var regexpInteger = regexp.MustCompile(`int|integer`)
	var regexpUnsignedInteger = regexp.MustCompile(`uint`)
	var regexpDatetime = regexp.MustCompile(`datetime`)
	var regexpBoolean = regexp.MustCompile(`bool`)
	var regexpReference = regexp.MustCompile(`reference`)

	isReference := false

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
	} else if regexpReference.MatchString(columnType) {
		columnType = "uint"
		extra = ""
		columnName = columnName + "_ID"
		isReference = true
	}

	return columnName, columnType, extra, isReference
}

func sliceToPath(path []string) string {
	var filePath string

	for index, element := range path {
		if index > 0 {
			filePath = filepath.Join(filePath, element)
		} else {
			filePath = element
		}
	}

	return filePath
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

func gormSpecificationForReference() string {
	return "`gorm:\"index\"`"
}
