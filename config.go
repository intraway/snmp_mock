package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	SnmpPort      uint16 `yaml:"snmp_port"`
	SnmpCommunity string `yaml:"snmp_community"`
	BaseOid       string `yaml:"base_oid"`
	AppPort       uint16 `yaml:"app_port"`
}

func LoadConfig(path string) (Config, error) {
	yamlstr, err := ioutil.ReadFile(path)
	// Initialize with default values
	yamlparsed := Config{
		SnmpPort:      161,
		SnmpCommunity: "public",
		AppPort:       8080,
		BaseOid:       "1.3.6",
	}
	if err != nil {
		return yamlparsed, err
	}

	err = yaml.Unmarshal(yamlstr, &yamlparsed)
	if err != nil {
		return yamlparsed, err
	}

	return yamlparsed, nil
}
