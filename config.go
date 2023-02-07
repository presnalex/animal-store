package main

import "github.com/presnalex/go-micro/v3/service"

func newConfig(name, version string) *Config {
	return &Config{
		Core:            &service.CoreConfig{},
		Server:          &service.ServerConfig{Name: name, Version: version},
		PostgresPrimary: &service.PostgresConfig{},
		Metric:          &service.MetricConfig{},
		Consul:          &service.ConsulConfig{},
	}
}

type Config struct {
	Core            *service.CoreConfig     `json:"core"`
	Server          *service.ServerConfig   `json:"server"`
	Consul          *service.ConsulConfig   `json:"consul"`
	PostgresPrimary *service.PostgresConfig `json:"postgres_primary"`
	Metric          *service.MetricConfig   `json:"metric"`
}
