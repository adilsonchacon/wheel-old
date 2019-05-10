package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"wheel.smart26.com/utils"
)

var db *gorm.DB
var Errors []string

func DbConnect() {
	var err error

	dbCconfig := loadDatabaseConfigFile()
	db, err = gorm.Open("postgres", stringfyDatabaseConfigFile(dbCconfig))

	if err != nil {
		utils.LoggerFatal().Println(err)
		panic("failed connect to database")
	} else {
		utils.LoggerInfo().Println("connected to the database successfully")
	}

	pool, err := strconv.Atoi(dbCconfig["pool"])
	if err != nil {
		utils.LoggerFatal().Println(err)
	} else {
		utils.LoggerInfo().Printf("database pool of connections: %d", pool)
	}

	db.DB().SetMaxIdleConns(pool)
}

func DbDisconnect() {
	defer db.Close()
}

func TableName(table interface{}) string {
	return db.NewScope(table).GetModelStruct().TableName(db)
}

func GetColumnType(table interface{}, columnName string) (string, error) {
	field, ok := db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Type().String(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func GetColumnValue(table interface{}, columnName string) (interface{}, error) {
	field, ok := db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Interface(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func SetColumnValue(table interface{}, columnName string, value string) error {
	field, ok := db.NewScope(table).FieldByName(columnName)

	if ok {
		columnType, _ := GetColumnType(table, columnName)
		valueInterface, _ := utils.ConvertStringToInterface(columnType, value)
		return field.Set(valueInterface)
	} else {
		return errors.New("column was not found")
	}
}

func ColumnsFromTable(table interface{}, all bool) []string {
	var columns []string
	fields := db.NewScope(table).Fields()

	for _, field := range fields {
		if !all && ((field.Names[0] == "Model") || (field.Relationship != nil)) {
			continue
		}
		columns = append(columns, field.DBName)
	}

	return columns
}

// PACKAGE METHODS

func loadDatabaseConfigFile() map[string]string {
	config := make(map[string]string)

	err := yaml.Unmarshal(readDatabaseConfigFile(), &config)
	if err != nil {
		utils.LoggerFatal().Printf("error: %v\n", err)
	}

	if config["pool"] == "" {
		config["pool"] = "5"
	}

	return config
}

func readDatabaseConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/database.yml")
	if err != nil {
		utils.LoggerFatal().Println(err)
	}

	return data
}

func stringfyDatabaseConfigFile(mapped map[string]string) string {
	var arr []string

	for key, value := range mapped {
		if key != "pool" {
			arr = append(arr, key+"='"+value+"'")
		}
	}

	return strings.Join(arr, " ")
}
