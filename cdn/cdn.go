package cdn

import (
	"fmt"

	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/auth"
	"github.com/ethodomingues/slow_example/model"
)

var routes = []*slow.Route{
	{
		Name:    "set",
		Url:     "/",
		Func:    auth.Manager(setCdn, true),
		Methods: []string{"POST", "PUT"},
	},
	{
		Name:    "get",
		Url:     "/{id}/{filename}",
		Func:    getCdn,
		Methods: []string{"GET"},
	},
	{
		Name:    "getByid",
		Url:     "/{id}",
		Func:    getCdnByID,
		Methods: []string{"GET"},
	},
}

func Load() *slow.Router {
	cdn := &slow.Router{
		Name:      "cdn",
		Prefix:    "/v1",
		Subdomain: "cdn",
	}
	cdn.Cors = &slow.Cors{
		AllowOrigin: "*",
	}
	cdn.AddAll(routes...)
	return cdn
}

func setCdn(ctx *slow.Ctx) {
	u := ctx.Global["user"].(*model.User)
	f := ctx.Request.Files["images"][0]
	cdn := model.CreateCdn(f, fmt.Sprint(u.ID))
	ctx.Response.JSON(map[string]int{"id": cdn.ID}, 201)
}

func getCdn(ctx *slow.Ctx) {
	id := ctx.Request.Args["id"]
	filename := ctx.Request.Args["filename"]

	cdn := model.FindOr404(id, "*model.Cdn", "filename = ?", filename).(*model.Cdn)
	ctx.Response.Headers.Set("Content-Type", cdn.ContentType)
	ctx.Response.Body.Write(cdn.Blob)
}

func getCdnByID(ctx *slow.Ctx) {
	id := ctx.Request.Args["id"]

	cdn := model.FindOr404(id, "*model.Cdn").(*model.Cdn)
	ctx.Response.Headers.Set("Content-Type", cdn.ContentType)
	ctx.Response.Body.Write(cdn.Blob)
}
