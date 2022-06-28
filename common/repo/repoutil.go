package repo

import (
	"errors"
	"fmt"
	UrlParse "net/url"
	"path"
	"strings"
	"time"

	"github.com/sydneyowl/GoOwl/common/command"
	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/logger"
)

//Errors that does not infect cicd.
type UncriticalError struct {
	Uerror error
	ID     string
}

var (
	token_supported []string = []string{"gogs", "github"}
)

// getHttpRepoURL returns url include username and password.
func getHttpRepoURL(url string, username string, password string) (string, error) {
	urlParse, err1 := UrlParse.Parse(url)
	if err1 != nil {
		return "", err1
	}
	username_prcessed := strings.ReplaceAll(
		username,
		"@",
		"%40",
	) //Replace @ if mail used as username
	return fmt.Sprintf(
		"%s://%s:%s@%s",
		urlParse.Scheme,
		username_prcessed,
		password,
		urlParse.Host+urlParse.Path,
	), nil
}

//getTokenRepoURL returns token format url.
func getTokenRepoURL(url string, token string) (string, error) {
	urlParse, err1 := UrlParse.Parse(url)
	if err1 != nil {
		return "", err1
	}
	return fmt.Sprintf(
		"%s://%s@%s",
		urlParse.Scheme,
		token,
		urlParse.Host+urlParse.Path,
	), nil
}

// getOauthRepoURL returns oauth format url.(github)
func getOauthRepoURL(url string, token string) (string, error) {
	urlParse, err1 := UrlParse.Parse(url)
	if err1 != nil {
		return "", err1
	}
	return fmt.Sprintf(
		"%s://oauth2:%s@%s",
		urlParse.Scheme,
		token,
		urlParse.Host+urlParse.Path,
	), nil
}

// CheckRepoConfig checks if any attr is empty.
func CheckRepoConfig(repoarray []config.Repo) (string, []UncriticalError, error) {
	var bser []UncriticalError
	for _, v := range repoarray {
		existsScript, _ := file.CheckPathExists(v.Buildscript)
		if v.Buildscript == "" || !existsScript {
			// UncriticalError
			bser = append(bser, UncriticalError{
				ID:     v.ID,
				Uerror: errors.New("buildScript is empty"),
			})
		}
		if v.Repoaddr == "" {
			return v.ID, nil, errors.New("repoaddr not set")
		}
		if v.Branch == "" {
			bser = append(bser, UncriticalError{
				ID:     v.ID,
				Uerror: errors.New("branch not set. Use master in default"),
			})
		}
		if v.Type == "" {
			return v.ID, nil, errors.New("hooktype not specified")
		}
		if v.Trigger == nil {
			v.Trigger = []string{"push"} //Set push
		}
		if isPublicRepo(v) {
			bser = append(bser, UncriticalError{
				ID: v.ID,
				Uerror: errors.New(
					"repo delare itself as public since neither username/password nor token is specified",
				),
			})
			continue //Ignore since it is an public repo
		}
		if v.Sshkeyaddr == "" && Checkprotocol(v) == "http" { //http protocol
			if v.Type == "github" && v.Token == "" {
				return v.ID, nil, errors.New("github supports personal token to access repo only")
			}
			if (v.Username == "" || v.Password == "") && v.Token == "" {
				return v.ID, nil, errors.New("no valid authorization method found in config")
			}
			if v.Token != "" && !config.CheckInSlice(token_supported, v.Type) {
				return v.ID, nil, fmt.Errorf("type %s does not support token authorization", v.Type)
			}
		} else if Checkprotocol(v) == "ssh" && v.Sshkeyaddr != "" { //ssh protocol
			if exists, _ := file.CheckPathExists(v.Sshkeyaddr); !exists {
				return v.ID, nil, errors.New("sshkey not found in " + v.Sshkeyaddr)
			}
		} else {
			bser = append(bser, UncriticalError{
				ID:     v.ID,
				Uerror: errors.New("mix use of http and ssh. GoOwl use ssh by default"),
			})
		}
		if v.Token != "" && (v.Username != "" || v.Password != "") { //exists at the same time
			bser = append(bser, UncriticalError{
				ID: v.ID,
				Uerror: errors.New(
					"both username and token are specified. GoOwl uses token in default",
				),
			})
		}
	}
	//return only if there're no more critial error gened.
	if len(bser) != 0 {
		return "", bser, nil
	}
	return "", nil, nil
}

// Checkprotocol checks protocol
func Checkprotocol(v config.Repo) string {
	if strings.Contains(v.Repoaddr, "http") {
		return "http"
	}
	return "ssh"
}

