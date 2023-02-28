package middlewares

import (
	"context"
	"errors"
	"github.com/Shresth92/audiophile/database/dbHelpers"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type ContextKey string

const ContextUserKey ContextKey = "userInfo"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		claims := &models.Claims{}
		if token == "" {
			utils.RespondError(w, http.StatusUnauthorized, errors.New("token not sent in header"), true, "token not sent in header")
			return
		} else {
			parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(utils.GetEnvValue("JwtKey")), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					utils.RespondError(w, http.StatusUnauthorized, err, true, "Invalid token signature")
					return
				}
				utils.RespondError(w, http.StatusBadRequest, err, true, "Token is expired")
				return
			}
			if !parseToken.Valid {
				utils.RespondError(w, http.StatusUnauthorized, err, true, "Invalid token")
				return
			}

			sessionId := claims.SessionId
			userId := claims.UserId
			sessionEndTime, err := dbHelpers.CheckSession(sessionId, userId)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, err, false, "something went wrong")
				return
			}

			currentTime := time.Now()
			if sessionEndTime.Before(currentTime) {
				utils.RespondError(w, http.StatusBadRequest, errors.New("you are already logged out"), true, "you are already logged out")
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
