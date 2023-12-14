package main

import (
	"template/handlers"
	"template/internal/config"
	"template/internal/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	log := cfg.NewLogger()

	// Init app
	app := echo.New()
	app.Debug = false
	app.HideBanner = true
	app.Logger = lecho.From(*log)
	app.Server.IdleTimeout = cfg.IdleTimeout
	app.Server.ReadTimeout = cfg.ReadTimeout
	app.Server.WriteTimeout = cfg.WriteTimeout

	if cfg.UseGoJSON {
		app.JSONSerializer = &json.GoJSONSerializer{}
	}

	// Init database
	/// ...

	// Use middleware (depending on cfg.Enviornment)
	app.Use(
		cfg.NewLoggerMiddleware(),
		middleware.Recover(),
		middleware.RequestID(),
	)

	// Make routes
	exampleHandler := handlers.NewExampleHandler(log)
	exampleHandler.Bind(app.Group("/users"))

	healthcheckHandler := handlers.NewHealthcheckHandler()
	healthcheckHandler.Bind(app.Group(""))

	// Run
	if err := app.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("fatal server error")
	}
}
