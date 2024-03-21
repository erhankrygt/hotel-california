package envvars

import (
	"fmt"
	"github.com/codingconcepts/env"
	"time"
)

// EnvVars represents environment variables
type EnvVars struct {
	Service      Service
	HTTPServer   HTTPServer
	MySql        MySql
	Localization Localization
	JWTToken     JWTToken
}

// Service represents service configurations
type Service struct {
	Environment string `env:"SERVICE_ENVIRONMENT" default:"dev"`
}

// HTTPServer represents http server configurations
type HTTPServer struct {
	Address         string        `env:"HTTP_SERVER_ADDRESS"`
	ReadTimeout     time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" default:"15s"`
	WriteTimeout    time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout     time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" default:"15s"`
	MaxHeaderBytes  int           `env:"HTTP_SERVER_MAX_HEADER_BYTES" default:"1048576"`
	ShutdownTimeout time.Duration `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"10s"`
}

// MySql represents mysql configurations
type MySql struct {
	URI            string `env:"MYSQL_URI" required:"true"`
	Database       string `env:"MYSQL_DATABASE" required:"true"`
	UserName       string `env:"MYSQL_USER_NAME" required:"true"`
	Port           string `env:"MYSQL_PORT" required:"true"`
	Password       string `env:"MYSQL_PASSWORD" required:"true"`
	ConnectTimeout int    `env:"MYSQL_CONNECT_TIMEOUT" default:"10"`
}

// Localization represents localization configurations
type Localization struct {
	LanguageFilesDirectory string `env:"LOCALIZATION_LANGUAGE_FILES_DIRECTORY" default:"internal/localization/language-files"`
}

// JWTToken represents jwt configurations
type JWTToken struct {
	Secret string `env:"JWT_TOKEN_SECRET" required:"true"`
}

// LoadEnvVars loads and returns environment variables
func LoadEnvVars() (*EnvVars, error) {
	s := Service{}
	if err := env.Set(&s); err != nil {
		return nil, fmt.Errorf("loading service environment variables failed, %s", err.Error())
	}

	hs := HTTPServer{}
	if err := env.Set(&hs); err != nil {
		return nil, fmt.Errorf("loading http server environment variables failed, %s", err.Error())
	}

	ms := MySql{}
	if err := env.Set(&ms); err != nil {
		return nil, fmt.Errorf("loading mysql environment variables failed, %s", err.Error())
	}

	l := Localization{}
	if err := env.Set(&l); err != nil {
		return nil, fmt.Errorf("loading localization environment variables failed, %s", err.Error())
	}

	jwt := JWTToken{}
	if err := env.Set(&jwt); err != nil {
		return nil, fmt.Errorf("loading jwt environment variables failed, %s", err.Error())
	}

	ev := &EnvVars{
		Service:      s,
		HTTPServer:   hs,
		MySql:        ms,
		Localization: l,
		JWTToken:     jwt,
	}

	return ev, nil
}
