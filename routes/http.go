package routes

import (
	"service/dto"
	g "service/global"
	"service/handlers"
	"service/handlers/auth_handlers"
	"service/middlewares"
	"service/middlewares/extra_middlewares"
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

	// Translator
	app.Use(extra_middlewares.Translator)

	// Panic
	app.Use(extra_middlewares.Panic)

	// Timeout
	app.Use(extra_middlewares.Timeout(time.Second * time.Duration(g.CFG.Timeout)))

	// RateLimiter
	app.Use(extra_middlewares.ConcurrentLimiter(200))

	// Creates a db for every db operation
	app.Use(extra_middlewares.CreateDbInstance)

	// Copression
	app.Use(iris.Compression)
}

func HTTP(app *iris.Application) {
	addMiddlewares(app)

	app.Get("/", handlers.Hello)

	authParty := app.Party("/api/auth")

	registerValidator := middlewares.Validate(dto.RegisterRequestValidator, dto.RegisterRequest{})
	authParty.Post("/register", registerValidator, auth_handlers.Register)

	loginValidator := middlewares.Validate(dto.LoginRequestValidator, dto.LoginRequest{})
	authParty.Post("/login", loginValidator, auth_handlers.Login)

	apiParty := app.Party("/api", middlewares.Auth)

	apiParty.Get("/me", handlers.Me)

	apiParty.Get("/users", handlers.Users)
}
