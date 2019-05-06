package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

// Config stores the handler's configuration and UI interface parameters
type Config struct {
	Version bool `short:"V" long:"version" description:"Display version."`

	Port int `short:"p" long:"port" description:"Port to listen on." default:"8080"`

	Log struct {
		Level  string `short:"l" long:"log-level" description:"Set log level ('debug', 'info', 'warn', 'error', 'fatal', 'panic')." env:"TERRABOARD_LOG_LEVEL" default:"info"`
		Format string `long:"log-format" description:"Set log format ('plain', 'json')." env:"TERRABOARD_LOG_FORMAT" default:"plain"`
	} `group:"Logging Options"`

	TerraDB struct {
		URL string `long:"terradb-url" env:"TERRABOARD_TERRADB_URL"`
	} `group:"TerraDB"`

	Authentication struct {
		LogoutURL string `long:"logout-url" env:"TERRABOARD_LOGOUT_URL" description:"Logout URL."`
	} `group:"Authentication"`
}

// LoadConfig loads the config from flags & environment
func LoadConfig(version string) *Config {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	if c.Version {
		fmt.Printf("Terraboard v%v\n", version)
		os.Exit(0)
	}

	return &c
}

// SetupLogging sets up logging for Terraboard
func (c Config) SetupLogging() (err error) {
	switch c.Log.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		errMsg := fmt.Sprintf("Wrong log level '%v'", c.Log.Level)
		return errors.New(errMsg)
	}

	switch c.Log.Format {
	case "plain":
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		errMsg := fmt.Sprintf("Wrong log format '%v'", c.Log.Format)
		return errors.New(errMsg)
	}

	return
}
