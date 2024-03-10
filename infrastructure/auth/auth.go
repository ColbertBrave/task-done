package auth

import (
	"io"
	"net/http"

	"github.com/cloud-disk/app/types/result"

	"github.com/cloud-disk/infrastructure/config"
	"github.com/cloud-disk/infrastructure/log"
)

var Auth HmacAuthenticator

func InitAuth() {
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
