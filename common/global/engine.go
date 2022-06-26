package global

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine //global engine
)

// routerDetail is the explanation for each route
type routerDetail struct {
	Route       string
	Explanation string
}

// SetEngine set global used engine
func SetEngine(eng *gin.Engine) {
	engine = eng
}

// GetEngine get global used engine
func GetEngine() *gin.Engine {
	return engine
}

// GetAllRouters get routers registered.
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
		if strings.Contains(v.Path, "repoid") {
			routerCurrent.Explanation = ":repoid should be replaced to get repo status"
		}
		detail = append(detail, routerCurrent)
	}
	return detail
}
