package apis

import (
	"fmt"
	"strings"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/hook"
	"github.com/sydneyowl/GoOwl/common/repo"
	"github.com/sydneyowl/GoOwl/common/stdout"

	"github.com/gin-gonic/gin"
)

//Common resp of hooks
func GogsHookReceiver(c *gin.Context) {
	fmt.Println("Hook received from gogs...")
	action := c.GetHeader("X-Gogs-Event")
	hook := hook.GogsHook{
		Pusher: hook.GogsPusher{},
	}
	err := c.ShouldBind(&hook)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": "InternalServerError", //InternalServerErrorErr
		})
		fmt.Println(stdout.Cyan("Warning: err binding struct!"))
		return
	}
	ref := strings.Split(hook.Ref, "/")
	triggerBranch := ref[len(ref)-1] //branch
	repoID := strings.Split(c.FullPath(), "/")[2]
	targetRepo, err := repo.SearchRepo(repoID)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": "InternalServerError", //InternalServerErrorErr
		})
		fmt.Println(stdout.Cyan("Warning: No repo found with id " + repoID))
		return
	}
	c.JSON(200, gin.H{
		"Status": "Success", //InternalServerErrorErr
	})
	//match trigger pull condition.
	if config.CheckInSlice(targetRepo.Trigger, action) && triggerBranch == targetRepo.Branch {
		po := repo.PullOptions{
			Remote: targetRepo.Repoaddr,
			Branch: targetRepo.Branch,
		}
		if repo.CheckProtocal(targetRepo) == "ssh" {
			po.Protocol = "ssh"
			po.Sshkey = targetRepo.Sshkeyaddr
		} else {
			po.Protocol = "http"
			po.Username = targetRepo.Username
			po.Password = targetRepo.Password
		}
		fmt.Println("----------------" + action + "----------------")
		fmt.Printf(
			"Pulling updated Repo:%s(%s),Hash: %s -> %s, %sed by %s......",
			targetRepo.ID,
			repo.GetRepoName(targetRepo),
			hook.Before[0:6],
			hook.After[0:6],
			action,
			hook.Pusher.Username,
		)
		if err := repo.Pull(repo.LocalRepoAddr(targetRepo), po); err != nil {
			c.JSON(500, gin.H{
				"Status": "InternalServerError", //InternalServerErrorErr
			})
			fmt.Println(
				stdout.Cyan(
					"Warning: Pull error :repo " + repoID + "(" + repo.GetRepoName(
						targetRepo,
					) + ") reports " + err.Error(),
				),
			)
			return
		}
		fmt.Println(stdout.Green("Done"))
		fmt.Printf(
			"Executing script %s under %s......\n-------------------------\n",
			targetRepo.Buildscript,
			repo.LocalRepoAddr(targetRepo),
		)
		standout, err := repo.RunScript(targetRepo)
		if err != nil {
			fmt.Println(
				stdout.Cyan(
					"-------------------------\nWarning: Executing script failed:" + err.Error(),
				),
			)
			return
		}
		fmt.Println(standout)
		fmt.Println("-------------------------\nCICD Done.")
		return
	}
	fmt.Printf(
		"Hook received but does not match trigger condition.(%s,%s)\n",
		targetRepo.ID,
		repo.GetRepoName(targetRepo),
	)
}
