package main

import (
	"template/handlers"
	"template/internal/config"
	"template/internal/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	log := cfg.NewLogger()
	log.Debug().Msg("debug logs are enabled")

	// Init app
	app := echo.New()
	app.Debug = false
	app.HideBanner = true
	app.HidePort = true
	app.Server.IdleTimeout = cfg.IdleTimeout
	app.Server.ReadTimeout = cfg.ReadTimeout
	app.Server.WriteTimeout = cfg.WriteTimeout

	if cfg.UseGoJSON {
		log.Info().Msg("using go-json")
		app.JSONSerializer = &json.GoJSONSerializer{}
	}

	// Init database
	/// ...

	// Use middleware (depending on cfg.Enviornment)
	app.Use(
		middleware.RequestID(),
		cfg.NewLoggerMiddleware(),
		middleware.Recover(),
	)

	// Make routes
	exampleHandler := handlers.NewExampleHandler(log)
	exampleHandler.Bind(app.Group("/users"))

	healthcheckHandler := handlers.NewHealthcheckHandler()
	healthcheckHandler.Bind(app.Group(""))

	// Run
	log.Info().Str("host", cfg.Host).Str("port", cfg.Port).Str("link", cfg.Host+":"+cfg.Port).Msg("server starting")
	if err := app.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("fatal server error")
	}
}
