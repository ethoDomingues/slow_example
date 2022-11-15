package auth

import (
	"encoding/base64"

	"github.com/ethodomingues/slow"
)

func Manager(f slow.Func, withCred bool) slow.Func {
	return func(ctx *slow.Ctx) {
		if withCred {
			Required(ctx)
			f(ctx)
		} else {
			if _, ok := ctx.Global["user"]; ok {
				if ctx.Request.ContentType == "application/json" {
					ctx.Response.JSON(map[string]string{"location": "/"}, 302)
				}
				ctx.Response.Redirect("/")
			}
		}
	}
}

func Required(ctx *slow.Ctx) {
	rq := ctx.Request
	rsp := ctx.Response

	if _, ok := ctx.Global["user"]; ok {
		return
	}

	n := base64.RawURLEncoding.EncodeToString([]byte(rq.RequestURL()))
	rsp.Redirect(
		slow.UrlFor("auth.login", true, map[string]string{"next": n}),
	)
}
