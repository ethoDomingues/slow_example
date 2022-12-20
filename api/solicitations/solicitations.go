package solicitations

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Url:     "/users/{userID}/solicitations",
		Name:    "solicitations",
		Methods: []string{"GET", "POST", "DELETE"},
		MapCtrl: slow.MapCtrl{
			"GET":    {Func: getSolicitations},
			"POST":   {Func: post},
			"DELETE": {Func: delete},
		},
	},
	slow.GET("/users/{userID}/solicitations/request", getRequestSolicitations),
	slow.GET("/users/{userID}/solicitations/received", getReceivedSolicitations),
}

func validateUser(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	userID := rq.Args["userID"]
	if userID != user.UID {
		rsp.Unauthorized()
	}
}

func getSolicitations(ctx *slow.Ctx) {
	validateUser(ctx)
	rq := ctx.Request
	rsp := ctx.Response
	user := ctx.Global["user"].(*models.User)
	sols := map[string]any{
		"request":  user.RequestSolicitations(rq),
		"received": user.ReceivedSolicitations(rq),
	}
	rsp.JSON(sols, 200)
}

func getReceivedSolicitations(ctx *slow.Ctx) {
	validateUser(ctx)
	rq := ctx.Request
	rsp := ctx.Response
	user := ctx.Global["user"].(*models.User)
	rsp.JSON(user.ReceivedSolicitations(rq), 200)
}

func getRequestSolicitations(ctx *slow.Ctx) {
	validateUser(ctx)
	rq := ctx.Request
	rsp := ctx.Response
	user := ctx.Global["user"].(*models.User)
	rsp.JSON(user.RequestSolicitations(rq), 200)
}

func post(ctx *slow.Ctx) {
	validateUser(ctx)
	rq := ctx.Request
	rsp := ctx.Response
	db := models.Session()

	user := ctx.Global["user"].(*models.User)
	userRec := rq.Form["user"].(string)
	if r := db.Find(&models.User{}, "uid = ?", userRec); r.RowsAffected == 0 {
		rsp.NotFound()
	}
	sol := &models.Solicitation{}
	if r := db.Where("rec = ? AND req = ?", user.UID, userRec).Find(sol); r.RowsAffected > 0 {
		rsp.JSON(sol.ToMap(), 200)
	} else if r := db.Where("req = ? AND rec = ?", user.UID, userRec).Find(sol); r.RowsAffected > 0 {
		rsp.JSON(sol.ToMap(), 200)
	}
	s := &models.Solicitation{
		Req: user.UID,
		Rec: userRec,
	}
	db.Create(s)
	rsp.JSON(s.ToMap(), 201)
}

func delete(ctx *slow.Ctx) {
	validateUser(ctx)
	rq := ctx.Request
	rsp := ctx.Response
	db := models.Session()

	user := ctx.Global["user"].(*models.User)
	userRec := rq.Form["user"].(string)
	if r := db.Find(&models.User{}, "uid = ?", userRec); r.RowsAffected == 0 {
		rsp.NotFound()
	}
	sol := &models.Solicitation{}
	if r := db.Where("req = ? AND rec = ?", user.UID, userRec).Find(sol); r.RowsAffected == 0 {
		rsp.NotFound()
	}
	db.Delete(sol)
	rsp.JSON(nil, 204)
}
