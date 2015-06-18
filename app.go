package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//"encoding/json"
	"net/http"
	"./metroapi"
)

var app *martini.ClassicMartini

func main() {
	app = martini.Classic()
	app.Use(render.Renderer())

	router := martini.NewRouter()
	router.Get(`/`, func(req *http.Request) (int, string) {
		return http.StatusOK, "Hello world!\n"
	})

	router.Get(`/providers`, func(req *http.Request, r render.Render) {
		providers := metroapi.GetProviders()
		//data, err := json.Marshal(providers); err == nil {
			r.JSON(200, providers)
	})

	app.Action(router.Handle)

	http.ListenAndServe(":3000", app)
}
