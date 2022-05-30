package logger

import (
	"bufio"
	"fmt"
	"os"
	_ "strconv"
	"time"

	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/global"

	"github.com/fatih/color"
)

func CreateOnNotExist(id string) error {
	if global.LoggingMethod != 1 {
		file.CreateDir(file.GetCwd() + "/log")
		now:=time.Now()
	curDate := fmt.Sprintf("%d-%d-%d",now.Year(),now.Month(),now.Day())
		err := file.CreateFile(file.GetCwd() + "/log/" + fmt.Sprintf("[%s]Repo_%s", curDate, id))
		if err != nil {
			return err
		}
	}
	return nil
}

func AppendLog(id string, msg string) error {
	now:=time.Now()
	curDate := fmt.Sprintf("%d-%d-%d",now.Year(),now.Month(),now.Day())
	filePtr, err := os.OpenFile(file.GetCwd()+"/log/"+fmt.Sprintf("[%s]Repo_%s", curDate, id), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			CreateOnNotExist(id)
			filePtr, _ = os.OpenFile(file.GetCwd()+"/log/"+fmt.Sprintf("[%s]Repo_%s", curDate, id), os.O_WRONLY|os.O_APPEND, 0666)
		} else {
			return err
		}
	}
	defer filePtr.Close()
	write := bufio.NewWriter(filePtr)
	write.WriteString(msg + "\n")
	write.Flush()
	return nil
}

// red return string in red
func red(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgRed).SprintFunc()(msg)
	}
	return msg
}

// green return string in green
func green(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgGreen).SprintFunc()(msg)
	}
	return msg
}

// yellow return string in yellow
func yellow(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgYellow).SprintFunc()(msg)
	}
	return msg
}

// blue return string in blue
func blue(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgBlue).SprintFunc()(msg)
	}
	return msg
}

// magenta return string in Magenta
func magenta(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgHiMagenta).SprintFunc()(msg)
	}
	return msg
}

// cyan return string in cyan
func cyan(msg string, addColor bool) string {
	if addColor && global.OS == "linux" {
		return color.New(color.FgCyan).SprintFunc()(msg)
	}
	return msg
}

func logFactory(msg interface{}, id string, level string) {
	CreateOnNotExist(id)
	timestr := time.Now().Format("2006-01-02 15:04:05")
	var logInfo string
	//info with color '1
	var infocolor string
	switch level {
	case "Debug":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Debug-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Debug-Repo %s] %s", timestr, id, msg)
		}
		infocolor = green(logInfo, true)
	case "Notice":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Notice-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Notice-Repo %s] %s", timestr, id, msg)
		}
		infocolor = blue(logInfo, true)
	case "Warning":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Warning-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Warning-Repo %s] %s", timestr, id, msg)
		}
		infocolor = yellow(logInfo, true)
	case "Error":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Error-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Error-Repo %s] %s", timestr, id, msg)
		}
		infocolor = cyan(logInfo, true)
	case "Fatal":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Fatal-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Fatal-Repo %s] %s", timestr, id, msg)
		}
		infocolor = red(logInfo, true)
	case "Info":
		if id == "GoOwl-MainLog" {
			logInfo = fmt.Sprintf("%s [Info-GoOwl] %s", timestr, msg)
		} else {
			logInfo = fmt.Sprintf("%s [Info-Repo %s] %s", timestr, id, msg)
		}
		infocolor =logInfo
	}
	if global.LoggingMethod == 1 {
		fmt.Println(infocolor)
	} else if global.LoggingMethod == 2 {
		AppendLog(id, logInfo)
	} else {
		fmt.Println(infocolor)
		AppendLog(id, logInfo)
	}
}
func Debug(msg interface{}, id string) {
	logFactory(msg,id,"Debug")
}
func Info(msg interface{},id string){
	logFactory(msg,id,"Info")
}
func Notice(msg interface{}, id string) {
	logFactory(msg,id,"Notice")
}
func Warning(msg interface{}, id string) {
	logFactory(msg,id,"Warning")
}
func Error(msg interface{}, id string) {
	logFactory(msg,id,"Error")
}
func Fatal(msg interface{}, id string) {
	logFactory(msg,id,"Fatal")
}