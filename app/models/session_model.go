package models

import (
	"strconv"
	"time"
	"wheel.smart26.com/utils"
)

type Session struct {
	ID            uint       `gorm:"primary_key"`
	UserID        uint       `gorm:"index"`
	Jti           string     `gorm:"type:varchar(255);unique_index"`
	App           string     `gorm:"type:varchar(255)"`
	Requests      int        `gorm:"not null;default:0"`
	ExpiresIn     int        `gorm:"not null;default:0"`
	Address       string     `gorm:"type:varchar(255)"`
	Active        bool       `gorm:"default:true"`
	LastRequestAt *time.Time `gorm:"default:null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ExpiresAt     time.Time
}

func sessionValidate(session *Session) bool {
	return true
}

func SessionUpdate(session *Session) bool {
	var id string
	mapUpdate := make(map[string]interface{})
	id = strconv.FormatUint(uint64(session.ID), 10)
	currentSession := SessionFind(id)

	if sessionValidate(session) && !db.NewRecord(session) {
		if currentSession.Requests != session.Requests {
			mapUpdate["requests"] = session.Requests
		}

		if currentSession.LastRequestAt != session.LastRequestAt {
			mapUpdate["last_request_at"] = session.LastRequestAt
		}

		if currentSession.Active != session.Active {
			mapUpdate["active"] = session.Active
		}

		if currentSession.Jti != session.Jti {
			mapUpdate["jti"] = session.Jti
		}

		if len(mapUpdate) > 0 {
			db.Model(&session).Updates(mapUpdate)
		}

		return true
	} else {
		Errors = append(Errors, "was not found")
		return false
	}
}

func SessionCreate(session *Session) bool {
	utils.LoggerInfo().Println("models: SessionCreate")

	if sessionValidate(session) && db.NewRecord(session) {
		db.Create(&session)

		if db.NewRecord(session) {
			Errors = append(Errors, "database error")
			return false
		}
		return true
	} else {
		return false
	}
}

func SessionSave(session *Session) bool {
	if db.NewRecord(session) {
		return SessionCreate(session)
	} else {
		return SessionUpdate(session)
	}
}

func SessionFind(id interface{}) Session {
	var session Session

	db.First(&session, id)

	return session
}

func SessionFindByToken(jti string) Session {
	var session Session

	db.Where("jti = ?", jti).First(&session)

	return session
}

func SessionDestroy(session *Session) bool {
	if db.NewRecord(session) {
		return false
	} else {
		db.Delete(&session)
		return true
	}
}

func SessionDeactivate(session *Session) bool {
	session.Active = false
	return SessionSave(session)
}

func SessionIncrementStats(session *Session) bool {
	t := time.Now()
	session.LastRequestAt = &t
	session.Requests = session.Requests + 1
	return SessionSave(session)
}
