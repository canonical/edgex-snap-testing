package utils

type Params struct {
	Snap      string
	Config    Config
	Net       Net
	Packaging Packaging
}

type Config struct {
	TestChangePort ConfigChangePort
}

type ConfigChangePort struct {
	App                      string
	DefaultPort              string
	TestLegacyEnvConfig      bool
	TestAppConfig            bool
	TestGlobalConfig         bool
	TestMixedGlobalAppConfig bool
}

type Net struct {
	StartSnap        bool // should be set to true if services aren't started by default
	TestOpenPorts    []string
	TestBindLoopback []string
}

type Packaging struct {
	TestSemanticSnapVersion bool
}
