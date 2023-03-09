package services

import (
	"github.com/Shresth92/audiophile/models"
	"time"
)

type PublicService interface {
	CheckSession(sessionID string, userID string) (time.Time, error)
	CheckUserExist(email string) (bool, error)
	GetUserDetails(email string) (models.Users, error)
	CreateUser(email string, password string) (string, error)
	CreateUserRole(userID string, role models.Roles) (string, error)
	GetSessionId(userID string) (string, error)
	Logout(userID string, sessionID string) error
}
