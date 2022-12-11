package post

import (
	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var Routes = []*slow.Route{
	{
		Url:  "/posts",
		Name: "getAllPost",
		Func: post,
		MapCtrl: slow.MapCtrl{
			"get": {Func: getAll},
		},
	},
	{
		Url:  "/users/{userID}/posts",
		Name: "setPost",
		Func: post,
		MapCtrl: slow.MapCtrl{
			"post": {Func: post},
			"get":  {Func: getUserPosts},
		},
	},
	{
		Url:  "/users/{userID}/posts/{postID}",
		Name: "post",
		Func: get,
		MapCtrl: slow.MapCtrl{
			"get":    {Func: get},
			"delete": {Func: delete},
		},
	},
}

func get(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	postID := ctx.Request.Args["postID"]
	userID := rq.Args["userID"]

	p, found := models.FindByID(postID, "owner = ? AND deleted = false", userID)
	if post, ok := p.(*models.Post); found && ok {
		rsp.JSON(post.ToMap(rq), 200)
	}
	rsp.NotFound()
}

func getUserPosts(ctx *slow.Ctx) {
	userID := ctx.Request.Args["userID"]
	models.FindOr404(userID, "*models.User")

	var posts = []*models.Post{}

	models.Session().Where("owner = ? AND deleted = false", userID).Find(&posts)
	ps := []map[string]any{}
	for _, p := range posts {
		ps = append(ps, p.ToMap(ctx.Request))
	}
	ctx.Response.JSON(ps, 200)
}

func getAll(ctx *slow.Ctx) {
	posts := []*models.Post{}
	models.Session().Where("deleted = false").Order("id DESC").Find(&posts)
	js := []any{}
	for i := 0; i < len(posts); i++ {
		js = append(js, posts[i].ToMap(ctx.Request))
	}
	ctx.Response.JSON(js, 200)
}

func post(ctx *slow.Ctx) {
	authAPI.Required(ctx)

	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)

	if rq.Args["userID"] != user.UID {
		rsp.BadRequest()
	}

	var text string
	var pShared *models.Post
	_text := rq.Form["text"]
	shared := rq.Form["shared"]
	profile := rq.Form["profile"]
	images := rq.Files["images"]

	if sh, ok := shared.(string); ok {
		pShared = models.FindOr404(sh, "*models.Post").(*models.Post)
		profile = nil
		images = nil
	}
	if (len(images) < 1) || (profile != nil && profile != "true") {
		profile = nil
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
	post := models.CreatePost(text, user.UID, sharedID, images)

	if profile == "true" {
		img := post.GetImages()[0]
		models.NewProfile(user.UID, img.UID())
	}
	rsp.JSON(post.ToMap(rq), 201)
}

// TODO: implemetar içaqi
// func put(ctx *slow.Ctx) {}

func delete(ctx *slow.Ctx) {
	authAPI.Required(ctx)
	rq := ctx.Request
	rsp := ctx.Response

	user := ctx.Global["user"].(*models.User)
	if rq.Args["userID"] != user.UID {
		rsp.BadRequest()
	}

	id := rq.Args["postID"]
	post := models.FindOr404(id, "*models.Post", "owner = ?", user.UID)
	post.(*models.Post).Delete()
	rsp.StatusCode = 204
}
