package app

import (
	"service/routes"

	"github.com/kataras/iris/v12"
)

func API() {
	// Print Info
	info()

	app := iris.New()
	app.Configure(iris.WithoutStartupLog)

	// Router Settings
	routes.HTTP(app)

	// Run App
	app.Listen(":8080")
}
