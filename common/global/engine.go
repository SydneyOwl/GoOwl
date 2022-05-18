package global

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine //global engine
)

//explanation for each route
type routerDetail struct {
	Route       string
	Explanation string
}

//Set global used engine
func SetEngine(eng *gin.Engine) {
	engine = eng
}

//get global used engine
func GetEngine() *gin.Engine {
	return engine
}
func GetAllRouters() []routerDetail {
	//show all routes
	routers := engine.Routes()
	var detail []routerDetail
	for _, v := range routers {
		routerCurrent := routerDetail{}
		routerCurrent.Route = v.Path
		routesplit := strings.Split(v.Path, "/")
		if strings.Contains(v.Path, "hook") {
			routerCurrent.Explanation = fmt.Sprintf(
				"Hook for repo %s,type:%s",
				routesplit[2],
				routesplit[1],
			)
		}
		detail = append(detail, routerCurrent)
	}
	return detail
}
