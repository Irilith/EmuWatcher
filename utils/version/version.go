package version

// Version explain
// 0.0.0
// first 0 is major version
// second 0 is minor version
// third 0 is patch version
// major.minor.patch
var (
	version = "DEV"
	commit  = "DEV"
)

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}
