package repo

import (
	"strconv"
	"time"

	"github.com/sydneyowl/GoOwl/common/command"
)

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
	// Git rep
	Type string
	// sshkey used for clone.
	Sshkey string
	// Under http protocol
	Username string
	Password string
	// Auth type
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
			if opt.Token != "" { //git pull ....
				if opt.Type == "github" {
					target, err := getOauthRepoURL(targetURL, opt.Token)
					if err != nil {
						return err
					}
					targetURL = target
				} else {
					target, err := getTokenRepoURL(targetURL, opt.Token)
					if err != nil {
						return err
					}
					targetURL = target
				}
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
	_,_, err := cmd.RunInDirWithTimeout(opt.Timeout, dst)
	return err
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
