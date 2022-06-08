package utils

type Params struct {
	Snap string
	App  string
	ConfigTests
	NetworkingTests
	PackagingTests
}

type ConfigTests struct {
	// TODO: pass port to tests and refactor to allow testing other config options
	DefaultServicePort string // used by config tests
	TestEnvConfig      bool
	TestAppConfig      bool
	TestGlobalConfig   bool
	TestMixedConfig    bool
}

type NetworkingTests struct {
	TestOpenPorts        []string
	TestBindAddrLoopback bool
}

type PackagingTests struct {
	TestSemanticSnapVersion bool
}
