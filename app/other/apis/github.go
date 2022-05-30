package apis

import (
	"fmt"
	"strings"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/hook"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"

	"github.com/gin-gonic/gin"
)

// GithubHookReceiver processes hook received in github format and pull the repo/run the script if condition matched.
func GithubHookReceiver(c *gin.Context) {
	fmt.Println("Hook received from github...")
	action := c.GetHeader("X-GitHub-Event")
	hook := hook.GithubHook{
		Pusher: hook.GithubPusher{},
	}
	err := c.ShouldBind(&hook)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": "InternalServerError", //InternalServerErrorErr
		})
		logger.Warning("Err binding struct!","GoOwl-MainLog")
		return
	}
	ref := strings.Split(hook.Ref, "/")
	triggerBranch := ref[len(ref)-1] //branch
	repoID := strings.Split(c.FullPath(), "/")[2]
	// if err!=nil{
	// 	c.JSON(500,gin.H{
	// 		"Status":"InternalServerError",//InternalServerErrorErr
	// 	})
	// 	fmt.Println(stdout.Cyan("Warning: Error converting id to int!"))
	// 	return
	// }
	targetRepo, err := repo.SearchRepo(repoID)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": "InternalServerError", //InternalServerErrorErr
		})
		logger.Warning("No repo found with id " + repoID,"GoOwl-MainLog")
		return
	}
	c.JSON(200, gin.H{
		"Status": "accepted", //InternalServerErrorErr
	})
	//match trigger pull condition.
	if config.CheckInSlice(targetRepo.Trigger, action) && triggerBranch == targetRepo.Branch {
		po := repo.PullOptions{
			Remote: targetRepo.Repoaddr,
			Branch: targetRepo.Branch,
		}
		if repo.Checkprotocol(targetRepo) == "ssh" {
			po.Protocol = "ssh"
			po.Sshkey = targetRepo.Sshkeyaddr
		} else {
			po.Protocol = "http"
			po.Token = targetRepo.Token
		}
		fmt.Println("----------------" + action + "----------------")
		fmt.Printf(
			"Pulling updated Repo:%s(%s),Hash: %s -> %s, %sed by %s......",
			targetRepo.ID,
			repo.GetRepoName(targetRepo),
			hook.Before[0:6],
			hook.After[0:6],
			action,
			hook.Pusher.Name,
		)
		if err := repo.Pull(repo.LocalRepoAddr(targetRepo), po); err != nil {
			c.JSON(500, gin.H{
				"Status": "InternalServerError", //InternalServerErrorErr
			})
				logger.Warning(
					"Warning: Pull error:" + err.Error(),targetRepo.ID,
				)
			return
		}
		fmt.Println("Done")
		logger.Info(fmt.Sprintf(
			"Executing script %s under %s......\n-------------------------\n",
			targetRepo.Buildscript,
			repo.LocalRepoAddr(targetRepo),
		),targetRepo.ID)
		standout, err := repo.RunScript(targetRepo)
		if err != nil {
				logger.Error(
					"-------------------------\nExecuting script failed:" + err.Error(),targetRepo.ID,
				)
		}
		logger.Info("Script output:"+standout,targetRepo.ID)
		logger.Info("-------------------------\nCICD Done.",targetRepo.ID)
		return
	}
	logger.Notice(fmt.Sprintf(
		"Hook received but does not match trigger condition.(%v,%v)\n",
		targetRepo.ID,
		repo.GetRepoName(targetRepo),
	),targetRepo.ID)
}
