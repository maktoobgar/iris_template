package routes

import (
	g "service/global"
	"service/handlers"
	"service/middlewares"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

// Applies all necessary middlewares
func addMiddlewares(app *iris.Application) {
	// Cors
	c := cors.New(cors.Options{
		AllowedOrigins: strings.Split(g.CFG.AllowOrigins, ","),
		AllowedHeaders: strings.Split(g.CFG.AllowHeaders, ","),
	})
	app.WrapRouter(c.ServeHTTP)

	// Json
	app.Use(middlewares.Json)

	// Translator
	app.Use(middlewares.Translator)

	// Panic
	app.Use(middlewares.Panic)

	// Timeout
	app.Use(middlewares.Timeout(time.Second * time.Duration(g.CFG.Timeout)))

	// RateLimiter
	app.Use(middlewares.ConcurrentLimiter(200))

	// Copression
	app.Use(iris.Compression)
}

func HTTP(app *iris.Application) {
	addMiddlewares(app)

	app.Get("/", handlers.Hello)
}
