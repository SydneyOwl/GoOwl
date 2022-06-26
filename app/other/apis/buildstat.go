package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/narqo/go-badge"
	"github.com/sydneyowl/GoOwl/common/global"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"
)

func ReportBuildStatus(c *gin.Context) {
	repoid := c.Param("repoid")
	targetRepo, err := repo.SearchRepo(repoid)
	var badgeData []byte
	if err != nil {
		logger.Notice("No repo found with this ID!", "GoOwl-MainLog")
		//returnnofoundhere
		badgeData, err = badge.RenderBytes("Repo", "Not Found", badge.ColorLightgray)
	} else {
		switch targetRepo.BuildStatus {
		case 0:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Unknown", badge.ColorYellowgreen)
		case 1:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Failed", badge.ColorOrange)
		case 2:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Waiting", badge.ColorYellow)
		case 3:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Passed", badge.ColorGreen)
		}
	}
	if err != nil {
		logger.Warning("Failed to generate badge", "GoOwl-MainLog")
		c.Writer.WriteString(global.ErrorBuildStatus)
		return
	}
	c.Writer.Header().Set("content-type", "image/svg+xml")
	c.Writer.WriteString(string(badgeData))
}
