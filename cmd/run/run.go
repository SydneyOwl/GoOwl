package run

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/database"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/global"
	"github.com/sydneyowl/GoOwl/common/logger"
	"github.com/sydneyowl/GoOwl/common/repo"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	LoggingMethod int
	yamlAddr      string
	AppRouters    = make([]func(), 0) // Storages routers
	skipRepoCheck bool
	StartCmd      = &cobra.Command{
		Use:     "run",
		Short:   "Run GoOwl",
		Example: "GoOwl run -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

// Specify yaml before run
func init() {
	addrConfig := file.GetCwd() + "/config/settings.yaml"
	StartCmd.Flags().
		StringVarP(&yamlAddr, "config", "c", addrConfig, "Run GoOwl using specified yaml config. Use $PWD/config/settings.yaml if not specified.")
	StartCmd.Flags().
		BoolVar(&skipRepoCheck, "skip-repocheck", false, "Skip check of repo config, including address and authorization.")
	StartCmd.Flags().
		IntVarP(&LoggingMethod, "log", "l", 1, "Specify logging method. 1->stdout 2->file 3->both.")
	StartCmd.Flags().
		StringVarP(&global.Sqlite3DBPosition, "database-location", "d", "./GoOwl.db", "Specify the position database used by GoOwl storage in. Default is ./GoOwl.db.")
	StartCmd.Flags().BoolVar(&global.SqlDebug, "enable-sqldebug", false, "Print all sql sentences")
}

// initGin starts gin framework.
func initGin() *gin.Engine {
	r := gin.Default()
	global.SetEngine(r)
	//	common.InitMiddleware(r)
	return r
}

// initCloneRepo clones repo on not exist
func initCloneRepo() bool {
	var exists bool
	//Clone repo unexists
	for _, v := range config.WorkspaceConfig.Repo {
		if repo.Checkprotocol(v) == "ssh" {
			logger.Notice("Manually answer yes if required.", "GoOwl-MainLog")
		}
		if err := repo.CloneOnNotExist(v); err != nil {
			global.RejectedRepo = append(global.RejectedRepo, v.ID)
			exists = true
			logger.Error("Error:"+err.Error(), v.ID)
		}
	}
	return exists
}

// initGinEngine init gin engine.
func initGinEngine() (engine *gin.Engine, suspend bool) {
	//set to release mode
	if config.ApplicationConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine = initGin()
	//init routes
	for _, f := range AppRouters {
		f() //init all routers!
	}
	//Welcome Page
	routers := global.GetAllRouters()
	fmt.Println("------------------------------------------------------")
	if len(routers) != 0 {
		fmt.Println("Here're all routes you've registered:")
		//rej
		for _, v := range global.RejectedRepo {
			fmt.Printf("[rejected] Repo %v(%v)\n", v, "failed to clone")
		}
	} else {
		fmt.Println("No route is registered. GoOwl suspend")
		suspend = true
		//rej
		for _, v := range global.RejectedRepo {
			fmt.Printf("[rejected] Repo %v(%v)\n", v, "failed to clone")
		}
		return
	}
	for _, v := range routers {
		fmt.Println(v.Route + "---------------->" + v.Explanation)
	}
	return
}

// run Runs main application.
func run() {
	global.LoggingMethod = LoggingMethod
	// fmt.Println(global.LoggingMethod)
	if global.LoggingMethod != 1 && global.LoggingMethod != 2 && global.LoggingMethod != 3 {
		global.LoggingMethod = 1 //reset
	}
	if global.LoggingMethod != 1 {
		err := file.CreateDir(config.WorkspaceConfig.Path + "/log")
		if err != nil {
			fmt.Println("Error creating log dir!")
			return
		}
		for _, v := range config.WorkspaceConfig.Repo {
			name := v.ID
			curDate := time.Now().Format("2006-01-02 15:04:05")
			err := file.CreateFile(fmt.Sprintf("[%s]Repo_%s", curDate, name))
			if err != nil {
				fmt.Println("Error creating log file!")
				return
			}
		}
	}
	config.InitConfig(&yamlAddr)
	//Check repo
	if !skipRepoCheck {
		repo.CheckRepo()
	} else {
		logger.Notice("Check skipped.", "GoOwl-MainLog")
	}
	//init all repo
	var iserr bool = initCloneRepo()
	if iserr {
		logger.Error(
			"Err occurred. Check and fix it if necessary. Those routes of repos that failed to clone will not be registered.",
			"GoOwl-MainLog",
		)
	}
	if engine, suspend := initGinEngine(); suspend {
		return
	} else {
		//goroutine here
		go func() {
			err := engine.Run(
				fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
			)
			if err != nil {
				logger.Fatal(
					"Cannot start Gin framework:"+err.Error(),
					"GoOwl-MainLog",
				)
			}
		}()
	}
	if err := database.InitDB(); err != nil {
		fmt.Println(err)
	}
	if global.LoggingMethod == 2 {
		fmt.Println("Log write to file only.")
	}
	//In order to use ^c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) //STUFF UNTIL CHAN IN
	<-quit
	fmt.Println("\nGoOwl Exit.")
}
