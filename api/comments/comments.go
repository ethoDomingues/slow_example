package comments

import (
	"fmt"

	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Url:  "/users/{userID}/posts/{postID}/comments",
		Name: "setComent",
		MapCtrl: slow.MapCtrl{
			"POST": {Func: auth.Manager(post, true)},
		},
	},
	{
		Url:         "/users/{userID}/posts/{postID}/comments/{commID}",
		Name:        "getComent",
		Middlewares: slow.NewMiddleware(auth.Required),
		MapCtrl: slow.MapCtrl{
			"GET":    {Func: get},
			"PUT":    {Func: get},
			"DELETE": {Func: auth.Manager(delete, true)},
		},
	},
}

func get(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response
	uid := rq.Args["userID"]
	pid := rq.Args["postID"]
	id := rq.Args["commID"]
	c, find := model.FindByID(id, "owner = ? AND post = ?", uid, pid)
	if find {
		if comm, ok := c.(*model.Comm); ok {
			rsp.JSON(comm.ToJSON(rq), 200)
		}
	}
	rsp.NotFound()
}

// func put(ctx *slow.Ctx) {}

func post(ctx *slow.Ctx) {

	rq := ctx.Request
	rsp := ctx.Response
	cUser := ctx.Global["user"].(*model.User)
	user := rq.Args["userID"]
	post := rq.Args["postID"]

	var img *slow.File
	text, ok := rq.Form["text"].(string)
	if !ok {
		text = ""
	}

	image := rq.Files["image"]
	if (post == "") || (text == "" && len(image) < 1) {
		rsp.BadRequest()
	}

	model.FindOr404(post, "*model.Post", "owner = ?", user)

	if len(image) > 0 {
		img = image[0]
	}
	comm := model.NewComm(text, fmt.Sprint(cUser.UID()), post, img)

	rsp.JSON(comm.ToJSON(rq), 201)
}

func delete(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	id := rq.Args["commID"]

	if rq.Args["userID"] != user.UID() {
		rsp.BadRequest()
	}
	c := model.FindOr404(id, "*model.Comm", "Owner = ?", user.UID())
	c.(*model.Comm).Delete()
	rsp.StatusCode = 204
}
