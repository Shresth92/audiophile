package utils

import (
	cloud "cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"github.com/Shresth92/audiophile/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
	"os"
	"strconv"
	"time"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func GetEnvValue(key string) string {
	value := os.Getenv(key)
	return value
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWTToken(userId string, sessionID string, role models.Roles) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		UserId:    userId,
		SessionId: sessionID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtKey := []byte(GetEnvValue("JwtKey"))
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, err
}

func GetFirebaseClient() (*models.App, error) {
	client := &models.App{}
	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(GetEnvValue("FirebaseConfig")))
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		return client, err
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		return client, err
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		return client, err
	}

	return client, nil
}

func GetLimitPage(ctx *gin.Context) (int, int, error) {
	var err error
	limit := 5
	page := 0
	limitQuery := ctx.Query("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		return limit, page, err
	}

	pageQuery := ctx.Query("page")
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		return limit, page, err
	}

	return limit, page, nil
}
