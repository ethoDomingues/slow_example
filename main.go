package main

import (
	"strings"

	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/api"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/cdn"
	"github.com/ethodomingues/slow_example/model"
)

func main() {
	app := slow.NewApp()
	app.SecretKey = "secret From Environment Variables" // Security: get from 'Environment Variables' => os.Getenv("SECRETKEY")
	app.TemplateFolder = "front/templates/"             // the default is '/template', but ...
	app.StaticFolder = "front/assets/"                  // the default is '/assets', but ...
	app.BeforeRequest = GetCurrentUser

	app.Get("/", home)
	app.Get("/login", login)
	app.Get("/profile", profile)

	app.Mount(cdn.Load())
	app.Mount(api.Load())
	app.Mount(auth.Load())

	db := model.GetDB()
	db.AutoMigrate(
		&model.User{},                // 1º
		&model.Cdn{}, &model.React{}, // 2º
		&model.Profile{}, // 3º
		&model.Comm{},    // 4º
		&model.Post{},    // 5º
	)

	app.ShowRoutes() // show all routes
	app.Listen()     // start the Listener
}

func home(ctx *slow.Ctx) {
	rsp := ctx.Response                   // the current request
	rsp.RenderTemplate("/home/home.html") // app.TemplateFolder + "/home/home.html" => front/templates/home/home.html
}

func login(ctx *slow.Ctx) {
	rsp := ctx.Response
	rsp.RenderTemplate("/auth/auth.html")
}

func profile(ctx *slow.Ctx) {
	rsp := ctx.Response
	rsp.RenderTemplate("/profile/profile.html")
}

func GetCurrentUser(ctx *slow.Ctx) {
	rq := ctx.Request

	if _, ok := ctx.Global["user"]; ok {
		return
	} else if tkn := rq.Raw.Header.Get("Authorization"); tkn != "" {
		if ntkn := strings.TrimPrefix(tkn, "Bearer "); ntkn != tkn {
			if j, ok := slow.ValidJWT(ntkn, ctx.App.SecretKey); ok {
				if u, ok := j.Payload["sub"]; ok {
					user, found := model.FindByID(u)
					if found {
						ctx.Global["user"] = user
						ctx.Global["token"] = j
						return
					}
				}
			} else if j != nil {
				if _, ok := j.Payload["_permanent"]; ok {
					if u, ok := j.Payload["sub"]; ok {
						user, found := model.FindByID(u)
						if found {
							ctx.Global["user"] = user
							ctx.Global["token"] = j
							return
						}
					}
				}
			}
		}
	} else if u, ok := ctx.Session.Get("sub"); ok {
		user, found := model.FindByID(u)
		if found {
			ctx.Global["user"] = user
			return
		}
	}
	return
}
