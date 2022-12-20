package cdn

import (
	"fmt"

	"github.com/ethodomingues/authAPI"
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/models"
)

var routes = []*slow.Route{
	{
		Name:    "set",
		Url:     "/",
		Func:    setCdn,
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
	authAPI.Required(ctx)
	u := ctx.Global["user"].(*models.User)
	f := ctx.Request.Files["images"][0]
	cdn := models.CreateCdn(f, fmt.Sprint(u.ID))
	ctx.Response.JSON(map[string]int{"id": cdn.ID}, 201)
}

func getCdn(ctx *slow.Ctx) {
	id := ctx.Request.Args["id"]
	filename := ctx.Request.Args["filename"]

	cdn := models.FindOr404(id, "*models.Cdn", "filename = ?", filename).(*models.Cdn)
	ctx.Response.Header.Set("Content-Type", cdn.ContentType)
	ctx.Response.Body.Write(cdn.Blob)
}

func getCdnByID(ctx *slow.Ctx) {
	id := ctx.Request.Args["id"]
	cdn := models.FindOr404(id, "*models.Cdn").(*models.Cdn)
	ctx.Response.Header.Set("Content-Type", cdn.ContentType)
	ctx.Response.Body.Write(cdn.Blob)
}
