package config

import (
	"os"
	"strconv"
)

type Config struct {
	WEBSERVER_PORT  int
	SHARED_LIB_PATH string
	MODEL_PATH_ONNX string
}

func NewConfig() *Config {
	return &Config{
		WEBSERVER_PORT:  parseEnvToInt("WEBSERVER_PORT", "5000"),
		SHARED_LIB_PATH: getEnv("SHARED_LIB_PATH", ""),
		MODEL_PATH_ONNX: getEnv("MODEL_PATH_ONNX", ""),
	}
}

func getEnv(env, defaultValue string) string {
	enviroment := os.Getenv(env)
	if enviroment == "" {
		return defaultValue
	}

	return enviroment
}

func parseEnvToInt(envName, defaultValue string) int {
	num, err := strconv.Atoi(getEnv(envName, defaultValue))
	if err != nil {
		return 0
	}
	return num
}
