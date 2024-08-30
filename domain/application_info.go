package domain

import "runtime"

// ApplicationInfo contains basic information about an application
type ApplicationInfo struct {
	AppGroup  string
	App       string
	Version   string
	Branch    string
	Revision  string
	BuildDate string
}

// OperatingSystem returns the operating system the application is running on
func (ai *ApplicationInfo) OperatingSystem() string {
	return runtime.GOOS
}

// GoVersion  returns the version of Go the application is built with
func (ai *ApplicationInfo) GoVersion() string {
	return runtime.Version()
}

// NewApplicationInfo returns a new ApplicationInfo instance
func NewApplicationInfo(appgroup, app, version, branch, revision, buildDate string) ApplicationInfo {
	return ApplicationInfo{
		AppGroup:  appgroup,
		App:       app,
		Version:   version,
		Branch:    branch,
		Revision:  revision,
		BuildDate: buildDate,
	}
}
