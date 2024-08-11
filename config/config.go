package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-yaml/yaml"
)

var DefaultConfigPath string = "./config.yaml"

type Configuration struct {
	VersionNumber        string `yaml:"version"`
	ScreenWidth          int    `yaml:"screenWidth"`
	ScreenHeight         int    `yaml:"screenHeight"`
	MaxTPS               int    `yaml:"maxTPS"`
	MaxParticles         int    `yaml:"maxParticles"`
	RainAmount           int    `yaml:"rainAmount"`
	ShowSkippedParticles bool   `yaml:"showSkippedParticles"`
	VerboseLogging       bool   `yaml:"verboseLogging"`
}

func (conf *Configuration) ReadConfig() error {

	// Read bytes in file
	configBytes, err := os.ReadFile(DefaultConfigPath)

	// Check for error after reading bytes
	if err != nil {
		// Ship it back to main to be handled
		return err
	}

	err = yaml.Unmarshal(configBytes, &conf)

	if err != nil {
		return err
	}
	return nil // config was read successfully
}

func (conf *Configuration) LogConfig() {
	val := reflect.TypeOf(conf).Elem() // dereference pointer
	typ := reflect.TypeOf(conf).Elem() // dereference pointer

	log.Println("============== Config ==============")
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldName := typ.Field(i).Name
		log.Printf("%v: %v\n", fieldName, fieldValue)
	}
	log.Println("====================================")
}
