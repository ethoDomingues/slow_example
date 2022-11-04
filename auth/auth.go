package auth

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/model"
)

var routes = []*slow.Route{
	{
		Url:     "/login",
		Name:    "login",
		Func:    login,
		Methods: []string{"POST"},
	},
	{
		Url:     "/register",
		Name:    "register",
		Func:    register,
		Methods: []string{"POST"},
	},
}

func Load() *slow.Router {
	auth := &slow.Router{
		Name:      "auth",
		Subdomain: "auth",
		Prefix:    "/v1",
		Cors: &slow.Cors{
			ExposeHeaders: []string{
				"Authorization",
			},
		},
	}
	auth.AddAll(routes...)
	auth.Cors = &slow.Cors{
		AllowOrigin:   "*",
		ExposeHeaders: []string{"Authorization"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
	}
	return auth
}

func login(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response
	keep := false

	user, pass, ok := rq.Raw.BasicAuth()
	if !ok {
		rsp.BadRequest()
	}

	if _keep, ok := rq.Form["keep"]; ok {
		if fmt.Sprint(_keep) == "true" {
			keep = true
		}
	}
	u := &model.User{}
	u.Query("Email = ?", user)
	if u.ID > 0 {
		if u.CheckPaswd(pass) {
			ctx.Global["user"] = u
			ctx.Session.Set("sub", u.UID())
			if keep {
				ctx.Session.Permanent = true
			}
			tkn := GetTkn(ctx, keep)
			nList := rq.Query["next"]

			if len(nList) > 0 {
				n, err := base64.RawURLEncoding.DecodeString(rq.Query["next"][0])
				if err == nil {
					rsp.JSON(map[string]any{
						"token":    tkn,
						"location": string(n),
					}, 200)
				}
			}
			rsp.JSON(map[string]any{"token": tkn}, 200)
		}
	}
	rsp.Unauthorized()
}

func register(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	var u *model.User
	var err error
	var name string
	var keep bool

	email, password, ok := rq.Raw.BasicAuth()

	profile := rq.Files["profile"]
	_name := rq.Form["name"]

	if _name == nil || !ok {
		rsp.BadRequest()
	} else {
		name = _name.(string)
	}

	if k, ok := rq.Form["keep"].(bool); ok {
		keep = k
	}

	if len(profile) > 0 {
		u, err = model.NewUser(map[string]any{"name": name, "email": email, "password": password}, profile[0])
	} else {
		rq.Form["email"] = email
		rq.Form["password"] = password
		u, err = model.NewUser(rq.Form, nil)
	}
	if err != nil {
		strErr := err.Error()
		switch strErr {
		default:
			rsp.BadRequest()
		case "Unauthorized":
			rsp.Unauthorized()
		}
	}
	db := model.GetDB()
	db.Create(u)
	ctx.Global["user"] = u

	if keep {
		ctx.Session.Set("sub", u.UID())
	}
	tkn := GetTkn(ctx, keep)
	nList := rq.Query["next"]

	if len(nList) > 0 {
		n, err := base64.RawURLEncoding.DecodeString(rq.Query["next"][0])
		if err == nil {
			rsp.JSON(map[string]any{
				"token":    tkn,
				"location": string(n),
			}, 200)
		}
	}
	rsp.JSON(map[string]any{"token": tkn}, 200)
}

func GetTkn(ctx *slow.Ctx, keep bool) string {
	u := ctx.Global["user"].(*model.User)

	j := slow.NewJWT()
	j.Secret = ctx.App.SecretKey

	j.Payload["sub"] = u.UID()

	if keep {
		j.Payload["exp"] = fmt.Sprint(time.Now().Add(time.Hour * 24 * 31).Unix())
		j.Payload["_permanent"] = "1"
	} else {
		j.Payload["exp"] = fmt.Sprint(time.Now().Add(time.Second * 5).Unix())
	}

	tkn := j.Sign()
	ctx.Response.Headers.Set("Authorization", "Bearer "+tkn)
	return "Bearer " + tkn
}
