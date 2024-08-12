package config

import (
	"log"
	"os"

	"git.smallzcomputing.com/sand-game/util"
	"github.com/go-yaml/yaml"
)

var DefaultConfigPath string = "./config.yaml"

type Configuration struct {
	VersionNumber        string       `yaml:"version"`
	ScreenSize           util.Vector2 `yaml:"screenSize"`
	ParticleColor        util.RGBA    `yaml:"particleColor"`
	BackgroundColor      util.RGBA    `yaml:"backgroundColor"`
	MaxTPS               int          `yaml:"maxTPS"`
	MaxParticles         int          `yaml:"maxParticles"`
	RainRate             int          `yaml:"rainRate"`
	ShowSkippedParticles bool         `yaml:"showSkippedParticles"`
	SkippedParticleColor util.RGBA    `yaml:"skippedParticleColor"`
	VerboseLogging       bool         `yaml:"verboseLogging"`
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

// Does not represent data captured after marshaling
func (conf *Configuration) LogConfig() {

	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("[ERROR] Failed to read config.yaml while logging.")
	}

	log.Println("============== Config ==============")
	log.Printf("%v\n", string(data))
	log.Println("====================================")
}
