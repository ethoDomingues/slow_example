package friends

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Name: "friends",
		Url:  "/friends",
		Ctrl: slow.Ctrl{
			"PUT": {Func: auth.Manager(put, true)},
		},
	},
}

func put(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	userRec, ok := rq.Form["user"].(string)

	if !ok || userRec == "" {
		rsp.BadRequest()
	}

	f := model.Friends{
		Req: user.UID(),
		Rec: userRec,
	}
	db := model.GetDB()
	db.Create(f)

}
