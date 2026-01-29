package meta

import "fmt"

var (
	Name       = "forklift"
	BuildDate  string
	CommitHash string
	Version    string
	Platform   string
	GoVersion  string
)

func VersionMultiline() string {
	return fmt.Sprintf(
		"%s %s (%s %s)\nBuild date: %s\nCommit hash: %s\n",
		Name,
		Version,
		GoVersion,
		Platform,
		BuildDate,
		CommitHash,
	)
}
