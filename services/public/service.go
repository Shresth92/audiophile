package public

import (
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"time"
)

type Service struct {
	repo *repository
}

func NewPublicService(db *internal.Database) *Service {
	return &Service{repo: newPublicRepository(db)}
}

func (s *Service) CheckSession(sessionID string, userID string) (time.Time, error) {
	return s.repo.checkSession(sessionID, userID)
}

func (s *Service) CheckUserExist(email string) (bool, error) {
	return s.repo.checkUserExist(email)
}

func (s *Service) GetUserDetails(email string) (models.Users, error) {
	return s.repo.getUserDetails(email)
}

func (s *Service) CreateUser(email string, password string) (string, error) {
	return s.repo.createUser(email, password)
}

func (s *Service) CreateUserRole(userId string, role models.Roles) (string, error) {
	return s.repo.createUserRole(userId, role)
}

func (s *Service) GetSessionId(userID string) (string, error) {
	return s.repo.getSessionId(userID)
}

func (s *Service) Logout(userID string, sessionID string) error {
	return s.repo.logout(userID, sessionID)
}
