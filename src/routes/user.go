package routes

import (
	"net/http"

	"github.com/Sahil2k07/Blog-App-Go/src/controllers"
	"github.com/Sahil2k07/Blog-App-Go/src/middlewares"
)

func UserRoutes(router *http.ServeMux) {

	router.HandleFunc("/user/login", controllers.Login)

	// Authenticated Routes
	router.Handle("/user/update-profile", middlewares.Auth(http.HandlerFunc(controllers.UpdateProfile)))
	router.Handle("/user/get-profile", middlewares.Auth(http.HandlerFunc(controllers.GetProfile)))

}
