package react

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Url:     "/users/{userID}/posts/{postID}/reacts",
		Name:    "reactPost",
		Methods: []string{"GET", "PUT"},
		Func:    auth.Manager(reactPost, true),
	},
	{
		Url:     "/users/{userID}/posts/{postID}/comments/{commID}/reacts",
		Name:    "reactComment",
		Methods: []string{"GET", "PUT"},
		Func:    auth.Manager(reactComment, true),
	},
}

func reactPost(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	postID := rq.Args["postID"]
	userID := rq.Args["userID"]

	model.FindOr404(postID, "*model.Post", "owner = ?", userID)

	if ctx.Request.Method == "PUT" {
		if user.UID() != userID {
			rsp.BadRequest()
		}
		model.NewReact(postID, user.UID())
	}
	rs := []*model.React{}
	rsJson := []map[string]any{}
	model.GetDB().Where("obj = ?", postID).Find(&rs)
	for _, r := range rs {
		rsJson = append(rsJson, r.ToJson())
	}
	rsp.JSON(rsJson, 200)
}

func reactComment(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	postID := rq.Args["postID"]
	commID := rq.Args["commID"]
	userID := rq.Args["userID"]

	model.FindOr404(commID, "*model.Comm", "owner = ? AND post = ?", userID, postID)

	if ctx.Request.Method == "PUT" {
		if user.UID() != userID {
			rsp.BadRequest()
		}
		model.NewReact(commID, user.UID())
	}
	rs := []model.React{}
	rsJson := []map[string]any{}
	model.GetDB().Where("obj = ?", commID).Find(&rs)
	for _, r := range rs {
		rsJson = append(rsJson, r.ToJson())
	}
	rsp.JSON(rsJson, 200)
}
