package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		switch route.Name {
		case "UserBannerGet":
			handler = AuthMiddleware(userOrAdminAccessCheck)(handler)
		case "BannerGet", "BannerPost", "BannerIdDelete", "BannerIdPatch":
			handler = AuthMiddleware(adminAccessCheck)(handler)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"BannerGet",
		strings.ToUpper("Get"),
		"/banner",
		BannersGet,
	},

	Route{
		"BannerIdDelete",
		strings.ToUpper("Delete"),
		"/banner/{id}",
		BannerIdDelete,
	},

	Route{
		"BannerIdPatch",
		strings.ToUpper("Patch"),
		"/banner/{id}",
		BannerIdPatch,
	},

	Route{
		"BannerPost",
		strings.ToUpper("Post"),
		"/banner",
		BannerPost,
	},

	Route{
		"UserBannerGet",
		strings.ToUpper("Get"),
		"/user_banner",
		UserBannerGet,
	},
}
