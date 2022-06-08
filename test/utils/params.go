package utils

type TestParams struct {
	Snap string
	App  string
	TestConfigs
	TestNetworking
	TestVersion
}

type TestConfigs struct {
	TestEnvConfig      bool
	TestAppConfig      bool
	TestGlobalConfig   bool
	TestMixedConfig    bool
	DefaultServicePort string
}

type TestNetworking struct {
	TestOpenPorts        []string
	TestBindAddrLoopback bool
}

type TestVersion struct {
	TestSemanticSnapVersion bool
}
