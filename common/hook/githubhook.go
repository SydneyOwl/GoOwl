package hook

type GithubHook struct {
	// branch
	Ref string
	//Before inst. HASHBEFOREACTION.
	Before string
	//After inst. hashAfterAction.
	After  string
	Pusher GithubPusher
}
type GithubPusher struct {
	Name string
}
