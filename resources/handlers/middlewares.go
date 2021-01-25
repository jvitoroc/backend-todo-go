package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jvitoroc/todo-api/resources/repo"
)

func authenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") == false {
			respondWithMessage("Authorization token not provided.", http.StatusUnauthorized, w)
			return
		}

		claims, err := parseToken(authHeader[7:])
		if err != nil {
			respondWithMessage("Given authorization token is invalid or expired: "+err.Error(), http.StatusUnauthorized, w)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims["userId"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkActivationState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var user *repo.User
		userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)

		if user, err = repo.GetUser(db, userId); err != nil {
			if err == sql.ErrNoRows {
				respondWithMessage(fmt.Sprintf(MSG_NOT_FOUND_ERROR, "User", userId), http.StatusNotFound, w)
			} else {
				respondWithMessage("An error ocurred while verifying your account: "+err.Error(), http.StatusUnauthorized, w)
			}
			return
		}

		if !user.Active {
			respond(
				map[string]interface{}{
					"message":      "Your account is not yet verified. Check your mailbox.",
					"userVerified": false,
					"data":         user,
				},
				http.StatusForbidden,
				w,
			)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
