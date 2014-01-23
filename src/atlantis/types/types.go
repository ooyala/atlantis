package types

import (
	"encoding/json"
	"os"
)

const (
	ContainerConfigDir  = "/etc/atlantis/config"
	ContainerConfigFile = "/etc/atlantis/config/config.json"
)

type ContainerConfig struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Env  string `json:"env"`
}

type AppConfig struct {
	HTTPPort       uint16                            `json:"http_port"`
	SecondaryPorts []uint16                          `json:"secondary_ports"`
	Container      *ContainerConfig                  `json:"container"`
	Dependencies   map[string]map[string]interface{} `json:"dependencies"`
}

func LoadAppConfig() (*AppConfig, error) {
	f, err := os.Open(ContainerConfigFile)
	if err != nil {
		return nil, err
	}
	cfg := AppConfig{}
	dec := json.NewDecoder(f)
	return &cfg, dec.Decode(&cfg)
}

func (a *AppConfig) Save(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	return enc.Encode(a)
}
