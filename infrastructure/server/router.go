package server

import (
	"github.com/task-done/app/user"
	"github.com/task-done/infrastructure/constants"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	return nil
}

func userAPI(server *Server) {
	server.ginEngine.GET(constants.PrefixURL+"/user-info", user.GetUserInfo)
}
