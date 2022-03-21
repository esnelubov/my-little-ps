package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type IConfig interface {
	GetString(key string) string
	GetBool(key string) bool
	GetDurationSec(key string) time.Duration
}

type Config struct {
	vprConfig *viper.Viper
}

func New(name string) *Config {
	vprConfig := viper.New()

	vprConfig.SetConfigName(name)
	vprConfig.SetConfigType("yaml")
	vprConfig.AddConfigPath(".")
	vprConfig.AddConfigPath("./..")
	vprConfig.AddConfigPath("~/")

	vprConfig.SetDefault("ip", "127.0.0.1:3000")
	vprConfig.SetDefault("dsn", "root:root@tcp(127.0.0.1:3306)/dropix_mobile_refill?charset=utf8mb4&parseTime=True")
	vprConfig.SetDefault("idleTimeoutSec", 15)
	vprConfig.SetDefault("readTimeoutSec", 15)
	vprConfig.SetDefault("writeTimeoutSec", 15)
	vprConfig.SetDefault("shutdownTimeoutSec", 30)
	vprConfig.SetDefault("location", "Europe/Moscow")

	if err := vprConfig.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &Config{
		vprConfig: vprConfig,
	}
}

func (c *Config) GetString(key string) string {
	return c.vprConfig.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.vprConfig.GetBool(key)
}

func (c *Config) GetDurationSec(key string) time.Duration {
	return time.Duration(c.vprConfig.GetInt(key)) * time.Second
}
