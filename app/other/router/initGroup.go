package router

import (
	"fmt"

	"github.com/sydneyowl/GoOwl/app/other/apis"
	"github.com/sydneyowl/GoOwl/common/config"

	"github.com/gin-gonic/gin"
)

func initgroup() {
	//hooks only!
	for _, v := range config.WorkspaceConfig.Repo {
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
}