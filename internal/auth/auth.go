package auth

import (
	"errors"
	"io"
	"net/http"

	"cloud-disk/internal/log"
)

var Auth HmacAuthenticator

func VerifyRequest(authenticator Authenticator, request *http.Request) error {
	sign, isOk := request.Header["Authorization"]
	if !isOk || len(sign) == 0 {
		return errors.New("no Authorization field in the request header")
	}

	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error("read the request body error:%s", err)
		return err
	}

	return authenticator.Verify(string(bytes), sign[0])
}
