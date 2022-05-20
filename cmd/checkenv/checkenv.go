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
	"github.com/sydneyowl/GoOwl/common/stdout"

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
		dockerStat += stdout.Green("True")
		fmt.Println("Checking platform......", detail+dockerStat)
	} else if indocker == 0 {
		warningCount++
		dockerStat += stdout.Red("False")
		fmt.Println("Checking platform......", detail+dockerStat)
		fmt.Println(stdout.Yellow("Warning:GoOwl is not running in docker. Do not run it in host directly!"))
	} else {
		dockerStat += stdout.Yellow("UNKNOWN")
		fmt.Println("Checking platform......", detail+dockerStat)
	}
	if global.OS != "linux" {
		warningCount++
		fmt.Println(
			stdout.Yellow(
				"Warning:GoOwl may crash in " + global.OS + " since it is still unstable.",
			),
		)
	} else {
		rootPermission := envcheck.CheckIsRoot()
		if rootPermission {
			fmt.Println("Checking Permission......", stdout.Green("OK"))
		} else {
			warningCount++
			fmt.Println("Checking Permission......", stdout.Red("NO"))
			fmt.Println(stdout.Yellow("Warning:Run GoOwl without root may crash!"))
		}
	}
	diskSpace := envcheck.CheckDiskSpace()
	if diskSpace {
		fmt.Println("Checking DiskSpace......", stdout.Green("OK"))
	} else {
		warningCount++
		fmt.Println("Checking DiskSpace......", stdout.Red("NO"))
		fmt.Println(stdout.Yellow("Warning:More then 2G Disk space is suggested for GoOwl!"))
	}
	memorySpace := envcheck.CheckMemory()
	if memorySpace {
		fmt.Println("Checking Memory......", stdout.Green("OK"))
	} else {
		warningCount++
		fmt.Println("Checking Memory......", stdout.Red("NO"))
		fmt.Println(stdout.Yellow("Warning:More then 1G Memory is suggested for GoOwl!"))
	}
}

//Check yaml config here；Create example if not exists;check host env only!
func appChecking() error {
	if readable, err := file.CheckYamlReadable(&yamlAddr); !readable {
		fmt.Println("Checking yaml......", stdout.Red("NO"))
		// Deleted since no necessary to generate.
		if os.IsNotExist(err) {
			// fmt.Println(stdout.Yellow("Warning:File not exist.Creating example.yaml for you......"))
			// if isExists, _ := file.CheckPathExists("./config"); isExists {
			// 	if err := config.ReleaseYaml("./config/example.yaml"); err != nil {
			// 		fmt.Println(stdout.Red("ERROR:Failed to release ./config/example.yaml"))
			// 		return errors.New("cannot read file")
			// 	}
			// }
			fmt.Println(stdout.Yellow("Warning:File not exist! Skip..."))
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
			fmt.Println(stdout.Yellow("Warning:File not readable! Skip..."))
			return errors.New("Config file not readable")
		}
	} else {
		rawConfig, err := config.LoadConfigFromYaml(yamlAddr) //returns raw viper obj
		if err := config.CheckViperErr(err); err != nil {
			fmt.Println("Checking yaml......", stdout.Red("NO"))
			fmt.Println(stdout.Magenta("FATAL:" + stdout.Red(err.Error())))
			return errors.New("config error")
		}
		fmt.Println("Checking Yaml......", stdout.Green("ok"))
		port := rawConfig.GetInt("settings.application.port")
		if envcheck.CheckConn("localhost:" + strconv.Itoa(port)) {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", stdout.Red("NO"))
			fmt.Println(stdout.Red("Fatal:Port " + strconv.Itoa(port) + " is being occupied!"))
			return errors.New("port being occupied")
		} else {
			fmt.Println("Checking Port "+strconv.Itoa(port)+".....", stdout.Green("OK"))
		}
		workspaceExists, _ := file.CheckPathExists(rawConfig.GetString("settings.workspace.path"))
		if workspaceExists {
			fmt.Println("Checking workspace......", stdout.Green("OK"))
		} else {
			fmt.Println("Checking workspace......", stdout.Red("NO"))
			fmt.Println(stdout.Yellow("Warning:Workspace does not exist!Create it manually first!"))
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
		fmt.Println(stdout.Cyan("Fatal error occured. Fix it before run checkenv."))
		return
	}
	if warningCount != 0 {
		warnOutput := fmt.Sprintf(
			"Checkenv %v with %v %v. Fix them if you are in production env.",
			stdout.Green("PASSED"),
			warningCount,
			stdout.Yellow("warning"),
		)
		fmt.Println(warnOutput)
	} else {
		fmt.Println("Checkenv", stdout.Green("PASSED")+". You are ready to GoOwl.")
	}
}
