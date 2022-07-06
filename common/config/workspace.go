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
	//0->Unknown 1->Failing 2->Waiting 3->Passing
	BuildStatus int
}

var WorkspaceConfig = new(WorkspaceSettings)
