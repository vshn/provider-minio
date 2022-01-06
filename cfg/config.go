package cfg

var (
	// Depending on how isolated the CLI subcommands are, you can also maintain a Config per subcommand and put in another package.

	// Config contains the values of the user-provided configuration combined with the default values.
	Config = &Configuration{}
)

// EnvPrefix is the global prefix to use for the keys in environment variables
// TODO: Adjust or clear env var prefix
const EnvPrefix = "BOOTSTRAP"

// Configuration holds a strongly-typed tree of the configuration.
type Configuration struct {
	// ExampleFlag demonstrates how the configuration can be used in the business logic.
	ExampleFlag string
}

// Env combines EnvPrefix with given suffix delimited by underscore.
func Env(suffix string) string {
	return EnvPrefix + "_" + suffix
}
