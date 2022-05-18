package config

type Total struct {
	Settings Settings
}
type Settings struct {
	Application *ApplicationSettings `yaml:"application"`
	Workspace   *WorkspaceSettings   `yaml:"workspace"`
}

var YamlConfig = &Total{
	Settings: Settings{
		Application: ApplicationConfig,
		Workspace:   WorkspaceConfig,
	},
}
