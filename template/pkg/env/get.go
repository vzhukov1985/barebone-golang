package env

import (
	"github.com/{{.orgName}}/{{.pkgRepoName}}/log"
	"os"
	"strconv"
	"time"
)

func GetString(env string, defaultVal string, isRequired bool) string {
	valStr := os.Getenv(env)
	if valStr != "" {
		return valStr
	}

	if isRequired {
		log.Fatalw("Required env not specified", "env", env)
	}

	return defaultVal
}

func GetInt(env string, defaultVal int, isRequired bool) int {
	valStr := os.Getenv(env)
	if valStr != "" {
		retVal, err := strconv.Atoi(valStr)
		if err != nil {
			log.Fatalw("Invalid env value. Must be integer", "env", env)
		}
		return retVal
	}

	if isRequired {
		log.Fatalw("Required env not specified", "env", env)
	}

	return defaultVal
}

func GetFloat(env string, defaultVal float64, isRequired bool) float64 {
	valStr := os.Getenv(env)
	if valStr != "" {
		retVal, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			log.Fatalw("Invalid env value. Must be float", "env", env)
		}
		return retVal
	}

	if isRequired {
		log.Fatalw("Required env not specified", "env", env)
	}

	return defaultVal
}

func GetDuration(env string, defaultVal int, d time.Duration, isRequired bool) time.Duration {
	valStr := os.Getenv(env)
	if valStr != "" {
		retVal, err := strconv.Atoi(valStr)
		if err != nil {
			log.Fatalw("Invalid env value. Must be int", "env", env)
		}
		return time.Duration(retVal) * d
	}

	if isRequired {
		log.Fatalw("Required env not specified", "env", env)
	}

	return time.Duration(defaultVal) * d
}

func GetBool(env string, defaultVal bool, isRequired bool) bool {
	valStr := os.Getenv(env)
	if valStr != "" {
		retVal, err := strconv.ParseBool(valStr)
		if err != nil {
			log.Fatalw("Invalid env value. Must be bool", "env", env)
		}
		return retVal
	}

	if isRequired {
		log.Fatalw("Required env not specified", "env", env)
	}

	return defaultVal
}
