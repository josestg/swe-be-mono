package config

import (
	"fmt"
	"os"
	"time"

	"github.com/josestg/swe-be-mono/pkg/env"
	"github.com/josestg/swe-be-mono/pkg/httpkit"
	"github.com/rs/cors"
)

// Config is a central configuration for the application.
type Config struct {
	AppInfo    AppInfo
	HttpCORS   cors.Options
	HttpServer httpkit.RunConfig
}

// New creates a new Config.
func New(appName, buildTime, buildVersion string) (*Config, error) {
	appInfo, err := NewAppInfo(appName, buildTime, buildVersion)
	if err != nil {
		return nil, fmt.Errorf("create app info: %w", err)
	}

	cfg := &Config{
		AppInfo: appInfo,
		HttpServer: httpkit.RunConfig{
			Port:                env.Int("HTTP_SERVER_PORT", 8080),
			ShutdownTimeout:     env.Duration("HTTP_SERVER_SHUTDOWN_TIMEOUT", 5*time.Second),
			RequestReadTimeout:  env.Duration("HTTP_SERVER_REQUEST_READ_TIMEOUT", 5*time.Second),
			RequestWriteTimeout: env.Duration("HTTP_SERVER_REQUEST_WRITE_TIMEOUT", 10*time.Second),
		},
		HttpCORS: cors.Options{
			AllowedOrigins:     env.StringList("HTTP_CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:     env.StringList("HTTP_CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH"}),
			AllowedHeaders:     env.StringList("HTTP_CORS_ALLOWED_HEADERS", []string{"*"}),
			AllowCredentials:   env.Bool("HTTP_CORS_ALLOW_CREDENTIALS", false),
			MaxAge:             env.Int("HTTP_CORS_MAX_AGE", 0),
			OptionsPassthrough: env.Bool("HTTP_CORS_OPTIONS_PASSTHROUGH", false),
			Debug:              env.Bool("HTTP_CORS_DEBUG", false),
		},
	}

	return cfg, nil
}

// AppInfo describes the basic information of the application.
type AppInfo struct {
	// Name is the name of the application.
	Name string
	// BuildTime is the time when the application was built.
	BuildTime string
	// BuildVersion is the version of the application when it was built.
	// If tag is not available, it will be the commit hash.
	BuildVersion string

	// Hostname is the hostname of the machine where the application is running.
	Hostname string
}

// NewAppInfo creates a new AppInfo and validates the build time.
func NewAppInfo(name, buildTime, buildVersion string) (AppInfo, error) {
	buildAt, err := time.Parse(time.RFC3339, buildTime)
	if err != nil {
		return AppInfo{}, fmt.Errorf("invalid build time: %w", err)
	}
	buildTime = buildAt.In(time.Local).String()

	hostname, err := os.Hostname()
	if err != nil {
		return AppInfo{}, fmt.Errorf("get hostname: %w", err)
	}

	info := AppInfo{
		Name:         name,
		Hostname:     hostname,
		BuildTime:    buildTime,
		BuildVersion: buildVersion,
	}

	return info, nil
}

// String returns the string representation of the App Info.
func (i AppInfo) String() string {
	return fmt.Sprintf("%s:%s (%s)", i.Name, i.BuildVersion, i.BuildTime)
}
