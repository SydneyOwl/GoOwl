package repo

import (
	"fmt"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/logger"
)

func StartPullAndWorkflow(repo config.Repo, hook Hook, action string) {
	po := PullOptions{
		Remote: repo.Repoaddr,
		Branch: repo.Branch,
		Type: repo.Type,
	}
	if Checkprotocol(repo) == "ssh" {
		po.Protocol = "ssh"
		po.Sshkey = repo.Sshkeyaddr
	} else {
		po.Protocol = "http"
		if repo.Token != "" {
			po.Token = repo.Token
		} else {
			po.Username = repo.Username
			po.Password = repo.Password
		}
	}
	logger.Info("----------------"+action+"----------------", repo.ID)
	name:=hook.Pusher.Username
	if hook.Pusher.Username==""{
		name=hook.Pusher.Name
	}
	logger.Info(fmt.Sprintf(
		"Pulling Repo:%s(%s),Hash: %s -> %s, %sed by %s......",
		repo.ID,
		GetRepoName(repo),
		hook.Before[0:6],
		hook.After[0:6],
		action,
		name), repo.ID,
	)
	if err := Pull(LocalRepoAddr(repo), po); err != nil {
		repo.BuildStatus = 1
		logger.Warning(
			"Pull error: "+err.Error(), repo.ID)
		return
	}
	//don't throw unrelated exception
	logger.Info("Done Pulling",repo.ID)
	repo.BuildStatus = 2
	logger.Info(fmt.Sprintf(
		"Executing script %s under %s......\n",
		repo.Buildscript,
		LocalRepoAddr(repo),
	), repo.ID)
	standout, err := RunScript(repo)
	if err != nil {
		repo.BuildStatus = 1
		logger.Error(
			"Executing script failed:"+err.Error(), repo.ID,
		)
	}
	logger.Info("Script output:"+standout, repo.ID)
	logger.Info("CICD Done.", repo.ID)
	repo.BuildStatus = 3
}