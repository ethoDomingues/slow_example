package dynamic

import (
	"fmt"

	"github.com/ethodomingues/slow"
)

var routes = []*slow.Route{
	{
		Url:  "/users/{user?}",
		Name: "get",
		Func: func(ctx *slow.Ctx) {
			fmt.Println(ctx.Request.Args)
		},
		Methods: []string{"GET"},
	},
}

func Load() *slow.Router {
	dyn := &slow.Router{
		Name: "dynamic",
		// Subdomain: "{dynamic}",
	}
	dyn.AddAll(routes...)
	return dyn
}
