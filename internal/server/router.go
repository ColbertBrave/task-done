package server

import (
	"cloud-disk/app/user"
	"cloud-disk/internal/constants"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	return nil
}

func userAPI(server *Server) {
	server.ginEngine.GET(constants.PrefixURL+"/user-info", user.GetUserInfo)
}
