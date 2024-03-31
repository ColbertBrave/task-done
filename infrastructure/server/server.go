package server

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/task-done/infrastructure/auth"
	"github.com/task-done/infrastructure/config"
	"github.com/task-done/infrastructure/log"
)

var server *Server

type Server struct {
	httpServer *http.Server
	ginEngine  *gin.Engine
}

type Option func(engine *gin.Engine)

func Init() {
	serverAddr := config.GetConfig().Server.Host + ":" + config.GetConfig().MySQL.Port
	if serverAddr == "" {
		log.Error("the server addr is empty")
		return
	}

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(GetCostTimeOfRequest())
	ginEngine.Use(Authenticate())

	server = &Server{
		httpServer: &http.Server{
			Addr:    serverAddr,
			Handler: ginEngine,
		},
		ginEngine: ginEngine,
	}
}

func Run() {
	if err := server.Start(); err != nil {
		log.Error("run the http server error|%s", err)
		return
	}
	log.Info("successfully run the http server")
}

func GetCostTimeOfRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		costTime := time.Since(startTime)
		log.Info("%s|%s|cost time %d ms", c.Request.Method, c.Request.URL, costTime)
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.VerifyRequest(auth.Auth, c.Request)
		if err != nil {
			log.Error("verify request error:%s", err)
			return
		}
		c.Next()
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	err = s.httpServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func GET(relativePath string, handlers ...gin.HandlerFunc) {
	if server == nil {
		Init()
	}

	server.ginEngine.GET(relativePath, handlers...)
}

func POST(relativePath string, handlers ...gin.HandlerFunc) {
	if server == nil {
		Init()
	}

	server.ginEngine.POST(relativePath, handlers...)
}

func Close() {
	if server == nil {
		log.Info("the http server is nil")
		return
	}

	if err := server.httpServer.Close(); err != nil {
		log.Error("close http server err|%s", err)
		return
	}
	log.Info("the http server is closed")
}
