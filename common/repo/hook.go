package repo

type Hook struct {
	// branch
	Ref string
	//Before refers to hash before push.
	Before string
	//After refers to hash after push.
	After  string
	Pusher Pusher
}
type Pusher struct {
	//Github
	Name string
	//Gogs
	Username string
}
