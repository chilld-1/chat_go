package api

import "gochat/api/router"

type API struct{}

func New() *API {
	return &API{}
}
func (a *API) Run() {
	r := router.SetupRouter()
	r.Run(":8081")
}
