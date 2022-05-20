package router

import (
	"github.com/sydneyowl/GoOwl/common/global"

	"github.com/gin-gonic/gin"
)

//mulitVersion
var (
	GogsRouterGroup   = make([]func(*gin.RouterGroup), 0)
	GithubRouterGroup = make([]func(*gin.RouterGroup), 0)
)

// GogsRouter defines group router name of hook of gogs.
func GogsRouter(eng *gin.Engine) {
	v1 := eng.Group("/gogs")
	//add more routergroup
	for _, f := range GogsRouterGroup {
		f(v1) //give them same address
	}
}

// GithubRouter defines group router name of hook of github.
func GithubRouter(eng *gin.Engine) {
	v1 := eng.Group("/github")
	//add more routergroup
	for _, f := range GithubRouterGroup {
		f(v1) //give them same address
	}
}

// initAllRouter simply init all routers
func InitAllRouter() {
	initgroup()
	engine := global.GetEngine()
	GogsRouter(engine) //engine gogs
	GithubRouter(engine)
}
