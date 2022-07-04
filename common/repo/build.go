package repo

import (
	"fmt"
	"strconv"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/database"
	"github.com/sydneyowl/GoOwl/common/logger"
)

func StartPullAndWorkflow(repo config.Repo, hook Hook, action string) {
	po := PullOptions{
		Remote: repo.Repoaddr,
		Branch: repo.Branch,
		Type:   repo.Type,
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
	name := hook.Pusher.Username
	if hook.Pusher.Username == "" {
		name = hook.Pusher.Name
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
		SetBuildStat(repo.ID, 1)
		logger.Warning(
			"Pull error: "+err.Error(), repo.ID)
			database.GetConn().Create(&config.BuildInfo{
				RepoID:      repo.ID,
				BuildStatus: 1,
				Output:      "Failed to pull repo.",
				TimeCost:    0,
			})
		return
	}
	//don't throw unrelated exception
	logger.Info("Done Pulling", repo.ID)
	SetBuildStat(repo.ID, 2)
	logger.Info(fmt.Sprintf(
		"Executing script %s under %s......\n",
		repo.Buildscript,
		LocalRepoAddr(repo),
	), repo.ID)
	cost, standout, err := RunScript(repo)
	if err != nil {
		SetBuildStat(repo.ID, 1)
		logger.Error(
			"Executing script failed:"+err.Error()+"\n time cost:"+strconv.Itoa(int(cost))+"ms", repo.ID,
		)

			database.GetConn().Create(&config.BuildInfo{
				RepoID:      repo.ID,
				BuildStatus: 1,
				Output:      err.Error(),
				TimeCost:    cost,
			})
		return
	} else {
		logger.Info("Script output:"+standout+"\n time cost:"+strconv.Itoa(int(cost))+"ms", repo.ID)
	}
	logger.Info("CICD Done.", repo.ID)
	SetBuildStat(repo.ID, 3)
	if err:=database.GetConn().Create(&config.BuildInfo{
		RepoID:      repo.ID,
		BuildStatus: 3,
		Output:      standout,
		TimeCost:    cost,
	}).Error;err!=nil{
		logger.Warning("Failed to write build data to db","GoOwl-MainLog")
	}
}
