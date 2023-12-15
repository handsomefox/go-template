package config

import (
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvStaging     Environment = "staging"
)

type Application struct {
	DatabaseConnectionString string      `env:"POSTGRES_CONN_STRING" env-required:"true"       yaml:"postgresConn"`
	Environment              Environment `env:"APP_ENV"              env-default:"development" yaml:"env"`
	JWTIssuer                string      `env:"JWT_ISSUER"           env-required:"true"       yaml:"jwtIssuer"`
	JWTSecret                string      `env:"JWT_SECRET"           env-required:"true"       yaml:"jwtSecret"`
	SessionKey               string      `env:"SESSION_KEY"          yaml:"sessionKey"`
	UseGoJSON                bool        `env:"GOJSON"               yaml:"useGoJson"`
	HTTPServer               `yaml:"httpServer"`
}

type HTTPServer struct {
	Host         string        `env:"HOST"        yaml:"host"`
	IdleTimeout  time.Duration `env-default:"60s" yaml:"idleTimeout"`
	Port         string        `env:"PORT"        yaml:"port"`
	ReadTimeout  time.Duration `env-default:"5s"  yaml:"readTimeout"`
	WriteTimeout time.Duration `env-default:"10s" yaml:"writeTimeout"`
}

func MustLoad() *Application {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal().Msg("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Str("CONFIG_PATH", configPath).Err(err).Send()
	}

	var cfg Application
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to read the config")
	}

	return &cfg
}

func (cfg *Application) NewLogger() *zerolog.Logger {
	const maxPathElements = 4
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		short := ""
		appended := 0
		prevSlash := len(file)
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i:prevSlash] + short
				prevSlash = i
				appended++
			}
			if appended == maxPathElements {
				if short[0] == '/' {
					short = short[1:]
				}
				break
			}
		}
		return short + ":" + strconv.Itoa(line)
	}

	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Stack().Logger()
	switch cfg.Environment {
	case EnvDevelopment:
		logger = logger.Level(zerolog.DebugLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case EnvStaging:
		logger = logger.Level(zerolog.InfoLevel)
	case EnvProduction:
		logger = logger.Level(zerolog.WarnLevel)
	}

	return &logger
}

func (cfg *Application) NewLoggerMiddleware() echo.MiddlewareFunc {
	l := cfg.NewLogger()
	return lecho.Middleware(lecho.Config{Logger: lecho.From(*l)})
}
