package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type (
	Claims struct {
		UserId    string `json:"userId"`
		SessionId string `json:"sessionId"`
		Role      Roles  `json:"role"`
		jwt.RegisteredClaims
	}
)
