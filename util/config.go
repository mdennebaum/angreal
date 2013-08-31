package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	file string
	DynMap
}

func NewConfig(file string) *Config {
	c := Config{file, DynMap{make(map[string]interface{})}}
	return &c
}

func (this *Config) Load() *Config {
	cf, e := ioutil.ReadFile(this.file)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var f interface{}
	json.Unmarshal(cf, &f)
	this.Map = f.(map[string]interface{})
	return this
}
