package model

import (
	"fmt"

	"github.com/ethodomingues/slow"
)

func NewSliceCdn(owner string, fs ...*slow.File) []*Cdn {
	cdns := []*Cdn{}
	for _, f := range fs {
		cdns = append(cdns, NewCdn(owner, f))
	}
	return cdns
}

func NewCdn(owner string, f *slow.File) *Cdn {
	return &Cdn{
		Owner:       owner,
		Blob:        f.Stream.Bytes(),
		Filename:    f.Filename,
		ContentType: f.ContentType,
	}
}

func CreateCdn(f *slow.File, owner string) *Cdn {
	cdn := NewCdn(owner, f)
	db := GetDB()
	db.Create(cdn)
	return cdn
}

type Cdn struct {
	Model
	Owner string

	Blob        []byte
	Filename    string
	ContentType string
}

func (c *Cdn) Url(rq *slow.Request) string {
	return rq.UrlFor("cdn.get", true, map[string]string{"id": c.UID(), "filename": c.Filename})
}

func (c *Cdn) Query(conds ...any) {
	db := GetDB()
	db.Find(c, conds...)
}

func (c *Cdn) Delete() {
	db := GetDB()
	pf := &Profile{}
	db.Where("cdn = ?", c.UID()).Delete(pf)
	db.Delete(c)
}

func (c *Cdn) UID() string { return "cdns@" + fmt.Sprint(c.ID) }

func (c *Cdn) ToJSON(rq *slow.Request) map[string]any {
	id := c.UID()
	return map[string]any{
		"id":        id,
		"url":       c.Url(rq),
		"owner":     c.Owner,
		"createdAt": c.Created(),
	}
}
