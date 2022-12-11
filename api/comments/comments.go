package comments

import (
	"fmt"

	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Url:  "/users/{userID}/posts/{postID}/comments",
		Name: "setComent",
		MapCtrl: slow.MapCtrl{
			"POST": {Func: post},
		},
	},
	{
		Url:         "/users/{userID}/posts/{postID}/comments/{commID}",
		Name:        "getComent",
		Middlewares: slow.NewMiddleware(authAPI.Required),
		MapCtrl: slow.MapCtrl{
			"GET":    {Func: get},
			"PUT":    {Func: get},
			"DELETE": {Func: delete},
		},
	},
}

func get(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response
	uid := rq.Args["userID"]
	pid := rq.Args["postID"]
	id := rq.Args["commID"]
	c, find := models.FindByID(id, "owner = ? AND post = ?", uid, pid)
	if find {
		if comm, ok := c.(*models.Comm); ok {
			rsp.JSON(comm.ToMap(rq), 200)
		}
	}
	rsp.NotFound()
}

// func put(ctx *slow.Ctx) {}

func post(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response
	cUser := ctx.Global["user"].(*models.User)
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

	models.FindOr404(post, "*models.Post", "owner = ?", user)

	if len(image) > 0 {
		img = image[0]
	}
	comm := models.NewComm(text, fmt.Sprint(cUser.UID), post, img)

	rsp.JSON(comm.ToMap(rq), 201)
}

func delete(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	id := rq.Args["commID"]

	if rq.Args["userID"] != user.UID {
		rsp.BadRequest()
	}
	c := models.FindOr404(id, "*models.Comm", "Owner = ?", user.UID)
	c.(*models.Comm).Delete()
	rsp.StatusCode = 204
}
