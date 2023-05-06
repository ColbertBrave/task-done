package server

import (
	"github.com/cloud-disk/app/user"
	"github.com/cloud-disk/internal/constants"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	return nil
}

func userAPI(server *Server) {
	server.ginEngine.GET(constants.PrefixURL+"/user-info", user.GetUserInfo)
}
