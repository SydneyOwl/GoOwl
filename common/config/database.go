package config

import "time"
type BuildInfo struct {
	RepoID string
	CreatedAt time.Time
	BuildStatus int
	Output string
	TimeCost int64
}
type TriggerInfo struct{
	IsAvailableAction uint8
	Action string
	Branch string
	HookType string
	CreatedAt time.Time
	TriggerBy string
	HashBeforeTrigger string
	HashAfterTrigger string
}