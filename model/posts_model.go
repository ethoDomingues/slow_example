package model

import (
	"fmt"
	"strings"

	"github.com/ethodomingues/slow"
)

func CreatePost(text, owner string, shared string, images []*slow.File) *Post {
	p := &Post{
		Owner:  owner,
		Text:   text,
		Shared: shared,
	}

	db := GetDB()
	db.Create(p)
	var cdnIds strings.Builder
	for i, img := range images {
		cdn := CreateCdn(img, owner)
		cdnIds.WriteString(fmt.Sprint(cdn.ID))
		if i < (len(images) - 1) {
			cdnIds.WriteString(";")
		}
	}
	p.Images = cdnIds.String()
	db.Save(p)
	fmt.Println(p.Images)
	return p
}

type Post struct {
	Model
	Text,
	Images,
	Owner,
	Shared string
	Deleted bool
}

func (p *Post) UID() string { return "posts@" + fmt.Sprint(p.ID) }

func (p *Post) GetOwner() *User {
	owner, _ := FindByID(p.Owner)
	return owner.(*User)
}

func (p *Post) GetShared() *Post {
	if p.Shared == "" {
		return nil
	}
	post, _ := FindByID(p.Shared)
	return post.(*Post)
}

func (p *Post) GetImages() []*Cdn {
	if p.Images == "" {
		return nil
	}
	imgIDs := strings.Split(p.Images, ";")
	cdns := []*Cdn{}
	for _, id := range imgIDs {
		cdn := &Cdn{}
		cdn.Query("ID = ?", id)
		cdns = append(cdns, cdn)
	}
	return cdns
}

func (p *Post) GetReact() []*React {
	reacts := []*React{}
	GetDB().Where("Obj = ?", p.UID()).Order("id DESC").Find(&reacts)
	return reacts
}

func (p *Post) GetComm() []*Comm {
	comms := []*Comm{}
	GetDB().Where("Post = ?", p.UID()).Order("id DESC").Find(&comms)
	return comms
}

func (p *Post) SharedCount() int {
	var count int64
	GetDB().Model(&Post{}).Where("shared = ?", p.UID()).Count(&count)
	return int(count)
}

func (p *Post) ToJSON(rq *slow.Request) map[string]any {
	if p.Deleted {
		return map[string]any{
			"id":          p.UID(),
			"text":        "",
			"owner":       p.GetOwner().ToJSON(rq),
			"images":      []any{},
			"reacts":      []any{},
			"shared":      nil,
			"deleted":     true,
			"comments":    []any{},
			"createdAt":   p.Created(),
			"sharedCount": p.SharedCount(),
		}
	}
	imgs := []map[string]any{}
	for _, img := range p.GetImages() {
		imgs = append(imgs, img.ToJSON(rq))
	}

	comms := []map[string]any{}
	for _, comm := range p.GetComm() {
		comms = append(comms, comm.ToJSON(rq))
	}

	rs := p.GetReact()
	reacts := []map[string]any{}
	for _, r := range rs {
		reacts = append(reacts, r.ToJson())
	}
	var shMap map[string]any
	sh := p.GetShared()
	if sh != nil {
		shMap = sh.ToJSON(rq)
	}
	return map[string]any{
		"id":          p.UID(),
		"text":        p.Text,
		"owner":       p.GetOwner().ToJSON(rq),
		"images":      imgs,
		"reacts":      reacts,
		"shared":      shMap,
		"deleted":     p.Deleted,
		"comments":    comms,
		"createdAt":   p.Created(),
		"sharedCount": p.SharedCount(),
	}
}

func (p *Post) ToJSONbasic(rq *slow.Request) map[string]any {
	return map[string]any{
		"id":        p.UID(),
		"owner":     p.GetOwner().ToJSON(rq),
		"createdAt": p.Created(),
	}
}

func (p *Post) Delete() {
	imgs := p.GetImages()
	for _, img := range imgs {
		img.Delete()
	}
	comms := p.GetComm()
	for _, comm := range comms {
		comm.Delete()
	}
	reacts := p.GetReact()
	for _, react := range reacts {
		react.Delete()
	}
	p.Text = ""
	p.Images = ""
	p.Shared = ""
	p.Deleted = true
	GetDB().Save(p)
}
