package utils

import (
	cloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"github.com/Shresth92/audiophile/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type clientError struct {
	ID            uuid.UUID `json:"id"`
	MessageToUser string    `json:"messageToUser"`
	DeveloperInfo string    `json:"developerInfo"`
	Err           string    `json:"error"`
	StatusCode    int       `json:"statusCode"`
	IsClientError bool      `json:"isClientError"`
}

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

func DecodeBody(body io.Reader, bodyObj interface{}) error {
	if err := json.NewDecoder(body).Decode(bodyObj); err != nil {
		return err
	}
	return nil
}

func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		RespondError(w, statusCode, err, false, "Something went wrong")
		return
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error, clientErr bool, additionalInfo ...string) {
	w.WriteHeader(statusCode)
	infoJoined := strings.Join(additionalInfo, "\n")
	if infoJoined == "" {
		infoJoined = err.Error()
	}
	errId := uuid.New()
	errMessage := &clientError{
		ID:            errId,
		MessageToUser: infoJoined,
		DeveloperInfo: infoJoined,
		Err:           err.Error(),
		StatusCode:    statusCode,
		IsClientError: clientErr,
	}
	Respond(w, statusCode, errMessage)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWTToken(userId uuid.UUID, sessionID uuid.UUID, role models.Roles) (string, error) {
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

func GetLimitPage(urlValues url.Values) (int, int, error) {
	var err error
	limit := 5
	page := 0
	limitQuery := urlValues.Get("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		return limit, page, err
	}

	pageQuery := urlValues.Get("page")
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		return limit, page, err
	}

	return limit, page, nil
}
