package model

import (
	"errors"
	"fmt"

	"github.com/ethodomingues/slow"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(data map[string]any, prof *slow.File) (*User, error) {
	email, ok := data["email"].(string)
	if !ok {
		return nil, errors.New("bad request")
	}

	name, ok := data["name"].(string)
	if !ok {
		return nil, errors.New("bad request")
	}

	password, ok := data["password"].(string)
	if !ok {
		return nil, errors.New("bad request")
	}

	u := &User{}
	u.Query("email = ?", email)
	if u.ID > 0 {
		return nil, errors.New("Unathorized")
	}

	u.Name = name
	u.Email = email
	u.HashPasswd(password)

	if prof != nil {
		c := CreateCdn(prof, u.Email)
		NewProfile(u.UID(), c.UID())
	}

	return u, nil
}

type User struct {
	Model
	Name, Email string
	Password    []byte
}

func (u *User) UID() string { return "users@" + fmt.Sprint(u.ID) }

func (u *User) Query(conds ...any) {
	GetDB().Find(u, conds...)
}

func (u *User) HashPasswd(passwd string) {
	bPass, err := bcrypt.GenerateFromPassword([]byte(passwd), 14)
	if err != nil {
		panic(err)
	}
	u.Password = bPass
}

func (u *User) CheckPaswd(passwd string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(passwd)) == nil
}

func (u *User) GetProfile() *Profile {
	var p Profile
	r := GetDB().Where("owner = ? AND active = ?", u.UID(), true).Find(&p)
	if r.RowsAffected > 0 {
		return &p
	}
	return nil
}

func (u *User) ToJSON(rq *slow.Request) map[string]any {
	var p map[string]any
	_p := u.GetProfile()
	if _p != nil {
		p = _p.ToJSON(rq)
	}
	return map[string]any{
		"id":        u.UID(),
		"name":      u.Name,
		"email":     u.Email,
		"createdAt": u.Created(),
		"updateAt":  u.Updated(),
		"profile":   p,
	}
}
