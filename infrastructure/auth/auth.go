package auth

import (
	"io"
	"net/http"

	"github.com/task-done/app/types/result"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/log"
)

var Auth HmacAuthenticator

func Init() {
	Auth.SecretKey = []byte(config.GetConfig().Auth.SecretKey)
}

func VerifyRequest(authenticator Authenticator, request *http.Request) error {
	sign, isOk := request.Header["Authorization"]
	if !isOk || len(sign) == 0 {
		return result.ErrNoAuthorization
	}

	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error("read the request body error:%s", err)
		return err
	}

	return authenticator.Verify(string(bytes), sign[0])
}
