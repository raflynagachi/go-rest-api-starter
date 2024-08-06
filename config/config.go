package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	App       App                  `json:"app"`
	Databases map[string]*Database `json:"databases"`
	JwtKey    string               `json:"jwt_key"`
	Redis     Redis                `json:"redis"`
}

var (
	readJsonConfig = ReadJsonConfig

	Env string
)

func init() {
	Env = os.Getenv(EnvKey)
	if Env == "" {
		Env = Development // default environment
	}
}

// LoadConfig load configuration based on env
func LoadConfig() (*Config, error) {
	path := fmt.Sprintf("env/%s.%s.json", ServiceName, Env)

	config, err := readJsonConfig(path)
	if err != nil {
		return nil, errors.Wrap(err, "LoadConfig.readJsonConfig")
	}

	return config, nil
}

// ReadJsonConfig mapping the json file to Config struct
func ReadJsonConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
