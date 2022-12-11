package solicitations

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Url:     "/users/{userID}/solicitations",
		Name:    "setSol",
		Methods: []string{"GET", "POST"},
		Func:    get,
	},
}

func get(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	userID := rq.Args["userID"]
	if userID != user.UID {
		rsp.Unauthorized()
	}
	sols := []*models.Solicitation{}
	models.Session().Where("rec = ? AND req = ?", userID, userID).Find(sols)
	slSol := models.SliceSol(sols)
	rsp.JSON(slSol.ToMap(userID, rq), 200)
}
