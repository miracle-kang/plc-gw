package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/miracle-kang/plc-gw/internal/pkg"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// HTTP Server configuration
	Server ServerConfig `yaml:"server"`
	// Logging configuration
	Logging LoggingConfig `yaml:"logging"`
	// PLC Configuration
	PLC PLCConfig `yaml:"plc"`
}

type ServerConfig struct {
	// The host port to bind HTTP to
	Port int `yaml:"port"`
}

type LoggingConfig struct {
	Path    string `yaml:"path"`
	MaxSize string `yaml:"max-size"`
	MaxFile int    `yaml:"max-file"`
}

type PLCConfig struct {
	// The Data Base Path to store the PLC data in
	BasePath       string `yaml:"basepath"`
	CheckInterval  int    `yaml:"check-interval"`
	TimeoutSeconds int    `yaml:"timeout-seconds"`
}

var DefaultConfig = &Config{
	Server: ServerConfig{
		Port: 8080,
	},
	Logging: LoggingConfig{
		Path:    "./data/logs",
		MaxSize: "100M",
		MaxFile: 10,
	},
	PLC: PLCConfig{
		BasePath:       "./data/plc",
		CheckInterval:  1,
		TimeoutSeconds: 30,
	},
}

// LoadedConfig is the global configuration
var LoadedConfig *Config

// NewConfig returns a new decoded Config struct
func newConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func parseConfigFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Default config path from app
	dir, _ := pkg.ExeDir()
	dftPath := filepath.Join(dir, "app.yaml")
	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", dftPath, "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

// LoadConfig will load the configuration file and return
func LoadConfig() *Config {
	if LoadedConfig != nil {
		return LoadedConfig
	}
	// Load the configuration file
	path, err := parseConfigFlags()
	if path == "" || err != nil {
		log.Println("Could not parse flags: ", err)
		LoadedConfig = DefaultConfig
		log.Println("Using default configuration")
		return LoadedConfig
	}
	config, err := newConfig(path)
	if err != nil {
		log.Println("Could not load config: ", err)
		LoadedConfig = DefaultConfig
		log.Println("Using default configuration")
		return LoadedConfig
	}
	log.Printf("Using configuration from: %v\n", path)
	log.Printf("Loadded configuration: %v\n", config)
	LoadedConfig = config
	return LoadedConfig
}
