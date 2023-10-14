package app

import (
	"fmt"
	"time"
)

// Info describes the basic information of the application.
type Info struct {
	// Name is the name of the application.
	Name string
	// BuildTime is the time when the application was built.
	BuildTime string
	// BuildVersion is the version of the application when it was built.
	// If tag is not available, it will be the commit hash.
	BuildVersion string
}

// NewInfo creates a new Info and validates the build time.
func NewInfo(name, buildTime, buildVersion string) (Info, error) {
	buildAt, err := time.Parse(time.RFC3339, buildTime)
	if err != nil {
		return Info{}, fmt.Errorf("invalid build time: %w", err)
	}

	buildTime = buildAt.In(time.Local).String()

	info := Info{
		Name:         name,
		BuildTime:    buildTime,
		BuildVersion: buildVersion,
	}

	return info, nil
}

// String returns the string representation of the App Info.
func (i Info) String() string {
	return fmt.Sprintf("%s:%s (%s)", i.Name, i.BuildVersion, i.BuildTime)
}
