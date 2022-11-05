package model

import (
	"fmt"
)

func NewReact(objID, owner string) (*React, bool) {
	var deleted bool
	var react = &React{}

	db := GetDB()
	_, found := FindByID(objID)
	if !found {
		return nil, false
	}
	db.Where("Obj = ? And Owner = ?", objID, owner).Find(react)
	if react.ID > 0 {
		db.Delete(react)
		deleted = true
	} else {
		react.Obj = objID
		react.Owner = owner
		db.Create(react)
		deleted = false
	}
	if deleted {
		return nil, deleted
	}
	return react, false
}

type React struct {
	Model
	Obj   string
	Owner string
}

func (r *React) ToJson() map[string]any {
	return map[string]any{
		"id":        r.UID(),
		"obj":       r.Obj,
		"owner":     r.Owner,
		"createdAt": r.Created(),
	}
}

func (r *React) Delete() {
	db := GetDB()
	db.Delete(r)
}

func (r *React) UID() string { return "reacts@" + fmt.Sprint(r.ID) }
