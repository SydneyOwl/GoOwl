package hook

type GogsHook struct {
	// branch
	Ref string
	//Before refers to hash before push.
	Before string
	//After refers to hash after push.
	After  string
	Pusher GogsPusher
}
type GogsPusher struct {
	Username string
}
