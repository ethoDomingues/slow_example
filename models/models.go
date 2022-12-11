package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethodomingues/slow"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) Created() string {
	return m.CreatedAt.Format(time.RFC3339)
}
func (m *Model) Updated() string {
	return m.UpdatedAt.Format(time.RFC3339)
}

type DB struct {
	*gorm.DB
	Tables map[string]any
}

func Session() *DB {
	var err error
	if db == nil {
		db, err = gorm.Open(sqlite.Open("database.db"))
		if err != nil {
			panic(err)
		}
	}
	return &DB{DB: db.Session(&gorm.Session{})}
}

func FindByID(id string, conds ...any) (obj any, found bool) {
	var db = Session()
	_id := strings.Split(id, "@")
	if len(_id) != 2 {
		return nil, false
	}
	switch _id[0] {
	default:
		return nil, false
	case "cdns":
		obj = &Cdn{}
	case "comms":
		obj = &Comm{}
	case "posts":
		obj = &Post{}
	case "reacts":
		obj = &React{}
	case "users":
		obj = &User{}
	}
	res := db.DB.Where("ID = ?", _id[1]).Find(obj, conds...)
	return obj, res.RowsAffected > 0
}

func FindOr(id string, typ string, code int, conds ...any) any {
	obj, found := FindByID(id, conds...)
	if !found || slow.TypeOf(obj) != typ {
		panic("abort:" + fmt.Sprint(code))
	}
	return obj
}

func FindOr404(id string, typ string, conds ...any) any {
	obj, found := FindByID(id, conds...)
	if !found {
		panic("abort:404")
	}
	if slow.TypeOf(obj) != typ {
		panic("abort:400")
	}
	return obj
}
