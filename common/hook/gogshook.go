package hook

type GogsHook struct {
	// branch
	Ref string
	//Before inst. HASHBEFOREACTION.
	Before string
	//After inst. hashAfterAction.
	After  string
	Pusher GogsPusher
}
type GogsPusher struct {
	Username string
}
