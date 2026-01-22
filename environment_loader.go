package bconf

import (
	"fmt"
	"os"
	"strings"
)

// NewEnvironmentLoader creates a new environment loader without a key prefix.
func NewEnvironmentLoader() *EnvironmentLoader {
	return NewEnvironmentLoaderWithKeyPrefix("")
}

// NewEnvironmentLoaderWithKeyPrefix creates a new environment loader with the specified key prefix.
func NewEnvironmentLoaderWithKeyPrefix(keyPrefix string) *EnvironmentLoader {
	return &EnvironmentLoader{KeyPrefix: keyPrefix}
}

// EnvironmentLoader loads configuration values from environment variables.
type EnvironmentLoader struct {
	KeyPrefix string
}

// Clone creates a copy of the EnvironmentLoader.
func (l *EnvironmentLoader) Clone() *EnvironmentLoader {
	newLoader := *l
	return &newLoader
}

// CloneLoader creates a copy of the loader as a Loader interface.
func (l *EnvironmentLoader) CloneLoader() Loader {
	return l.Clone()
}

// Name returns the name of this loader.
func (l *EnvironmentLoader) Name() string {
	return "bconf_environment"
}

// Get retrieves a single field value from environment variables.
func (l *EnvironmentLoader) Get(fieldSetKey, fieldKey string) (any, bool) {
	return os.LookupEnv(l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
}

// GetMap retrieves multiple field values from environment variables.
func (l *EnvironmentLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]any {
	values := map[string]any{}

	for _, fieldKey := range fieldKeys {
		value, found := os.LookupEnv(l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
		if found {
			values[fieldKey] = value
		}
	}

	return values
}

// HelpString returns a help string describing where this field can be configured.
func (l *EnvironmentLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("Environment key: '%s'", l.environmentKey(fmt.Sprintf("%s_%s", fieldSetKey, fieldKey)))
}

func (l *EnvironmentLoader) environmentKey(key string) string {
	envKey := ""
	if l.KeyPrefix != "" {
		envKey = fmt.Sprintf("%s_%s", l.KeyPrefix, key)
	} else {
		envKey = key
	}

	return strings.ToUpper(envKey)
}
