package apis

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/database"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"
	"gorm.io/gorm"
)

type repoInfo struct {
	RepoName      string
	RepoID        string
	RepoAddr      string
	RepoLocalAddr string
	RepoSize      int64
}
type triggerInfo struct {
	LastAvailableAction  string
	LastAction           string
	AllowedAction        string
	LAATriggerBy         string
	LATriggeredBy        string
	LAAHashBeforeTrigger string
	LAAHashAfterTrigger  string
}
type buildInfo struct {
	BuildStatus string
	Output      string
	SvgAddr     string
	ScriptAddr  string
	TimeCost    int64
}
type repoStat struct {
	RepoInfo    repoInfo
	TriggerInfo triggerInfo
	BuildInfo   buildInfo
}

// repoStatget fills info of repostat but is still unable in this version.
func repoStatGet(repoid string, rs *repoStat) {
	//repo
	rs.RepoInfo.RepoID = repoid
	info, err := repo.SearchRepo(repoid)
	if err != nil {
		rs.RepoInfo.RepoName = "Repo does not exists"
		logger.Info("Repo "+repoid+" not found!", "Goowl-MainLog")
		return
	} else {
		rs.RepoInfo.RepoName = repo.GetRepoName(info)
		rs.RepoInfo.RepoAddr = info.Repoaddr
		rs.RepoInfo.RepoLocalAddr = repo.LocalRepoAddr(info)
		if repoSize, err := file.CalcSize(repo.LocalRepoAddr(info)); err != nil {
			logger.Error("Failed to calculate size of repo "+repoid+" !", repoid)
			rs.RepoInfo.RepoSize = 0
		} else {
			rs.RepoInfo.RepoSize = repoSize
		}

	}
	//trigger
	availableActon := &config.TriggerInfo{}
	if err := database.GetConn().Where("is_available_action=? and repo_id=?", 1, repoid).Order("created_at desc").Find(availableActon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Notice("Trigger of repo not found!", repoid)
			rs.TriggerInfo.LastAvailableAction = "No Record Found"
		} else {
			logger.Error("error to to get db trig info!", "GoOwl-MainLog")
			rs.TriggerInfo.LastAvailableAction = "Unknown"
		}
	} else {
		rs.TriggerInfo.LastAvailableAction = availableActon.Action
		rs.TriggerInfo.LAATriggerBy = availableActon.TriggerBy
		rs.TriggerInfo.LAAHashBeforeTrigger = availableActon.HashBeforeTrigger
		rs.TriggerInfo.LAAHashAfterTrigger = availableActon.HashAfterTrigger
	}
	action := &config.TriggerInfo{}
	if err := database.GetConn().Where("repo_id=?", repoid).Order("created_at desc").Find(action).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Notice("Trigger of repo not found!", repoid)
			rs.TriggerInfo.LastAction = "No Record Found"
		} else {
			logger.Error("error to to get db trig info!", "GoOwl-MainLog")
			rs.TriggerInfo.LastAction = "Unknown"
		}
	} else {
		rs.TriggerInfo.LATriggeredBy = action.TriggerBy
		rs.TriggerInfo.LastAction = action.Action
	}
	rs.TriggerInfo.AllowedAction = func() string {
		var trigger string
		for _, v := range info.Trigger {
			trigger += v + ","
		}
		return trigger[0 : len(trigger)-1]
	}()
	//buildinfo
	build := &config.BuildInfo{}
	if err := database.GetConn().Where("repo_id=?", repoid).Order("created_at desc").Find(build).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Notice("buildstatus of repo not found!", repoid)
			rs.BuildInfo.BuildStatus = "No Record Found"
		} else {
			logger.Error("error to to get db buildss info!", "GoOwl-MainLog")
			rs.BuildInfo.BuildStatus = "Unknown"
		}
	} else {
		switch build.BuildStatus {
		case 0, 4:
			rs.BuildInfo.BuildStatus = "Unknown"
		case 1:
			rs.BuildInfo.BuildStatus = "Failed"
		case 3:
			rs.BuildInfo.BuildStatus = "Passed"
		}
		rs.BuildInfo.Output = build.Output
		rs.BuildInfo.TimeCost = build.TimeCost
		rs.BuildInfo.SvgAddr = "http(s)://domain:port/status/" + repoid + "status.svg"
		rs.BuildInfo.ScriptAddr = info.Buildscript
	}
}
func RepoStat(c *gin.Context) {
	repoid := c.Param("repoid")
	repoStatus := &repoStat{}
	repoStatGet(repoid, repoStatus)
	c.JSON(200, gin.H{
		"info": repoStatus,
	})
}
