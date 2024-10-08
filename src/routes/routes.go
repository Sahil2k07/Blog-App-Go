package routes

import (
	"net/http"

	"github.com/Sahil2k07/Blog-App-Go/src/controllers"
)

func AppRoutes() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("/", controllers.RootHandler)

	// Register all route handlers
	AuthRoutes(router)

	UserRoutes(router)

	BlogRoutes(router)

	return router

}
