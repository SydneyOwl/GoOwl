package router

import (
	"fmt"

	"github.com/sydneyowl/GoOwl/app/other/apis"
	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/global"

	"github.com/gin-gonic/gin"
)

// initGroup register hooks to specified route. Should in format "domain/type/id/hook"
func initgroup() {
	//hooks only!
	for _, v := range config.WorkspaceConfig.Repo {
		//reject repo
		if config.CheckInSlice(global.RejectedRepo, v.ID) {
			continue
		}
		route := fmt.Sprintf("/%s/hook", v.ID)
		//GogsRegister
		if v.Type == "gogs" {
			GogsRouterGroup = append(GogsRouterGroup, func(rg *gin.RouterGroup) {
				rg.POST(route, apis.GogsHookReceiver)
			})
		}

		//githubRegister
		if v.Type == "github" {
			GithubRouterGroup = append(GithubRouterGroup, func(rg *gin.RouterGroup) {
				rg.POST(route, apis.GithubHookReceiver)
			})
		}
	}
	routeStats := "/:repoid/status.svg"
	StatusRouterGroup = append(StatusRouterGroup, func(rg *gin.RouterGroup) {
		rg.GET(routeStats, apis.ReportBuildStatus)
	})
}
