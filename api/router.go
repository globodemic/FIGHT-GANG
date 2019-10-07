package api

import (
	"net/http"

	mux "github.com/gorilla/mux"
)

//NewRouter creates routes
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
