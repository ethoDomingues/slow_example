package react

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Url:     "/users/{userID}/posts/{postID}/reacts",
		Name:    "reactPost",
		Methods: []string{"GET", "PUT"},
		Func:    reactPost,
	},
	{
		Url:     "/users/{userID}/posts/{postID}/comments/{commID}/reacts",
		Name:    "reactComment",
		Methods: []string{"GET", "PUT"},
		Func:    reactComment,
	},
}

func reactPost(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response
	user := ctx.Global["user"].(*models.User)
	postID := rq.Args["postID"]
	userID := rq.Args["userID"]

	models.FindOr404(postID, "*models.Post", "owner = ?", userID)

	if ctx.Request.Method == "PUT" {
		if user.UID != userID {
			rsp.BadRequest()
		}
		models.NewReact(postID, user.UID)
	}
	rs := []*models.React{}
	rsJson := []map[string]any{}
	models.Session().Where("obj = ?", postID).Find(&rs)
	for _, r := range rs {
		rsJson = append(rsJson, r.ToMap())
	}
	rsp.JSON(rsJson, 200)
}

func reactComment(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	postID := rq.Args["postID"]
	commID := rq.Args["commID"]
	userID := rq.Args["userID"]

	models.FindOr404(commID, "*models.Comm", "owner = ? AND post = ?", userID, postID)

	if ctx.Request.Method == "PUT" {
		if user.UID != userID {
			rsp.BadRequest()
		}
		models.NewReact(commID, user.UID)
	}
	rs := []models.React{}
	rsJson := []map[string]any{}
	models.Session().Where("obj = ?", commID).Find(&rs)
	for _, r := range rs {
		rsJson = append(rsJson, r.ToMap())
	}
	rsp.JSON(rsJson, 200)
}
