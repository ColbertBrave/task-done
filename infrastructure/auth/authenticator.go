package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cloud-disk/infrastructure/log"
)

type Authenticator interface {
	Sign(body string, expireTime int64) (string, error)
	Verify(body string, sign string) error
}

type HmacAuthenticator struct {
	SecretKey []byte
}

func (h HmacAuthenticator) Sign(body string, expireTime int64) (string, error) {
	hmacHash := hmac.New(sha256.New, h.SecretKey)
	expireTimeStamp := strconv.FormatInt(expireTime, 10)
	_, err := io.WriteString(hmacHash, body+":"+expireTimeStamp)
	if err != nil {
		log.Error("fail to call io.WriteString:%s", err)
		return "", err
	}

	sign := base64.URLEncoding.EncodeToString(hmacHash.Sum(nil)) + ":" + expireTimeStamp
	return sign, nil
}

func (h HmacAuthenticator) Verify(body string, sign string) error {
	signSlice := strings.Split(sign, ":")
	if signSlice[len(signSlice)-1] == "" {
		return errors.New("empty expire time field in the sign")
	}

	// 签名信息的最后一位为过期时间
	expireTime, err := strconv.ParseInt(signSlice[len(signSlice)-1], 10, 64)
	if err != nil {
		return err
	}

	if expireTime < time.Now().Unix() && expireTime >= 0 {
		return errors.New("sign is expired")
	}

	verifySign, err := h.Sign(body, expireTime)
	if err != nil {
		return err
	}

	if verifySign != sign {
		return errors.New("sign is not correct")
	}
	return nil
}
