package main

type config struct {
	API apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port string `yaml:"port"`
}
