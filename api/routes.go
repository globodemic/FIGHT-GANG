package api

import (
	"net/http"

	authorization "./auths"
	controller "./controllers"
	ws "./websocket"
)

//Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes slices
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Test",
		"GET",
		"/static",
		controller.TestPage,
	},
	Route{
		"CreatePlayer",
		"POST",
		"/api/players",
		controller.CreatePlayer,
	},
	Route{
		"GetPlayer",
		"GET",
		"/api/players/{id}",
		authorization.ValidateMiddleware(controller.GetPlayer),
	},
	Route{
		"DeletePlayer",
		"DELETE",
		"/api/players/{id}",
		authorization.ValidateMiddleware(controller.DeletePlayer),
	},
	Route{
		"PutPlayer",
		"PUT",
		"/api/players/{id}",
		authorization.ValidateMiddleware(controller.PutPlayer),
	},
	Route{
		"GetPlayers",
		"GET",
		"/api/players",
		authorization.ValidateMiddleware(controller.GetPlayers),
	},
	// Route{
	// 	"Foo",
	// 	"GET",
	// 	"/arena",
	// 	ws.ServeWs,
	// },
}
