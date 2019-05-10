package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
)

type Counter struct {
	Entries int
}

func GenericPaginateQuery(table interface{}, criteria map[string]string, page interface{}, perPage interface{}) (*gorm.DB, int, int, int) {
	var currentPage, totalPages, entriesPerPage int
	var counter Counter

	currentPage = handleCurrentPage(page)
	entriesPerPage = handleEntriesPerPage(perPage)

	query, values := BuildSearchEngine(table, criteria, "AND")

	db.Table(TableName(table)).Select("COUNT(*) AS entries").Where(query, values...).Scan(&counter)

	totalPages = counter.Entries / entriesPerPage
	if (counter.Entries % entriesPerPage) > 0 {
		totalPages++
	}

	offset := (currentPage - 1) * entriesPerPage
	pagination := db.Offset(offset).Limit(entriesPerPage).Where(query, values...)

	return pagination, currentPage, totalPages, counter.Entries
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
}
