package pagination

var Path = []string{"commons", "app", "model", "pagination", "pagination.go"}

var Content = `package pagination

import (
	"{{ .AppRepository }}/commons/app/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
)

type Counter struct {
	Entries int
}

func Query(db *gorm.DB, table interface{}, page interface{}, perPage interface{}) (*gorm.DB, int, int, int) {
	var currentPage, totalPages, entriesPerPage int
	var counter Counter

	currentPage = handleCurrentPage(page)
	entriesPerPage = handleEntriesPerPage(perPage)

	db.Table(model.TableName(table)).Order("", true).Select("COUNT(*) AS entries").Scan(&counter)

	totalPages = counter.Entries / entriesPerPage
	if (counter.Entries % entriesPerPage) > 0 {
		totalPages++
	}

	offset := (currentPage - 1) * entriesPerPage
	db = db.Offset(offset).Limit(entriesPerPage)

	return db, currentPage, totalPages, counter.Entries
}

func handleCurrentPage(page interface{}) int {
	var currentPage int
	var err error

	switch auxPage := page.(type) {
	case int:
		currentPage = auxPage
	case string:
		currentPage, err = strconv.Atoi(auxPage)
		if err != nil {
			currentPage = 1
		}
	default:
		currentPage = 1
	}

	return currentPage
}

func handleEntriesPerPage(perPage interface{}) int {
	var entriesPerPage int
	var err error

	switch auxPerPage := perPage.(type) {
	case int:
		entriesPerPage = auxPerPage
	case string:
		entriesPerPage, err = strconv.Atoi(auxPerPage)
		if err != nil {
			entriesPerPage = 20
		}
	default:
		entriesPerPage = 20
	}

	return entriesPerPage
}`
