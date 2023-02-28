package middlewares

import (
	"errors"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"net/http"
)

func CheckUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(ContextUserKey).(*models.Claims).Role
		if role == models.User {
			next.ServeHTTP(w, r)
		} else {
			utils.RespondError(w, http.StatusUnauthorized, errors.New("you are not user"), true, "You are not admin")
			return
		}
	})
}
