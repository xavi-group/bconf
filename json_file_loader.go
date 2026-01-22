package bconf

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

// JSONUnmarshal defines the function signature for JSON unmarshaling.
type JSONUnmarshal func(data []byte, v any) error

// NewJSONFileLoader creates a new JSON file loader with default settings.
func NewJSONFileLoader() *JSONFileLoader {
	return NewJSONFileLoaderWithAttributes(nil)
}

// NewJSONFileLoaderWithAttributes creates a new JSON file loader with the specified decoder and file paths.
func NewJSONFileLoaderWithAttributes(decoder JSONUnmarshal, filePaths ...string) *JSONFileLoader {
	if decoder == nil {
		decoder = json.Unmarshal
	}

	return &JSONFileLoader{
		Decoder:   decoder,
		FilePaths: filePaths,
	}
}

// JSONFileLoader loads configuration values from JSON files.
type JSONFileLoader struct {
	Decoder   JSONUnmarshal
	FilePaths []string
	fileMaps  []map[string]any
}

// Clone creates a copy of the JSONFileLoader.
func (l *JSONFileLoader) Clone() *JSONFileLoader {
	clone := *l

	clone.FilePaths = slices.Clone(l.FilePaths)
	clone.fileMaps = nil

	return &clone
}

// CloneLoader creates a copy of the loader as a Loader interface.
func (l *JSONFileLoader) CloneLoader() Loader {
	return l.Clone()
}

// Name returns the name of this loader.
func (l *JSONFileLoader) Name() string {
	return "bconf_jsonfile"
}

// Get retrieves a single field value from the loaded JSON files.
func (l *JSONFileLoader) Get(fieldSetKey, fieldKey string) (any, bool) {
	fileMaps := l.getFileMaps()

	if len(fileMaps) < 1 {
		return "", false
	}

	return l.findValueInMaps(fieldSetKey, fieldKey, fileMaps)
}

// GetMap retrieves multiple field values from the loaded JSON files.
func (l *JSONFileLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]any {
	values := map[string]any{}

	fileMaps := l.getFileMaps()

	if len(fileMaps) < 1 {
		return values
	}

	for _, fieldKey := range fieldKeys {
		val, found := l.findValueInMaps(fieldSetKey, fieldKey, fileMaps)
		if found {
			values[fieldKey] = val
		}
	}

	return values
}

// HelpString returns a help string describing where this field can be configured.
func (l *JSONFileLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("JSON attribute: %s.%s", fieldSetKey, fieldKey)
}

func (l *JSONFileLoader) findValueInMaps(fieldSetKey, fieldKey string, fileMaps []map[string]any) (any, bool) {
	for _, fileMap := range fileMaps {
		fieldSetAny, found := fileMap[fieldSetKey]
		if !found {
			continue
		}

		fieldSetMap, ok := fieldSetAny.(map[string]any)
		if !ok {
			continue
		}

		value, ok := fieldSetMap[fieldKey]
		if !ok {
			continue
		}

		return l.convertValue(value), true
	}

	return nil, false
}

func (l *JSONFileLoader) convertValue(value any) any {
	switch v := value.(type) {
	case []any:
		parts := make([]string, len(v))
		for i, elem := range v {
			parts[i] = l.scalarToString(elem)
		}

		return strings.Join(parts, ",")
	case map[string]any:
		return l.convertToMapStringAny(v)
	default:
		return l.scalarToString(value)
	}
}

func (l *JSONFileLoader) convertToMapStringAny(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))
	maps.Copy(result, m)

	return result
}

func (l *JSONFileLoader) scalarToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}

		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (l *JSONFileLoader) getFileMaps() []map[string]any {
	if l.fileMaps != nil {
		return l.fileMaps
	}

	l.loadFileMaps()

	return l.fileMaps
}

func (l *JSONFileLoader) loadFileMaps() {
	l.fileMaps = []map[string]any{}

	for _, path := range l.FilePaths {
		fileBytes, err := os.ReadFile(path) //nolint:gosec // File paths are provided by the user
		if err != nil {
			continue
		}

		fileMap := map[string]any{}
		if err := l.Decoder(fileBytes, &fileMap); err != nil {
			continue
		}

		l.fileMaps = append(l.fileMaps, fileMap)
	}
}
