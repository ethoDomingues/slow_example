package friends

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Name: "friends",
		Url:  "/friends",
		MapCtrl: slow.MapCtrl{
			"PUT": {Func: put},
		},
	},
}

func put(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	userRec, ok := rq.Form["user"].(string)

	if !ok || userRec == "" {
		rsp.BadRequest()
	}

	f := models.Friends{
		Req: user.UID,
		Rec: userRec,
	}
	db := models.Session()
	db.Create(f)

}
