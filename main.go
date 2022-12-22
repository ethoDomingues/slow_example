package main

import (
	"fmt"

	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/api"
	"github.com/ethodomingues/slow_example/cdn"
	"github.com/ethodomingues/slow_example/models"
)

func main() {
	app := slow.NewApp()
	app.SecretKey = "secret From Environment Variables" // Security: get from 'Environment Variables' => os.Getenv("SECRETKEY")
	app.Servername = "boatchazul.com.br"
	app.TemplateFolder = "front/templates/" // the default is '/template', but ...
	app.StaticFolder = "front/assets/"      // the default is '/assets', but ...
	app.BeforeRequest = beforeRequest

	app.GET("/", home)

	app.Mount(api.Load())
	app.Mount(cdn.Load())

	models.Session().AutoMigrate(
		&models.User{},
		&models.Cdn{},
		&models.Profile{},
		&models.React{},
		&models.Comm{},
		&models.Post{},
		&models.Solicitation{},
	)
	app.Listen(":5000")
}

func home(ctx *slow.Ctx) {
	rsp := ctx.Response
	rsp.RenderTemplate("/home/home.html")
}

func beforeRequest(ctx *slow.Ctx) {
	authAPI.GetCurrentUser(ctx)
	fmt.Println(ctx.Global["userID"])
	fmt.Println(ctx.Request.Header)
	if u, ok := ctx.Global["userID"]; ok {
		db := models.Session()
		user := &models.User{}
		db.Find(user, "uid = ?", u)
		if user.UID == "" {
			if d, ok := ctx.Global["userData"]; ok {
				data := d.(map[string]any)

				user.UID = ctx.Global["userID"].(string)
				user.Name = data["name"].(string)
				user.Email = data["email"].(string)
				db.Create(user)
			} else {
				ctx.Global["userInfo"] = "need create account"
			}
		}
		ctx.Global["user"] = user
	}
}