// isPublicRepo checks if the repo is public. Used in checkinfo.
func isPublicRepo(v config.Repo) bool {
	return Checkprotocol(v) == "http" && v.Token == "" && v.Username == "" && v.Password == ""
}

// isPublicRepo checks if the repo is public. Used in cloneoptions.
func (opts CloneOptions) isPublicRepo() bool {
	return opts.Protocol == "http" && opts.Token == "" && opts.Username == "" && opts.Password == ""
}

// isPublicRepo checks if the repo is public. Used in pulloptions.
func (opts PullOptions) isPublicRepo() bool {
	return opts.Protocol == "http" && opts.Token == "" && opts.Username == "" && opts.Password == ""
}

// LocalRepoAddr returns path repo storage in.
func LocalRepoAddr(repo config.Repo) string {
	repoarr := strings.Split(repo.Repoaddr, "/")
	reponame := repoarr[len(repoarr)-1]
	realpath := path.Join(config.WorkspaceConfig.Path, reponame)
	return realpath
}

// GetRepoName returns reponame.
func GetRepoName(repo config.Repo) string {
	repoarr := strings.Split(repo.Repoaddr, "/")
	return repoarr[len(repoarr)-1]
}

// GetRepoOriginal returns reponame withoutgit.
func GetRepoOriginalName(repo config.Repo) string {
	name := GetRepoName(repo)
	return name[0 : len(name)-4]
}

//CloneOnNotExist clone repo not exist locally
func CloneOnNotExist(repo config.Repo) error {
	localAddr := LocalRepoAddr(repo)
	exists, err := file.CheckPathExists(localAddr)
	if err != nil {
		return err
	}
	if exists {
		logger.Info("Repo "+repo.ID+" already exists. Passing......", "GoOwl-MainLog")
		return nil
	}

	logger.Info("Cloning repo"+repo.ID+"...", repo.ID)
	option := CloneOptions{
		Branch: repo.Branch,
	}
	if Checkprotocol(repo) == "http" {
		option.Protocol = "http"
		if repo.Token != "" {
			option.Token = repo.Token
		} else {
			option.Username = repo.Username
			option.Password = repo.Password
		}
	} else {
		option.Protocol = "ssh"
		option.Sshkey = repo.Sshkeyaddr
	}
	return clone(repo.Repoaddr, localAddr, option)
}

//Searchinfo returns repo with specified id.
func SearchRepo(ID string) (config.Repo, error) {
	for _, v := range config.WorkspaceConfig.Repo {
		if v.ID == ID {
			return v, nil
		}
	}
	return config.Repo{
		ID: "",
	}, errors.New("no found")
}

//Runscript run script inside repo dir.
func RunScript(repo config.Repo) (string, error) {
	if repo.Buildscript == "" {
		return "", fmt.Errorf(
			"buildscript of repo %s (%v)is empty. CI suspended",
			repo.ID,
			GetRepoName(repo),
		)
	}
	command := command.CICDCommand(repo.Buildscript)
	result, err := command.RunInDirWithTimeout(time.Hour, LocalRepoAddr(repo)) //Hour of the timeout
	return string(result), err
}

// IsDuplcatedRepo check if repo is dupl in config.
func IsDuplcatedRepo(repos []config.Repo) (bool, error) {
	for i := 0; i < len(repos); i++ {
		for j := i + 1; j < len(repos); j++ {
			if repos[i].Repoaddr == repos[j].Repoaddr {
				return true, fmt.Errorf("duplcate repo address found:%v", repos[i].Repoaddr)
			}
			if repos[i].ID == repos[j].ID {
				return true, fmt.Errorf("duplcate repo id found:%v", repos[i].ID)
			}
		}
	}
	return false, nil
}
func CheckRepo() {
	if repeated, err := IsDuplcatedRepo(config.WorkspaceConfig.Repo); repeated {
		fmt.Println(err.Error())
		return
	}
	ID, uncritialerror, err := CheckRepoConfig(config.WorkspaceConfig.Repo)
	if err != nil {
		logger.Error("Repo has an invaild config:"+err.Error(), ID)
		return
	}
	if len(uncritialerror) > 0 {
		for _, v := range uncritialerror {
			logger.Warning(
				"repo  has an invaild config:"+v.Uerror.Error()+",check if it is correct.",
				v.ID,
			)
		}
	}
}

//SetBuildStat modify build status of repo.
func SetBuildStat(id string, stat int) {
	for i, v := range config.WorkspaceConfig.Repo {
		if v.ID == id {
			config.WorkspaceConfig.Repo[i].BuildStatus = stat
			break
		}
	}
}
