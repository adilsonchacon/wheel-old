package ordering

var Path = []string{"commons", "app", "model", "ordering", "ordering.go"}

var Content = `package ordering

import (
	"github.com/jinzhu/gorm"
)

func Query(db *gorm.DB, order string) *gorm.DB {
	if order != "" {
		db = db.Order(order)
	}

	return db
}`
