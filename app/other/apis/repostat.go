package apis

import (
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"
)
type repoInfo struct{
	RepoName string
	RepoID string
	RepoAddr string
	RepoLocalAddr string
	RepoSize float32
}
type triggerInfo struct{
	LastAvailableAction string
	LastAction string
	AllowedAction string
	LAATriggerBy string
	LATriggeredBy string
	HashBeforeTrigger string
	HashAfterTrigger string
}
type buildInfo struct{
	BuildStatus string
	Output string
	SvgAddr string
	ScriptAddr string
}
type repoStat struct{
	RepoInfo repoInfo
	TriggerInfo triggerInfo
	BuildInfo buildInfo
}
// repoStatget fills info of repostat but is still unable in this version.
func repoStatGet(repoid string, rs *repoStat){
	rs.RepoInfo.RepoID=repoid
	if info,err:=repo.SearchRepo(repoid);err!=nil{
		rs.RepoInfo.RepoName="Repo does not exists"
		logger.Info("Repo "+repoid+" not found!","Goowl-MainLog")
		return 
	}else{
		rs.RepoInfo.RepoName=repo.GetRepoName(info)
	}
}