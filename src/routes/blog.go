package routes

import (
	"net/http"

	"github.com/Sahil2k07/Blog-App-Go/src/controllers"
	"github.com/Sahil2k07/Blog-App-Go/src/middlewares"
)

func BlogRoutes(router *http.ServeMux) {

	router.HandleFunc("/blog/get-blog/{id}", controllers.GetBlog)
	router.HandleFunc("/blog/get-blogs", controllers.GetAllBlogs)

	// Authenticated Routes
	router.Handle("/blog/user-blogs", middlewares.Auth(http.HandlerFunc(controllers.GetUserBlogs)))
	router.Handle("/blog/update-blog/{id}", middlewares.Auth(http.HandlerFunc(controllers.UpdateBlog)))
	router.Handle("/blog/delete-blog/{id}", middlewares.Auth(http.HandlerFunc(controllers.DeleteBlog)))
	router.Handle("/blog/create-blog", middlewares.Auth(http.HandlerFunc(controllers.CreateBlog)))

}
