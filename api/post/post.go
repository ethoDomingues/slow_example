package post

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var Routes = []*slow.Route{
	{
		Url:  "/posts",
		Name: "getAllPost",
		Func: post,
		Ctrl: slow.Ctrl{
			"get": {Func: getAll},
		},
	},
	{
		Url:  "/users/{userID}/posts",
		Name: "setPost",
		Func: post,
		Ctrl: slow.Ctrl{
			"post": {Func: auth.Manager(post, true)},
			"get":  {Func: getUserPosts},
		},
	},
	{
		Url:  "/users/{userID}/posts/{postID}",
		Name: "post",
		Func: get,
		Ctrl: slow.Ctrl{
			"get":    {Func: get},
			"delete": {Func: auth.Manager(delete, true)},
		},
	},
}

func get(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	postID := ctx.Request.Args["postID"]
	userID := rq.Args["userID"]

	p, found := model.FindByID(postID, "owner = ? AND deleted = false", userID)
	if post, ok := p.(*model.Post); found && ok {
		rsp.JSON(post.ToJSON(rq), 200)
	}
	rsp.NotFound()
}

func getUserPosts(ctx *slow.Ctx) {
	userID := ctx.Request.Args["userID"]
	model.FindOr404(userID, "*model.User")

	var posts = []*model.Post{}

	model.GetDB().Where("owner = ? AND deleted = false", userID).Find(&posts)
	ps := []map[string]any{}
	for _, p := range posts {
		ps = append(ps, p.ToJSON(ctx.Request))
	}
	ctx.Response.JSON(ps, 200)
}

func getAll(ctx *slow.Ctx) {
	posts := []*model.Post{}
	model.GetDB().Where("deleted = false").Order("id DESC").Find(&posts)
	js := []any{}
	for i := 0; i < len(posts); i++ {
		js = append(js, posts[i].ToJSON(ctx.Request))
	}
	ctx.Response.JSON(js, 200)
}

func post(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)

	if rq.Args["userID"] != user.UID() {
		rsp.BadRequest()
	}

	var text string
	var pShared *model.Post
	_text := rq.Form["text"]
	shared := rq.Form["shared"]
	profile := rq.Form["profile"]
	images := rq.Files["images"]
	if sh, ok := shared.(string); ok {
		pShared = model.FindOr404(sh, "*model.Post").(*model.Post)
		profile = nil
		images = nil

	}

	if slow.TypeOf(_text) == "string" || _text == "" {
		text = _text.(string)
	}
	if (text == "") && (images == nil || len(images) < 1) && pShared == nil {
		rsp.BadRequest()
	}
	sharedID := ""
	if pShared != nil {
		sharedID = pShared.UID()
		if pShared.Shared != "" {
			sharedID = pShared.Shared
		}
	}
	post := model.CreatePost(text, user.UID(), sharedID, images)

	if profile != nil {
		img := post.GetImages()[0]
		model.NewProfile(user.UID(), img.UID())
	}
	rsp.JSON(post.ToJSON(rq), 201)
}

// TODO: implemetar içaqi
func put(ctx *slow.Ctx) {}

func delete(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*model.User)
	if rq.Args["userID"] != user.UID() {
		rsp.BadRequest()
	}

	id := rq.Args["postID"]
	post := model.FindOr404(id, "*model.Post", "owner = ?", user.UID())
	post.(*model.Post).Delete()
	rsp.StatusCode = 204
}
