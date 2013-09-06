package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	filePath string
	DynMap
}

func NewConfig(filePath string) *Config {
	c := Config{filePath, DynMap{}}
	return &c
}

func (this *Config) Load() {
	cf, err := ioutil.ReadFile(this.filePath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(cf, &this.Map)
	if this.Map == nil {
		fmt.Printf("Failed to parse config\n")
		os.Exit(1)
	}
}
