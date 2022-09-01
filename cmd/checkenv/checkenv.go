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
		fmt.Println(
			"Warning:GoOwl may crash in " + global.OS + " since it is still unstable.",
		)
	} else {
		rootPermission := envcheck.CheckIsRoot()
		if rootPermission {
			fmt.Println("Checking Permission......", "OK")
		} else {
			warningCount++
			fmt.Println("Checking Permission......", "NO")
			fmt.Println("Warning: Run GoOwl without root may crash!")
		}
	}
	diskSpace := envcheck.CheckDiskSpace()
	if diskSpace {
		fmt.Println("Checking DiskSpace......", "OK")
	} else {
		warningCount++
		fmt.Println("Checking DiskSpace......", "NO")
		fmt.Println("Warning: More then 2G Disk space is suggested for GoOwl!")
	}
	memorySpace := envcheck.CheckMemory()
	if memorySpace {
		fmt.Println("Checking Memory......", "OK")
	} else {
		warningCount++
		fmt.Println("Checking Memory......", "NO")
		fmt.Println("Warning:More then 1G Memory is suggested for GoOwl!")
	}
}

// appChecking check yaml config here；Create example if not exists;check host env only!
func appChecking() error {
	if readable, err := file.CheckYamlReadable(&yamlAddr); !readable {
		fmt.Println("Checking yaml......", "NO")
		// Deleted since no necessary to generate.
		if os.IsNotExist(err) {
			fmt.Println("Warning: File (yaml) not exist! Skip...")
			return errors.New("config file not found")
		} else {
			fmt.Println("Warning:File not readable! Skip...")
			return errors.New("config file not readable")
		}
	} else {
		rawConfig, err := config.LoadConfigFromYaml(yamlAddr) //returns raw viper obj
		if err := config.CheckViperErr(err); err != nil {
			fmt.Println("Checking yaml......", "NO")
			return errors.New("config error")
		}
		fmt.Println("Checking Yaml......", "ok")
		port := rawConfig.GetInt("settings.application.port")
		if envcheck.CheckConn("localhost:" + strconv.Itoa(port)) {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", "NO")
			fmt.Println("Warning: Port " + strconv.Itoa(port) + " is being occupied!")
			return errors.New("port being occupied")
		} else {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", "OK")
		}
		workspaceExists, _ := file.CheckPathExists(rawConfig.GetString("settings.workspace.path"))
		if workspaceExists {
			fmt.Println("Checking workspace......", "OK")
		} else {
			fmt.Println("Checking workspace......", "NO")
			fmt.Println("Warning: Workspace does not exist!Create it manually first!")
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
		fmt.Println("Fatal error occurred. Fix it before run checkenv.")
		return
	}
	if warningCount != 0 {
		warnOutput := fmt.Sprintf(
			"Checkenv %v with %v %v. Fix them if you are in production env.",
			"PASSED",
			warningCount,
			"warning",
		)
		fmt.Println(warnOutput)
	} else {
		fmt.Println("checkenv PASSED. You are ready to GoOwl.")
	}
}
