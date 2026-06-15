package config

// ServiceName is overridden by each cmd/*/main.go before calling LoadConfig.
var ServiceName = "go-rest-api-starter"

const (
	EnvKey      = "env"
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)
