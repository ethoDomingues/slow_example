package user

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	slow.GET("/users/whoami", whoami),
	slow.GET("/users/{userID}", getUser),
}

func whoami(ctx *slow.Ctx) {
	authAPI.Required(ctx)
	if u, ok := ctx.Global["user"]; ok {
		user := u.(*models.User)
		ctx.Response.JSON(user.ToMap(ctx.Request), 200)
	}
	ctx.Response.Redirect("http://auth.boatchazul.com.br/v1/auth")
}

func getUser(ctx *slow.Ctx) {
	authAPI.Required(ctx)
	id := ctx.Request.Args["userID"]
	u := models.FindOr404(id, "*models.User").(*models.User)
	ctx.Response.JSON(u.ToMap(ctx.Request), 200)
}
