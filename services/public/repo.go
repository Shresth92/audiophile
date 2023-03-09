package public

import (
	"errors"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type repository struct {
	*internal.Database
}

func newPublicRepository(db *internal.Database) *repository {
	return &repository{Database: db}
}

func (r *repository) checkSession(sessionId string, userId string) (time.Time, error) {
	session := models.Session{}
	err := r.Database.DB.
		Model(&models.Session{}).
		Where("id = ? AND user_id = ?", sessionId, userId).
		Find(&session).
		Error
	return session.EndedAt, err
}

func (r *repository) checkUserExist(email string) (bool, error) {
	user := models.Users{}
	err := r.Database.DB.
		Model(&models.Users{}).
		Where("email = ? AND archived_at is null", email).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (r *repository) getUserDetails(email string) (models.Users, error) {
	user := models.Users{}
	err := r.Database.DB.
		Model(&models.Users{}).
		Where("email = ? AND archived_at is null", email).
		First(&user).
		Error
	return user, err
}

func (r *repository) createUser(email string, password string) (string, error) {
	userId := uuid.New().String()
	user := models.Users{
		Id:       userId,
		Email:    email,
		Password: password,
	}
	err := r.Database.DB.
		Model(&models.Users{}).
		Create(&user).
		Error
	return user.Id, err
}

func (r *repository) createUserRole(userId string, role models.Roles) (string, error) {
	roleId := uuid.New().String()
	roleInstance := models.UserRole{
		Id:     roleId,
		UserId: userId,
		Role:   role,
	}
	err := r.Database.DB.
		Model(&models.UserRole{}).
		Create(&roleInstance).
		Error
	return roleInstance.Id, err
}

func (r *repository) getSessionId(userId string) (string, error) {
	sessionID := uuid.New().String()
	session := models.Session{
		ID:      sessionID,
		UserId:  userId,
		EndedAt: time.Now().Add(60 * time.Minute),
	}
	err := r.Database.DB.
		Model(&models.Session{}).
		Create(&session).
		Error
	return session.ID, err
}

func (r *repository) logout(userID string, sessionId string) error {
	err := r.Database.DB.
		Model(&models.Session{}).
		Where("user_id = ? AND id = ?", userID, sessionId).
		Update("end_at", time.Now()).
		Error
	return err
}
