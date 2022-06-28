package checkenv

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/envcheck"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/global"
	"github.com/sydneyowl/GoOwl/common/logger"

	"github.com/spf13/cobra"
)

var (
	yamlAddr     string
	warningCount int
	StartCmd     = &cobra.Command{
		Use:     "checkenv",
		Short:   "Check whether GoOwl is compatible with this platform",
		Example: "GoOwl checkenv -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			runEnvCheck()
		},
	}
)

func init() {
	StartCmd.Flags().
		StringVarP(&yamlAddr, "checkenv", "c", "", "Check settings with specified yaml")
}

// basicChecking checks environment.
func basicChecking() {
	//write platform to global
	global.OS = envcheck.CheckOS()
	arch := envcheck.CheckArch()
	fmt.Println("Checking SysArch......", arch)
	dockerStat := ",InDocker:"
	detail := envcheck.CheckPlatform()
	indocker := envcheck.CheckDocker()
	if indocker == 1 {
		dockerStat += "True"
		fmt.Println("Checking platform......", detail+dockerStat)
	} else if indocker == 0 {
		warningCount++
		dockerStat += "False"
		fmt.Println("Checking platform......", detail+dockerStat)
		fmt.Println("Warning:GoOwl is not running in docker. Do not run it in host directly!")
	} else {
		dockerStat += "UNKNOWN"
		fmt.Println("Checking platform......", detail+dockerStat)
	}
	if global.OS != "linux" {
		warningCount++
		logger.Warning(
			"Warning:GoOwl may crash in "+global.OS+" since it is still unstable.", "GoOwl-MainLog",
		)
	} else {
		rootPermission := envcheck.CheckIsRoot()
		if rootPermission {
			fmt.Println("Checking Permission......", "OK")
		} else {
			warningCount++
			fmt.Println("Checking Permission......", "NO")
			logger.Warning("Run GoOwl without root may crash!", "GoOwl-MainLog")
		}
	}
	diskSpace := envcheck.CheckDiskSpace()
	if diskSpace {
		fmt.Println("Checking DiskSpace......", "OK")
	} else {
		warningCount++
		fmt.Println("Checking DiskSpace......", "NO")
		logger.Warning("More then 2G Disk space is suggested for GoOwl!", "GoOwl-MainLog")
	}
	memorySpace := envcheck.CheckMemory()
	if memorySpace {
		fmt.Println("Checking Memory......", "OK")
	} else {
		warningCount++
		fmt.Println("Checking Memory......", "NO")
		logger.Warning("Warning:More then 1G Memory is suggested for GoOwl!", "GoOwl-MainLog")
	}
}

//appChecking check yaml config here；Create example if not exists;check host env only!
func appChecking() error {
	if readable, err := file.CheckYamlReadable(&yamlAddr); !readable {
		fmt.Println("Checking yaml......", "NO")
		// Deleted since no necessary to generate.
		if os.IsNotExist(err) {
			// fmt.Println(stdout.Yellow("Warning:File not exist.Creating example.yaml for you......"))
			// if isExists, _ := file.CheckPathExists("./config"); isExists {
			// 	if err := config.ReleaseYaml("./config/example.yaml"); err != nil {
			// 		fmt.Println(stdout.Red("ERROR:Failed to release ./config/example.yaml"))
			// 		return errors.New("cannot read file")
			// 	}
			// }
			logger.Warning("File not exist! Skip...", "GoOwl-MainLog")
			return errors.New("Config file not found")
			// } else {
			// 	if err := os.Mkdir("./config", 0777); err != nil {
			// 		fmt.Println(stdout.Red("ERROR:Failed to create ./config"))
			// 		return errors.New("cannot read file")
			// 	}
			// 	if err := config.ReleaseYaml("./config/example.yaml"); err != nil {
			// 		fmt.Println(stdout.Red("ERROR:Failed to create ./config/example.yaml"))
			// 		return errors.New("cannot read file")
			// 	}
			// }
			// fmt.Println(stdout.Yellow("Done.Now modify ./config/example.yaml and run again."))
			// return errors.New("cannot read file")
		} else {
			logger.Warning("Warning:File not readable! Skip...", "GoOwl-MainLog")
			return errors.New("config file not readable")
		}
	} else {
		rawConfig, err := config.LoadConfigFromYaml(yamlAddr) //returns raw viper obj
		if err := config.CheckViperErr(err); err != nil {
			fmt.Println("Checking yaml......", "NO")
			logger.Fatal(err.Error(), "GoOwl-MainLog")
			return errors.New("config error")
		}
		fmt.Println("Checking Yaml......", "ok")
		port := rawConfig.GetInt("settings.application.port")
		if envcheck.CheckConn("localhost:" + strconv.Itoa(port)) {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", "NO")
			logger.Fatal("Port "+strconv.Itoa(port)+" is being occupied!", "GoOwl-MainLog")
			return errors.New("port being occupied")
		} else {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", "OK")
		}
		workspaceExists, _ := file.CheckPathExists(rawConfig.GetString("settings.workspace.path"))
		if workspaceExists {
			fmt.Println("Checking workspace......", "OK")
		} else {
			fmt.Println("Checking workspace......", "NO")
			logger.Warning("Workspace does not exist!Create it manually first!", "GoOwl-MainLog")
			return errors.New("workspace not exist")
		}
		return nil
	}
}

// runEnvCheck runs checkenv process.
func runEnvCheck() {
	basicChecking()
	fmt.Println("-------------------------------------------")
	if err := appChecking(); err != nil { //All fatal errors.
		logger.Notice("Fatal error occurred. Fix it before run checkenv.", "GoOwl-MainLog")
		return
	}
	if warningCount != 0 {
		warnOutput := fmt.Sprintf(
			"Checkenv %v with %v %v. Fix them if you are in production env.",
			"PASSED",
			warningCount,
			"warning",
		)
		logger.Warning(warnOutput, "GoOwl-MainLog")
	} else {
		logger.Warning("Checkenv PASSED"+". You are ready to GoOwl.", "GoOwl-MainLog")
	}
}
