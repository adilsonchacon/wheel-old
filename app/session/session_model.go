package session

import (
	"errors"
	"time"
	"wheel.smart26.com/app/entity"
	"wheel.smart26.com/commons/db"
)

const NotFound = "session was not found"

func Find(id interface{}) (entity.Session, error) {
	var session entity.Session
	var err error

	db.Conn.First(&session, id)
	if db.Conn.NewRecord(session) {
		err = errors.New(NotFound)
	}

	return session, err
}

func IsValid(session *entity.Session) (bool, []error) {
	return true, []error{}
}

func Update(session *entity.Session) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentSession, findErr := Find(session.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(session)

	if valid {
		columns := db.ColumnsFromTable(session, false)
		for _, column := range columns {
			newValue, _ = db.GetColumnValue(session, column)
			currentValue, _ = db.GetColumnValue(currentSession, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue
			}
		}

		if len(mapUpdate) > 0 {
			db.Conn.Model(&session).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(session *entity.Session) (bool, []error) {
	valid, errs := IsValid(session)
	if valid && db.Conn.NewRecord(session) {
		db.Conn.Create(&session)

		if db.Conn.NewRecord(session) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(session *entity.Session) (bool, []error) {
	if db.Conn.NewRecord(session) {
		return Create(session)
	} else {
		return Update(session)
	}
}

func Destroy(session *entity.Session) bool {
	if db.Conn.NewRecord(session) {
		return false
	} else {
		db.Conn.Delete(&session)
		return true
	}
}

func FindByJti(jti string) (entity.Session, error) {
	var session entity.Session
	var err error

	db.Conn.Where("jti = ?", jti).First(&session)
	if db.Conn.NewRecord(session) {
		session = entity.Session{}
		err = errors.New(NotFound)
	}

	return session, err
}

func Deactivate(session *entity.Session) (bool, []error) {
	session.Active = false
	return Save(session)
}

func IncrementStats(session *entity.Session) (bool, []error) {
	t := time.Now()
	session.LastRequestAt = &t
	session.Requests = session.Requests + 1
	return Save(session)
}
