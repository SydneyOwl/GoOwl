package repo

import (
	"errors"
	"fmt"
	UrlParse "net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/sydneyowl/GoOwl/common/command"
	"github.com/sydneyowl/GoOwl/common/config"
	"github.com/sydneyowl/GoOwl/common/file"
)

//Errors that does not infect cicd.
type UncriticalError struct {
	Uerror error
	ID     string
}

//Gogs code under mit.
type CloneOptions struct {
	//hook type
	Type string
	//specify protocol
	Protocol string
	// Indicates whether the repository should be cloned as a mirror.
	Mirror bool
	// sshkey used for clone.
	Sshkey string
	// Indicates whether the repository should be cloned in bare format.
	Bare bool
	// Indicates whether to suppress the log output.
	Quiet bool
	// The branch to checkout for the working tree when Bare=false.
	Branch string
	//Under http protocol
	Username string
	Password string
	//token (github)
	Token string
	// The number of revisions to clone.
	Depth uint64
	// The timeout duration before giving up for each shell command execution. The
	// default timeout duration will be used when not supplied.
	Timeout time.Duration
}

type PullOptions struct {
	//specify protocol
	Protocol string
	// Indicates whether to rebased during pulling.
	Rebase bool
	//
	Type string
	// sshkey used for clone.
	Sshkey string
	//Under http protocol
	Username string
	Password string
	//github only
	Token string
	// Indicates whether to pull from all remotes.
	All bool
	// The remote to pull updates from when All=false.
	Remote string
	// The branch to pull updates from when All=false and Remote is supplied.
	Branch string
	// The timeout duration before giving up for each shell command execution. The
	// default timeout duration will be used when not supplied.
	Timeout time.Duration
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

// Pull pulls updates for the repository.
func Pull(dst string, opts ...PullOptions) error {
	var opt PullOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	var cmd *command.Command
	targetURL := opt.Remote
	if opt.Protocol == "ssh" {
		cmd = command.SSHCommand(opt.Sshkey, "pull") //to clone things using ssh
	} else { //https
		cmd = command.NewCommand("pull")
		if !opt.isPublicRepo() {
			if targetURL != "" { //git pull ....
				if opt.Token != "" {
					target, err := getOauthRepoURL(targetURL, opt.Token)
					if err != nil {
						return err
					}
					targetURL = target
				} else {
					target, err := getHttpRepoURL(targetURL, opt.Username, opt.Password)
					if err != nil {
						return err
					}
					targetURL = target
				}
				// username := strings.ReplaceAll(opt.Username, "@", "%40") //Replace @ if mail used as username
				// targetURL = fmt.Sprintf("%s://%s:%s@%s", urlParse.Scheme, username, opt.Password, urlParse.Host+urlParse.Path)
			}
		}
	}
	if opt.Rebase {
		cmd.AddArgs("--rebase")
	}
	if opt.All {
		cmd.AddArgs("--all")
	}
	if !opt.All && opt.Remote != "" {
		cmd.AddArgs(targetURL)
		if opt.Branch != "" {
			cmd.AddArgs(opt.Branch)
		}
	}
	_, err := cmd.RunInDirWithTimeout(opt.Timeout, dst)
	return err
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
				ID:     v.ID,
				Uerror: errors.New("repo delare itself as public since neither username/password nor token is specified"),
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

// clone clones the repository from remote URL to the destination.
func clone(url, dst string, opts ...CloneOptions) error {
	var opt CloneOptions
	var targetURL = url
	if len(opts) > 0 {
		opt = opts[0]
	}
	var cmd *command.Command
	//token->username
	if opt.Protocol == "ssh" {
		cmd = command.SSHCommand(opt.Sshkey, "clone") //to clone things using ssh
	} else { //httpprot
		cmd = command.NewCommand("clone")
		if !opt.isPublicRepo() {
			// Ignore empty since checked.
			if opt.Token != "" {
				if opt.Type == "github" {
					target, err := getOauthRepoURL(url, opt.Token)
					if err != nil {
						return err
					}
					targetURL = target
				} else {
					target, err := getTokenRepoURL(url, opt.Token)
					if err != nil {
						return err
					}
					targetURL = target
				}
			} else {
				target, err := getHttpRepoURL(url, opt.Username, opt.Password)
				if err != nil {
					return err
				}
				targetURL = target
			}
		}
	}
	if opt.Mirror {
		cmd.AddArgs("--mirror")
	}
	if opt.Bare {
		cmd.AddArgs("--bare")
	}
	if opt.Quiet {
		cmd.AddArgs("--quiet")
	}
	if !opt.Bare && opt.Branch != "" {
		cmd.AddArgs("-b", opt.Branch)
	}
	if opt.Depth > 0 {
		cmd.AddArgs("--depth", strconv.FormatUint(opt.Depth, 10))
	}
	_, err := cmd.AddArgs(targetURL, dst).RunWithTimeout(opt.Timeout)
	return err
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

//CloneOnNotExist clone repo not exist locally
func CloneOnNotExist(repo config.Repo) error {
	localAddr := LocalRepoAddr(repo)
	exists, err := file.CheckPathExists(localAddr)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Repo", repo.ID, "already exists. Passing......")
		return nil
	}

	fmt.Println("Cloning repo", repo.ID, "...")
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
	}, nil
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
