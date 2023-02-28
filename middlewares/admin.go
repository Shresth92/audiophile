package middlewares

import (
	"errors"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"net/http"
)

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(ContextUserKey).(*models.Claims).Role
		if role == models.Admin {
			next.ServeHTTP(w, r)
		} else {
			utils.RespondError(w, http.StatusUnauthorized, errors.New("you are not admin"), true, "You are not admin")
			return
		}
	})
}
