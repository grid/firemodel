package version

const dev = "(dev)"

var (
	Version string = dev
)

func IsDev() bool {
	return Version == dev
}
