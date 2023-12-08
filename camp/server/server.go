package server

import (
	"os"
	"video_service/camp/utils"

	"github.com/gin-gonic/gin"
	"github.com/mufe/golang-base/camp/xlog"
)

var s *server

type server struct {
	r *gin.Engine

	apis []utils.HttpApi
}

// Get 注册api
func Get(path string, funcs ...gin.HandlerFunc) {
	s.r.GET(path, funcs...)
	s.apis = append(s.apis, utils.HttpApi{Pattern: path})
}

// Post 注册api
func Post(path string, funcs ...gin.HandlerFunc) {
	s.r.POST(path, funcs...)
	s.apis = append(s.apis, utils.HttpApi{Pattern: path})
}

// Put 注册api
func Put(path string, funcs ...gin.HandlerFunc) {
	s.r.PUT(path, funcs...)
	s.apis = append(s.apis, utils.HttpApi{Pattern: path})
}

// Delete 注册api
func Delete(path string, funcs ...gin.HandlerFunc) {
	s.r.DELETE(path, funcs...)
	s.apis = append(s.apis, utils.HttpApi{Pattern: path})
}

func init() {
	var mode = os.Getenv(gin.EnvGinMode)
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: xlog.GinLogFormatter,
		Output:    xlog.NewGinLogger(),
	}), gin.Recovery())
	s = &server{r: r}
}

// Start 开启服务
func Start(port string) {
	go utils.KongRegister{}.StartRegister(s.apis, "")
	s.r.Run(":" + port)
}


// Start 开启服务
func StartWithOutRegister(port string) {
	s.r.Run(":" + port)
}
