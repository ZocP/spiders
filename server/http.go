package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"qa_spider/config"
	"qa_spider/server/content"
)

type HTTPServer struct {
	config *config.Config
	log    *zap.Logger
	engine *gin.Engine
	ctn    map[interface{}]*content.Content
}

func (s *HTTPServer) Run() error {
	if err := s.engine.Run(s.config.Server.Port); err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) Stop() {
	s.log.Sync()
}

func InitHTTPServer(config *config.Config, logger *zap.Logger) Server {
	s := &HTTPServer{
		config: config,
		log:    logger,
		engine: gin.Default(),
		ctn:    make(map[interface{}]*content.Content),
	}
	//init content services
	s.initContent()

	//set mode
	gin.SetMode(gin.DebugMode)

	//init internal dependencies

	//init handlers
	s.regHandlers()

	//allow cors
	if config.Server.AllowCors {
		logger.Info("Server allow cors enabled")
		s.engine.Use(Cors())
	} else {
		logger.Info("Server allow cors disabled")
	}

	return s
}

func (s *HTTPServer) initContent() {
}

//router initialize
func (s *HTTPServer) regHandlers() {
	s.ctn[config.QueryQA] =
}

//Cors management

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Max-Age", "600")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
