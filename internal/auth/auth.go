package auth

// func VerifySignOfRequest(authenticator Authenticator) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var err error
// 		switch c.Request.Method {
// 		case "PATCH", "POST",	"PUT":

// 	}
// }

// func verifyRequest(authenticator Authenticator, request *http.Request) error {
// 	sign, isOk := request.Header["Authorization"]
// 	if !isOk || len(sign) == 0 {
// 		return errors.New("no Authorization field in the request header")
// 	}

// 	body := []byte{}
// 	if

// 	return authenticator.Verify(request.Body, sign[0])
// }
