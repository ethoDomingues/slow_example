package api

import (
	"github.com/ethodomingues/slow"
	"github.com/ethodomingues/slow_example/api/comments"
	"github.com/ethodomingues/slow_example/api/post"
	"github.com/ethodomingues/slow_example/api/react"
	"github.com/ethodomingues/slow_example/api/solicitations"
	"github.com/ethodomingues/slow_example/api/user"
)

func Load() *slow.Router {
	api := &slow.Router{
		Name:      "api",
		Prefix:    "/v1",
		Subdomain: "api",
	}

	api.AddAll(user.Routes...)
	api.AddAll(post.Routes...)
	api.AddAll(comments.Routes...)
	api.AddAll(react.Routes...)
	api.AddAll(solicitations.Routes...)

	api.Cors = &slow.Cors{
		AllowOrigin: `*`,
		AllowHeaders: []string{
			"Authorization",
			"X-Session-Token",
			"Content-Type"},
		ExposeHeaders: []string{"Authorization", "X-Session-Token"},
	}
	return api
}
