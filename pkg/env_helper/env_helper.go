package env_helper

import (
	"os"
	"strconv"
	"strings"
)

func GetStringEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func GetIntEnv(key string, defaultValue int) int {
	value := defaultValue
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		value = v
	}
	return value
}

func GetBoolEnv(key string, defaultValue bool) bool {
	envVal := strings.ToLower(os.Getenv(key))
	if defaultValue {
		return envVal != "false"
	} else {
		return envVal == "true"
	}
}
