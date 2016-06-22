package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ab22/env"
)

// Config contains all of the configuration variables for the application.
type Config struct {
	SessionCookie struct {
		Name     string `envDefault:"__session"`
		LifeTime time.Duration
	}

	App struct {
		Port string `env:"PORT" envDefault:"1337"`
		Env  string `env:"ENV" envDefault:"DEV"`
	}

	DB struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASS" envDefault:"1234"`
		Name     string `env:"DB_NAME" envDefault:"abcd"`
	}
}

// NewConfig initializes a new Config structure.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err = cfg.validate(); err != nil {
		return nil, err
	}

	cfg.LifeTime = time.Minute * 30

	return cfg, nil
}

func (c *Config) validate() error {
	var errorMsg = "config: field [%v] was not set!"

	// SessionCookie validation.
	if c.SessionCookie.Name == "" {
		return fmt.Errorf(errorMsg, "SessionCookie.Name")
	}

	// App Validation
	if c.App.Port == "" {
		return fmt.Errorf(errorMsg, "App.Port")
	}

	if c.App.Env == "" {
		return fmt.Errorf(errorMsg, "App.Env")
	}

	//Db validation
	if c.Db.Host == "" {
		return fmt.Errorf(errorMsg, "Db.Host")
	}

	if c.Db.Port == 0 {
		return fmt.Errorf(errorMsg, "Db.Port")
	}

	if c.Db.User == "" {
		return fmt.Errorf(errorMsg, "Db.User")
	}

	if c.Db.Password == "" {
		return fmt.Errorf(errorMsg, "Db.Password")
	}

	if c.Db.Name == "" {
		return fmt.Errorf(errorMsg, "Db.Name")
	}

	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	log.Println("----------------------------------")
	log.Println("    Application Port:", c.App.Port)
	log.Println("         Environment:", c.App.Env)
	log.Println(" Session Cookie Name:", c.SessionCookie.Name)
	log.Println("       Database Host:", c.Db.Host)
	log.Println("       Database Port:", c.Db.Port)
	log.Println("       Database Name:", c.Db.Name)
	log.Println("----------------------------------")
}
