package models

import (
	"fmt"

	"github.com/ethodomingues/slow"
)

func NewComm(text, owner, post string, image *slow.File) *Comm {
	c := &Comm{
		Text:  text,
		Owner: owner,
		Post:  post,
	}
	db := Session()
	db.Create(c)
	if image != nil {
		cdn := CreateCdn(image, owner)
		c.Image = cdn.UID()
		db.Save(c)
	}

	return c
}

type Comm struct {
	Model
	Text  string
	Image string
	Owner string
	Post  string
}

func (c *Comm) GetOwner() *User {
	u, _ := FindByID(c.Owner)
	return u.(*User)
}

func (c *Comm) GetPost() *Post {
	p, _ := FindByID(c.Post)
	return p.(*Post)
}

func (c *Comm) GetImage() *Cdn {
	if c.Image == "" {
		return nil
	}
	cdn, _ := FindByID(c.Image)
	return cdn.(*Cdn)
}

func (c *Comm) GetReact() []*React {
	reacts := []*React{}
	Session().Where("Obj = ?", c.UID()).Find(&reacts)
	return reacts
}

func (c *Comm) ToMap(rq *slow.Request) map[string]any {
	cdn := c.GetImage()
	var img map[string]any
	if cdn != nil {
		img = cdn.ToMap(rq)
	}
	rs := c.GetReact()
	reacts := []map[string]any{}
	for _, r := range rs {
		reacts = append(reacts, r.ToMap())
	}
	return map[string]any{
		"id":        c.UID(),
		"text":      c.Text,
		"post":      c.GetPost().ToJSONbasic(rq),
		"image":     img,
		"owner":     c.GetOwner().ToMap(rq),
		"createdAt": c.Created(),
		"reacts":    reacts,
	}
}

func (c *Comm) UID() string { return "comms@" + fmt.Sprint(c.ID) }

func (c *Comm) Delete() {
	db := Session()
	if c.Image != "" {
		c.GetImage().Delete()
	}
	reacts := c.GetReact()
	for _, react := range reacts {
		react.Delete()
	}
	db.Delete(c)
}
