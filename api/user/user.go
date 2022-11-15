package user

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Url:         "/users/whoami",
		Name:        "getUser",
		Func:        whoami,
		Methods:     []string{"GET"},
		Middlewares: []slow.Func{auth.Required},
	},
	{
		Url:         "/users/{userID}",
		Name:        "getUser",
		Func:        getUser,
		Methods:     []string{"GET"},
		Middlewares: []slow.Func{auth.Required},
	},
}

func whoami(ctx *slow.Ctx) {
	if u, ok := ctx.Global["user"]; ok {
		var tkn string
		if j, ok := ctx.Global["token"]; ok {
			_, keep := j.(*slow.JWT).Payload["_permanent"]
			tkn = auth.GetTkn(ctx, keep)

		} else {
			tkn = auth.GetTkn(ctx, ctx.Session.Permanent)
		}
		ctx.Response.JSON(
			map[string]any{
				"user":  u.(*model.User).ToJSON(ctx.Request),
				"token": tkn,
			}, 200)
	}
	ctx.Response.Redirect(slow.UrlFor("auth.login", true, nil))
}

func getUser(ctx *slow.Ctx) {
	id := ctx.Request.Args["userID"]
	u := model.FindOr404(id, "*model.User").(*model.User)
	ctx.Response.JSON(u.ToJSON(ctx.Request), 200)
}
