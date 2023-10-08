package app

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
