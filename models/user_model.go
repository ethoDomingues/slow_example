package models

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
)

type User struct {
	authAPI.User
}

func (u *User) GetProfile() *Profile {
	var p Profile
	r := Session().Where("owner = ? AND active = ?", u.UID, true).Find(&p)
	if r.RowsAffected > 0 {
		return &p
	}
	return nil
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
