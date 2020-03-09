package essentials

import (
	"net/http"
	u "lens/utils"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
	"qasmoke/models"
	"fmt"
	"os"
	"context"
)

var JwtAuthentication = func(next http.Handler) (http.Handler){
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)){
		noAuth := []string{"/api/login"}
		requestPath := r.URL.Path

		for _, value := range noAuth {
			if value == requestPath{
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		token := r.Header.get("Authorization")

		if token == "" {
			response = u.Message(false, "No auth token provided")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("content-type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenSplitted := string.Split(token, " ")
		if len(tokenSplitted) != 2 {
			response = u.Message(false, "token is not in the correct format")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("content-type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] 
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { 
			response = u.Message(false, "token is not in the correct format")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { 
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		
		fmt.Sprintf("User %", tk.Username) 
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) 
	}
}