package routes

import (
	"net/http"

	"github.com/Sahil2k07/Blog-App-Go/src/controllers"
)

func AuthRoutes(router *http.ServeMux) {

	router.HandleFunc("/auth/resend-otp", controllers.ReSendOtp)
	router.HandleFunc("/auth/signup", controllers.SignUp)
	router.HandleFunc("/auth/verify-user", controllers.VerifyUser)

}
