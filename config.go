package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type LanguageConfig struct {
	Interpret string `yaml:",omitempty"`
	Compile   string `yaml:",omitempty"`
}

type Config struct {
	SolutionDir     string `yaml:"problem_dir"`
	TestDir         string `yaml:"test_dir"`
	DefaultLanguage string `yaml:"default_language"`
	Editor          string
	Languages       map[string]LanguageConfig
	IgnorePatterns  []string `yaml:"ignore"`
}

func baseDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(user.HomeDir, ".online-judge-helper")
}

func loadConfig() (*Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(filepath.Join(baseDir(), "config.yml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if config.SolutionDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		config.SolutionDir = pwd
	}
	if config.TestDir == "" {
		config.TestDir = filepath.Join(baseDir(), "tests")
	}
	return &config, err
}
