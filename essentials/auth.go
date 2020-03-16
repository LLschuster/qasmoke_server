package essentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	u "github.com/llschuster/qasmoke/utils"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuth := []string{"/api/v1/login"}
		requestPath := r.URL.Path

		for _, value := range noAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		token := r.Header.Get("Authorization")

		if token == "" {
			response = u.Message(false, "No auth token provided")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("content-type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenSplitted := strings.Split(token, " ")
		if len(tokenSplitted) != 2 {
			response = u.Message(false, "token is not in the correct format")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("content-type", "application/json")
			u.Respond(w, response)
			return
		}

		// tokenPart := tokenSplitted[1]
		// tk := &models.Token{}

		// token, err := jwt.ParseWithClaims(tokenPart, "729", func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("token_password")), nil
		// })
		var err interface{} = nil
		if err != nil {
			response = u.Message(false, "token is not in the correct format")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// if !token.Valid {
		// 	response = u.Message(false, "Token is not valid.")
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Header().Add("Content-Type", "application/json")
		// 	u.Respond(w, response)
		// 	return
		// }

		fmt.Sprintf("User ")
		ctx := context.WithValue(r.Context(), "user", 1)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
