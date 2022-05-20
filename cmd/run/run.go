package run

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/global"
	"github.com/sydneyowl/GoOwl/common/repo"
	"github.com/sydneyowl/GoOwl/common/stdout"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	yamlAddr   string
	AppRouters = make([]func(), 0) // Storages routers
	StartCmd   = &cobra.Command{
		Use:     "run",
		Short:   "Run GoOwl as backend",
		Example: "GoOwl run -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

//Specify yaml before run
func init() {
	StartCmd.Flags().
		StringVarP(&yamlAddr, "run", "c", "", "Run GoOwl using specified yaml config. Use $PWD/config/settings.yaml if not specified.")
}

// run Run main application.
func run() {
	if readable, err := file.CheckYamlReadable(&yamlAddr); !readable {
		fmt.Println(stdout.Magenta("FATAL:" + err.Error()))
		return
	}
	rawConfig, err := config.LoadConfigFromYaml(yamlAddr) //returns raw viper obj
	if err := config.CheckViperErr(err); err != nil {
		fmt.Println(stdout.Magenta(err.Error()))
		return
	}
	if err := rawConfig.Unmarshal(config.YamlConfig); err != nil {
		fmt.Println(stdout.Magenta("Unknown Error occurred!"))
		return
	}
	//Check repo
	if repeated, err := repo.IsDuplcatedRepo(config.WorkspaceConfig.Repo); repeated {
		fmt.Println(err.Error())
		return
	}
	ID, uncritialerror, err := repo.CheckRepoConfig(config.WorkspaceConfig.Repo)
	if err != nil {
		fmt.Println(stdout.Magenta("Repo " + ID + " has an invaild config:" + err.Error()))
		return
	}
	if len(uncritialerror) > 0 {
		for _, v := range uncritialerror {
			fmt.Println(
				stdout.Magenta(
					"repo " + v.ID + " has an invaild config:" + v.Uerror.Error() + ",check if it is correct.",
				),
			)
		}
	}
	//Clone repo unexists
	fmt.Println(stdout.Yellow("Manually answer yes if required."))
	for _, v := range config.WorkspaceConfig.Repo {
		if err := repo.CloneOnNotExist(v); err != nil {
			fmt.Println(stdout.Cyan("Error:" + err.Error()))
		}
	}
	//set to release mode
	if config.ApplicationConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := initGin()
	//init routes
	for _, f := range AppRouters {
		f() //init all routers!
	}
	//Welcome Page
	routers := global.GetAllRouters()
	fmt.Println("------------------------------------------------------")
	fmt.Println("Welcome to GoOwl! Here're all routes you've registered:")
	for _, v := range routers {
		fmt.Println(v.Route + "---------------->" + v.Explanation)
	}
	//goroutine to use interreput
	go engine.Run(
		fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
	)
	//In order to use ^c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) //STUFF UNTIL CHAN IN
	<-quit
	fmt.Println("\nGoOwl Exit.")
}

//set engine to global
func initGin() *gin.Engine {

	r := gin.Default()
	global.SetEngine(r)
	//	common.InitMiddleware(r)
	return r
}
