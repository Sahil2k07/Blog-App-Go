package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sahil2k07/Blog-App-Go/src/utils"
)

type UserAuthDetails struct {
	Id        string
	Email     string
	ProfileId string
	Verified  bool
}

var UserContext = &struct{}{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		var tokenString string

		if err == nil {
			tokenString = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			utils.UnAuthorized(w, "Missing or Un-Authorized Token")
			return
		}

		// Validate the token using the ValidateJWT function
		id, email, profileId, verified, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.InternalServerError(w, err.Error())
			return
		}

		if !verified {
			utils.UnAuthorized(w, "User's email is not verified")
			return
		}

		user := &UserAuthDetails{
			Id:        id,
			Email:     email,
			ProfileId: profileId,
			Verified:  verified,
		}

		ctx := context.WithValue(r.Context(), UserContext, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
