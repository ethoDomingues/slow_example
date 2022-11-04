package model

import "github.com/ethodomingues/slow"

func Newfriends(rec, req string) *Friends {
	db := GetDB()
	if res := db.Model(&Solicitation{}).Find(""); res.RowsAffected < 1 {
		if res := db.Model(&Solicitation{}).Find(""); res.RowsAffected < 1 {
		}
	}
	f := &Friends{
		Req: req,
		Rec: rec,
	}
	return f
}

type Friends struct {
	Model
	Req string
	Rec string
}

func (f *Friends) Delete() {}

func (f *Friends) GetAllFor(user string) []*User {
	friends := []*Friends{}

	db := GetDB()
	db.Where("Rec = ? AND Req = ?", user).Find(&friends)

	users := []*User{}
	for _, f := range friends {
		var u *User
		if f.Rec == user {
			u_, _ := FindByID(f.Req)
			u = u_.(*User)
		} else {
			u_, _ := FindByID(f.Rec)
			u = u_.(*User)
		}
		users = append(users, u)
	}
	return users
}

type SliceUser []*User

func (su *SliceUser) ToJSON(rq *slow.Request) map[string]any {
	users := []map[string]any{}
	ids := []string{}
	for _, u := range *su {
		ids = append(ids, u.UID())
		users = append(users, u.ToJSON(rq))
	}
	return map[string]any{
		"users": users,
		"ids":   ids,
	}
}
