package solicitations

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Url:     "/users/{userID}/solicitations",
		Name:    "setSol",
		Methods: []string{"GET", "POST"},
		Func:    auth.Manager(get, true),
	},
}

func get(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	userID := rq.Args["userID"]
	if userID != user.UID() {
		rsp.Unauthorized()
	}
	sols := []*model.Solicitation{}
	model.GetDB().Where("rec = ? AND req = ?", userID, userID).Find(sols)
	slSol := model.SliceSol(sols)
	rsp.JSON(slSol.ToJSON(userID, rq), 200)
}
