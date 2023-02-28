package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type (
	Claims struct {
		UserId    uuid.UUID `json:"userId"`
		SessionId uuid.UUID `json:"sessionId"`
		Role      Roles     `json:"role"`
		jwt.RegisteredClaims
	}
)
