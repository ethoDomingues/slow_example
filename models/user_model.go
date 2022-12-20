package models

import (
	"github.com/ethodomingues/slow"
)

type User struct {
	Model

	Name,
	Email,
	UID string
}

func (u *User) GetProfile() *Profile {
	var p Profile
	r := Session().Where("owner = ? AND active = ?", u.UID, true).Find(&p)
	if r.RowsAffected > 0 {
		return &p
	}
	return nil
}

func (u *User) RequestSolicitations(rq *slow.Request) map[string]any {
	db := Session()
	req := map[string]any{}
	sols := []Solicitation{}

	db.Where("req = ?", u.UID, u.UID).Find(&sols)
	for _, sol := range sols {
		_user, _ := FindByID(sol.Rec)
		user := _user.(*User)
		req[user.UID] = user.ToMap(rq)
	}
	return req
}

func (u *User) ReceivedSolicitations(rq *slow.Request) map[string]any {
	db := Session()
	rec := map[string]any{}
	sols := []Solicitation{}

	db.Where("rec = ?", u.UID, u.UID).Find(&sols)
	for _, sol := range sols {
		_user, _ := FindByID(sol.Req)
		user := _user.(*User)
		rec[user.UID] = user.ToMap(rq)
	}
	return rec
}

func (u *User) ToMap(rq *slow.Request) map[string]any {
	var p map[string]any
	_p := u.GetProfile()
	if _p != nil {
		p = _p.ToMap(rq)
	}
	return map[string]any{
		"id":        u.UID,
		"name":      u.Name,
		"email":     u.Email,
		"createdAt": u.Created(),
		"updateAt":  u.Updated(),
		"profile":   p,
	}
}
