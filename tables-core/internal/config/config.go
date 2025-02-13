package config

import (
    "log"
    "os"
    "path"
    "tables/pkg/config"
)

type Config struct {
    Server      *ServerConfig `json:"server"`
    Gulmarketdb *Database     `json:"gulmarketdb"`
    Meter       *MeterConfig  `json:"metric"` // meter - to expose metrics
}

func NewAppConfig() (*Config, error) {
    cfg := &Config{
        Server:      &ServerConfig{},
        Gulmarketdb: &Database{},
        Meter:       &MeterConfig{},
    }

    pathPrefix := os.Getenv("CONFIG_PATH")
    wd, wdErr := os.Getwd()
    if wdErr != nil {
        log.Fatal(wdErr)
    }
    if len(pathPrefix) == 0 {
        pathPrefix = path.Join(wd, "/helm/local/")
    }
    configPath := path.Join(pathPrefix, "config.json")
    log.Printf("CONFIG_PATH: | %s |", configPath)
    if err := config.ParseFileConfig(configPath, cfg); err != nil {
        log.Fatalf("Failed to parse file to config. Error: %v\n", err)
    }
    return cfg, nil
}

type ServerConfig struct {
    Addr    string `json:"addr" env:"HTTP_PORT" default:":9090"`
    Name    string `json:"name"`
    Version string `json:"version" env:"VERSION"`
    Proxy   string `json:"proxy"`
    Timeout string `json:"timeout"`
}

type Database struct {
    Url string `json:"url"`
}

type MeterConfig struct {
    Addr string `json:"addr" default:"0.0.0.0:8080"`
}
