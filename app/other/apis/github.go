package apis

import (
	"fmt"
	"strings"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/database"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"

	"github.com/gin-gonic/gin"
)

// GithubHookReceiver processes hook received in github format and pull the repo/run the script if condition matched.
func GithubHookReceiver(c *gin.Context) {
	fmt.Println("Hook received from github...")
	action := c.GetHeader("X-GitHub-Event")
	hook := repo.Hook{
		Pusher: repo.Pusher{},
	}
	err := c.ShouldBind(&hook)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": "InternalServerError", //InternalServerErrorErr
		})
		logger.Warning("Err binding struct!", "GoOwl-MainLog")
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
		logger.Warning("No repo found with id "+repoID, "GoOwl-MainLog")
		return
	}
	c.JSON(200, gin.H{
		"Status": "accepted", //InternalServerErrorErr
	})
	trigger := &config.TriggerInfo{
		RepoID:            repoID,
		Action:            action,
		Branch:            targetRepo.Branch,
		HookType:          "Github",
		TriggerBy:         hook.Pusher.Name,
		HashBeforeTrigger: hook.Before,
		HashAfterTrigger:  hook.After,
	}
	//match trigger pull condition.
	if config.CheckInSlice(targetRepo.Trigger, action) && triggerBranch == targetRepo.Branch {
		trigger.IsAvailableAction = 1
		repo.StartPullAndWorkflow(targetRepo, hook, action)
	} else {
		trigger.IsAvailableAction = 0
		logger.Notice(fmt.Sprintf(
			"Hook received but does not match trigger condition.(%v,%v)\n",
			targetRepo.ID,
			repo.GetRepoName(targetRepo),
		), targetRepo.ID)
		trigger.IsAvailableAction = 0
	}
	if err := database.GetConn().Create(trigger).Error; err != nil {
		logger.Warning("Failed to write trigger data to db.", "GoOwl-MainLog")
	}
}
