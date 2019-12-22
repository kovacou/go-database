// Copyright © 2019 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kovacou/go-env"
)

// Default environment configuration.
const (
	defaultDriver      = "mysql"
	defaultProtocol    = "tcp"
	defaultCharset     = "utf8mb4"
	defaultHost        = "172.18.0.1"
	defaultPort        = "3306"
	defaultMaxIdle     = 1
	defaultMaxOpen     = 2
	defaultMaxLifetime = 1800 * time.Second
)

// Environment store the configuration to open a new connection.
type Environment struct {
	Alias string
	DSN   string `env:"DATABASE_DSN"`

	Driver         string        `env:"DATABASE_DRIVER"`
	Protocol       string        `env:"DATABASE_PROTOCOL"`
	Host           string        `env:"DATABASE_HOST"`
	Port           string        `env:"DATABASE_PORT"`
	User           string        `env:"DATABASE_USER"`
	Pass           string        `env:"DATABASE_PASS"`
	Charset        string        `env:"DATABASE_CHARSET"`
	Schema         string        `env:"DATABASE_SCHEMA"`
	Mode           string        `env:"DATABASE_MODE"`
	ParseTime      bool          `env:"DATABASE_PARSETIME"`
	Autoconnect    bool          `env:"DATABASE_AUTOCONNECT"`
	MaxOpen        int           `env:"DATABASE_MAXOPEN"`
	MaxIdle        int           `env:"DATABASE_MAXIDLE"`
	MaxLifetime    time.Duration `env:"DATABASE_MAXLIFETIME"`
	ProfilerEnable bool          `env:"DATABASE_PROFILER_ENABLE"`
	ProfilerOutput string        `env:"DATABASE_PROFILER_OUTPUT"`
	Verbose        bool          `env:"DATABASE_VERBOSE"`
	Debug          bool          `env:"DATABASE_DEBUG"`
	ErrorNoRows    bool          `env:"DATABASE_ERROR_NOROWS"`
}

// Boot load the default environment configuration.
func (e *Environment) Boot() {
	env.Unmarshal(e)
}

// Load custom environment variable based on given alias.
func (e *Environment) Load(alias string) {
	if e.Alias = strings.ToUpper(strings.TrimSpace(alias)); e.Alias == "" {
		return
	}

	for _, key := range []string{"DSN", "DRIVER", "PROTOCOL", "HOST", "PORT", "USER", "PASS", "CHARSET", "SCHEMA", "MODE", "AUTOCONNECT", "MAXOPEN", "MAXIDLE", "MAXLIFETIME", "PARSETIME", "ERROR_NOROWS"} {
		if v, ok := env.Lookup(fmt.Sprintf("DATABASE_%s_%s", e.Alias, key)); ok {
			switch key {
			case "DSN":
				e.DSN = v
			case "DRIVER":
				e.Driver = v
			case "PROTOCOL":
				e.Protocol = v
			case "HOST":
				e.Host = v
			case "USER":
				e.User = v
			case "PASS":
				e.Pass = v
			case "CHARSET":
				e.Charset = v
			case "SCHEMA":
				e.Schema = v
			case "MODE":
				e.Mode = v
			case "PORT":
				e.Port = v
			case "AUTOCONNECT":
				e.Autoconnect = toBool(v)
			case "PARSETIME":
				e.ParseTime = toBool(v)
			case "MAXOPEN":
				e.MaxOpen = toInt(v)
			case "MAXIDLE":
				e.MaxIdle = toInt(v)
			case "MAXLIFETIME":
				e.MaxLifetime = toDuration(v)
			case "VERBOSE":
				e.Verbose = toBool(v)
			case "DEBUG":
				e.Debug = toBool(v)
			case "ERROR_NOROWS":
				e.ErrorNoRows = toBool(v)
			}
		}
	}
}

// Validate the environment configuration.
func (e *Environment) Validate() error {
	// Managing required values.
	if e.DSN == "" {
		if e.User == "" {
			return errors.New("you must provide Environment.User")
		}

		if e.Pass == "" {
			return errors.New("you must provide Environment.Pass ; connections without password are not allowed by the package")
		}
	}

	// Managing default values.
	if e.Driver == "" {
		e.Driver = defaultDriver
	}

	if e.Port == "" {
		e.Port = defaultPort
	}

	if e.Protocol == "" {
		e.Protocol = defaultProtocol
	}

	if e.Charset == "" {
		e.Charset = defaultCharset
	}

	if e.Host == "" {
		e.Host = defaultHost
	}

	if e.Schema == "" {
		e.Schema = e.User
	}

	if e.MaxIdle == 0 {
		e.MaxIdle = defaultMaxIdle
	}

	if e.MaxOpen == 0 {
		e.MaxOpen = defaultMaxOpen
	}

	if e.MaxLifetime <= time.Second*0 {
		e.MaxLifetime = defaultMaxLifetime
	}

	return nil
}

// String use the existing source or write a new one based on inputs.
func (e *Environment) String() string {
	if e.DSN != "" {
		return e.DSN
	}

	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=%s",
		e.User,
		e.Pass,
		e.Protocol,
		e.Host,
		e.Port,
		e.Schema,
		e.Charset,
	)

	if e.Driver == "mysql" && e.ParseTime {
		dsn += "&parseTime=true"
	}
	return dsn
}
