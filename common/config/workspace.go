package config

type WorkspaceSettings struct {
	Path string
	Repo []Repo `yaml:"repo"`
}
type Repo struct {
	ID          string
	Type        string
	Trigger     []string
	Repoaddr    string
	Branch      string
	Sshkeyaddr  string
	Username    string
	Password    string
	Token       string
	Buildscript string
}

var WorkspaceConfig = new(WorkspaceSettings)
