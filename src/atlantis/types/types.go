/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
