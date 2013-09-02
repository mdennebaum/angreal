package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
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
	var f interface{}
	e := json.Unmarshal(cf, &f)
	if e == nil {
		this.Map = f.(map[string]interface{})
	} else {
		fmt.Printf("Failed to parse config: %v\n", err)
		os.Exit(1)
	}
}
