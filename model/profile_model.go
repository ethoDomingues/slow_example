package model

import (
	"fmt"

	"github.com/ethodomingues/slow"
)

func NewProfile(user, cdn string) *Profile {
	p := &Profile{
		Owner:  user,
		Cdn:    cdn,
		Active: true,
	}
	db := GetDB()
	var p1 *Profile
	res := db.Where("owner = ? AND active = ?", user, true).Find(p1)
	if res.RowsAffected > 0 {
		p1.Active = false
		db.Save(p1)
	}
	db.Create(p)
	return p
}

type Profile struct {
	Model
	Cdn    string
	Owner  string
	Active bool
}

func (p *Profile) UID() string { return "profiles@" + fmt.Sprint(p.ID) }

func (p *Profile) Url(rq *slow.Request) string {
	_any, _ := FindByID(p.Cdn)
	cdn := _any.(*Cdn)
	return cdn.Url(rq)
}

func (p *Profile) ToJSON(rq *slow.Request) map[string]any {
	return map[string]any{
		"id":        p.UID(),
		"active":    p.Active,
		"createdAt": p.CreatedAt.UTC().String(),
		"owner":     p.Owner,
		"url":       p.Url(rq),
	}
}
