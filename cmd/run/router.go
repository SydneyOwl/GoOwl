package run

import "github.com/sydneyowl/GoOwl/app/other/router"

func init() {
	//add initrouter func to here append
	AppRouters = append(AppRouters, router.InitAllRouter)
}
