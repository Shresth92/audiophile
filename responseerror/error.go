package responseerror

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	ServerErrorMsg = "Internal Server Error occurred when processing request."
)

type clientError struct {
	ID            string        `json:"id"`
	MessageToUser string        `json:"messageToUser"`
	DeveloperInfo string        `json:"developerInfo"`
	Err           string        `json:"error"`
	StatusCode    int           `json:"statusCode"`
	IsClientError bool          `json:"isClientError"`
	Request       *http.Request `json:"-"`
}

func newClientError(err error, statusCode int, req *http.Request, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}

	var errString string
	if err != nil {
		errString = err.Error()
	}

	return &clientError{
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
		Request:       req,
	}
}

func RespondClientErr(ctx *gin.Context, err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("messageToUser: %v with error %v", messageToUser, err)
	clientError := newClientError(err, statusCode, ctx.Request, messageToUser, additionalInfoForDevs...)
	ctx.JSON(statusCode, clientError)
}

func RespondGenericServerErr(ctx *gin.Context, err error, additionalInfoForDevs ...string) {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	serverErr := &clientError{
		MessageToUser: ServerErrorMsg,
		DeveloperInfo: additionalInfoJoined,
		Err:           err.Error(),
		StatusCode:    http.StatusInternalServerError,
		IsClientError: false,
		Request:       ctx.Request,
	}
	ctx.JSON(http.StatusInternalServerError, serverErr)
}
