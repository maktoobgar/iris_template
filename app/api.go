package app

import (
	g "service/global"
	"service/routes"

	"github.com/kataras/iris/v12"
)

func API() {
	// Print Info
	info()

	app := iris.New()
	app.Configure(iris.WithoutStartupLog)

	// Router Settings
	g.App = app
	routes.HTTP(app)

	RunClonesAndServer(app)
}
