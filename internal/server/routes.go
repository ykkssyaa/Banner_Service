package server

import (
	"BannerService/internal/consts"
	"BannerService/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes struct {
	routesList   []Route
	allowedRoles []string
}

// createSubRouter Создание ветки роутов со своим middleware
func createSubRouter(router *mux.Router, routes Routes) {

	subRouter := router.PathPrefix("/").Subrouter()

	// Применяем middleware для подроутера
	subRouter.Use(middleware.TokenMiddleware(routes.allowedRoles...))

	for _, route := range routes.routesList {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		subRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

func NewRouter(server *HttpServer) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	var adminRoutes = Routes{
		[]Route{
			{
				"BannerGet",
				strings.ToUpper("Get"),
				"/banner",
				server.BannerGet,
			},

			{
				"BannerIdDelete",
				strings.ToUpper("Delete"),
				"/banner/{id}",
				server.BannerIdDelete,
			},

			{
				"BannerIdPatch",
				strings.ToUpper("Patch"),
				"/banner/{id}",
				server.BannerIdPatch,
			},

			{
				"BannerPost",
				strings.ToUpper("Post"),
				"/banner",
				server.BannerPost,
			},
			{
				"GetBannerVersions",
				strings.ToUpper("Get"),
				"/banner/{id}",
				server.GetBannerVersions,
			},
			{
				"SetBannerVersion",
				strings.ToUpper("Post"),
				"/banner/{id}",
				server.SetBannerVersion,
			},
		},
		[]string{consts.AdminRole},
	}
	var userRoutes = Routes{
		[]Route{
			{
				"UserBannerGet",
				strings.ToUpper("Get"),
				"/user_banner",
				server.UserBannerGet,
			},
		},
		[]string{
			consts.AdminRole,
			consts.UserRole,
		},
	}

	createSubRouter(router, adminRoutes)
	createSubRouter(router, userRoutes)

	return router
}
