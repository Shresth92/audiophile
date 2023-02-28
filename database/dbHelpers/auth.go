package dbHelpers

import (
	"errors"
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var (
	session models.Session
)

func CheckSession(sessionId uuid.UUID, userId uuid.UUID) (time.Time, error) {
	err := database.Db.
		Model(&models.Session{}).
		Where("id = ? AND user_id = ?", sessionId, userId).
		Find(&session).
		Error
	return session.EndedAt, err
}

func CheckUserExist(email string) (bool, error) {
	user := models.Users{}
	err := database.Db.
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

func GetUserDetails(email string) (models.Users, error) {
	user := models.Users{}
	err := database.Db.
		Model(&models.Users{}).
		Where("email = ? AND archived_at is null", email).
		First(&user).
		Error
	return user, err
}

func CreateUser(tx *gorm.DB, email string, password string) (uuid.UUID, error) {
	userId := uuid.New()
	user := models.Users{
		Id:       userId,
		Email:    email,
		Password: password,
	}
	err := tx.
		Model(&models.Users{}).
		Create(&user).
		Error
	return user.Id, err
}

func CreateUserRole(tx *gorm.DB, userId uuid.UUID, role models.Roles) (uuid.UUID, error) {
	roleId := uuid.New()
	roleInstance := models.UserRole{
		Id:     roleId,
		UserId: userId,
		Role:   role,
	}
	err := tx.
		Model(&models.UserRole{}).
		Create(&roleInstance).
		Error
	return roleInstance.Id, err
}

func GetSessionId(userId uuid.UUID) (uuid.UUID, error) {
	sessionID := uuid.New()
	session := models.Session{
		ID:      sessionID,
		UserId:  userId,
		EndedAt: time.Now().Add(60 * time.Minute),
	}
	err := database.Db.
		Model(&models.Session{}).
		Create(&session).
		Error
	return session.ID, err
}

func LogOut(userID uuid.UUID, sessionId uuid.UUID) error {
	err := database.Db.
		Model(&models.Session{}).
		Where("user_id = ? AND id = ?", userID, sessionId).
		Update("end_at", time.Now()).
		Error
	return err
}

func AddAddress(userId uuid.UUID, newAddress *models.Address) (uuid.UUID, error) {
	addressId := uuid.New()
	address := models.Address{
		Id:      addressId,
		UserId:  userId,
		Area:    newAddress.Area,
		City:    newAddress.City,
		State:   newAddress.State,
		ZipCode: newAddress.ZipCode,
		Contact: newAddress.Contact,
		LatLong: newAddress.LatLong,
	}
	err := database.Db.
		Model(&models.Address{}).
		Create(&address).
		Error
	return addressId, err
}
