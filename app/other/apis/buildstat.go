package apis

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/narqo/go-badge"
	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/database"
	"github.com/sydneyowl/GoOwl/common/global"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"
	"gorm.io/gorm"
)

//ReportBuildStatus returns svg reporting build status.
func ReportBuildStatus(c *gin.Context) {
	repoid := c.Param("repoid")
	targetRepo, err := repo.SearchRepo(repoid)
	var badgeData []byte
	if err != nil {
		logger.Notice("No repo found with this ID!", "GoOwl-MainLog")
		//returnnofoundhere
		badgeData, err = badge.RenderBytes("Repo", "Not Found", badge.ColorLightgray)
	} else {
		buildstat := &config.BuildInfo{}
		if err := database.GetConn().Select("build_status").Where("repo_id=?", repoid).Order("created_at desc").First(buildstat).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				buildstat.BuildStatus = 4
			} else {
				buildstat.BuildStatus = 5
			}
		}
		switch buildstat.BuildStatus {
		case 0, 4:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Unknown", badge.ColorYellowgreen)
		case 1:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Failed", badge.ColorOrange)
		case 2:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Waiting", badge.ColorYellow)
		case 3:
			badgeData, err = badge.RenderBytes(repo.GetRepoOriginalName(targetRepo)+" Build", "Passed", badge.ColorGreen)
		case 5:
			logger.Warning("No repo found in current database!", "GoOwl-MainLog")
			badgeData, err = badge.RenderBytes("GoOWL", "ERROR", badge.ColorRed)
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
