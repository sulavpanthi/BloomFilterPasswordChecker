package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once         sync.Once
	globalConfig *Config
)

func GetConfig() *Config {
	return globalConfig
}

func Initialize() error {
	var err error
	once.Do(func() {
		globalConfig, err = LoadConfig()
	})
	return err
}

type Config struct {
	ExpectedElements         uint64
	FalsePositiveProbability float64
	BloomFilterFileName      string
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var errMsgs []string
	for _, err := range e {
		errMsgs = append(errMsgs, err.Error())
	}
	return fmt.Sprintf("validation failed:\n- %s", strings.Join(errMsgs, "\n- "))
}

func (c *Config) Validate() error {
	var errors ValidationErrors

	// config validations
	if err := c.validateExpectedElements(); err != nil {
		errors = append(errors, err...)
	}

	if err := c.validateFalsePositiveProbability(); err != nil {
		errors = append(errors, err...)
	}

	if len(errors) > 0 {
		return errors
	}
	return nil

}

func (c *Config) validateExpectedElements() ValidationErrors {
	var errors ValidationErrors

	if c.ExpectedElements < 1 {
		errors = append(errors, ValidationError{
			Field:   "ExpectedElements",
			Message: "must be greater than 1",
		})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (c *Config) validateFalsePositiveProbability() ValidationErrors {
	var errors ValidationErrors

	if c.FalsePositiveProbability > 1 && c.FalsePositiveProbability < 0 {
		errors = append(errors, ValidationError{
			Field:   "FalsePositiveProbability",
			Message: "must be between 0 and 1",
		})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func loadEnvVars() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err.Error())
	}

	config := &Config{}

	// Server configs
	config.ExpectedElements = getEnvAsUInt64("EXPECTED_ELEMENTS", 100)
	config.FalsePositiveProbability = getEnvAsFloat64("FALSE_POSITIVE_PROBABILITY", 0.05)
	config.BloomFilterFileName = getEnvRequired("BLOOM_FILTER_FILE_NAME")

	return config, nil
}

func LoadConfig() (*Config, error) {
	config, err := loadEnvVars()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err.Error())
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

	return config, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

func getEnvAsUInt64(key string, defaultVal uint64) uint64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseUint(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsFloat64(key string, defaultVal float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}
