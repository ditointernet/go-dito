package env

import (
	"fmt"
	"os"
	"strconv"
)

// GetString extracts a String value from the given environment variable
func GetString(name string, defaultValue ...string) string {
	value := os.Getenv(name)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return value
}

// MustGetString extracts a String value from the given environment variable
// It exits the application if not present
func MustGetString(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fmt.Printf("%s can't be empty\n", name)
		os.Exit(1)
	}
	return value
}

// GetInt extracts an Int value from the given environment variable
func GetInt(name string, defaultValue ...int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return value
}

// MustGetInt extracts an Int value from the given environment variable
// It exits the application if not present
func MustGetInt(name string) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		fmt.Printf("%s must contain a float value!\n", name)
		os.Exit(1)
	}
	return value
}

// GetFloat extracts a Float value from the given environment variable
func GetFloat(name string, defaultValue ...float64) float64 {
	value, err := strconv.ParseFloat(os.Getenv(name), 64)
	if err != nil && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return value
}

// MustGetFloat extracts a Float value from the given environment variable
// It exits the application if not present
func MustGetFloat(name string) float64 {
	value, err := strconv.ParseFloat(os.Getenv(name), 64)
	if err != nil {
		fmt.Printf("%s must contain a float value!\n", name)
		os.Exit(1)
	}
	return value
}

// GetBool extracts a Bool value from the given environment variable
func GetBool(name string, defaultValue ...bool) bool {
	value, err := strconv.ParseBool(os.Getenv(name))
	if err != nil && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return value
}

// MustGetBool extracts a Bool value from the given environment variable
// It exits the application if not present
func MustGetBool(name string) bool {
	value, err := strconv.ParseBool(os.Getenv(name))
	if err != nil {
		fmt.Printf("%s must contain a boolean value! (true or false)\n", name)
		os.Exit(1)
	}
	return value
}
