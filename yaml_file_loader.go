package bconf

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLUnmarshal func(data []byte, v any) error

type YAMLFileLoaderOption func(*YAMLFileLoader)

func WithYAMLDecoder(decoder YAMLUnmarshal) YAMLFileLoaderOption {
	return func(l *YAMLFileLoader) {
		l.Decoder = decoder
	}
}

func WithYAMLFilePaths(paths ...string) YAMLFileLoaderOption {
	return func(l *YAMLFileLoader) {
		l.FilePaths = append(l.FilePaths, paths...)
	}
}

func NewYAMLFileLoader(opts ...YAMLFileLoaderOption) *YAMLFileLoader {
	loader := &YAMLFileLoader{
		Decoder: yaml.Unmarshal,
	}

	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

type YAMLFileLoader struct {
	Decoder   YAMLUnmarshal
	FilePaths []string
	fileMaps  []map[string]any
}

func (l *YAMLFileLoader) Clone() *YAMLFileLoader {
	clone := *l

	clone.FilePaths = slices.Clone(l.FilePaths)
	clone.fileMaps = nil

	return &clone
}

func (l *YAMLFileLoader) CloneLoader() Loader {
	return l.Clone()
}

func (l *YAMLFileLoader) Name() string {
	return "bconf_yamlfile"
}

func (l *YAMLFileLoader) Get(fieldSetKey, fieldKey string) (any, bool) {
	maps := l.getFileMaps()

	if len(maps) < 1 {
		return "", false
	}

	return l.findValueInMaps(fieldSetKey, fieldKey, maps)
}

func (l *YAMLFileLoader) GetMap(fieldSetKey string, fieldKeys []string) map[string]any {
	values := map[string]any{}

	maps := l.getFileMaps()

	if len(maps) < 1 {
		return values
	}

	for _, fieldKey := range fieldKeys {
		val, found := l.findValueInMaps(fieldSetKey, fieldKey, maps)
		if found {
			values[fieldKey] = val
		}
	}

	return values
}

func (l *YAMLFileLoader) HelpString(fieldSetKey, fieldKey string) string {
	return fmt.Sprintf("YAML attribute: %s.%s", fieldSetKey, fieldKey)
}

func (l *YAMLFileLoader) findValueInMaps(fieldSetKey, fieldKey string, maps []map[string]any) (any, bool) {
	for _, fileMap := range maps {
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

func (l *YAMLFileLoader) convertValue(value any) any {
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

func (l *YAMLFileLoader) convertToMapStringAny(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))
	maps.Copy(result, m)

	return result
}

func (l *YAMLFileLoader) scalarToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
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

func (l *YAMLFileLoader) getFileMaps() []map[string]any {
	if l.fileMaps != nil {
		return l.fileMaps
	}

	l.loadFileMaps()

	return l.fileMaps
}

func (l *YAMLFileLoader) loadFileMaps() {
	l.fileMaps = []map[string]any{}

	for _, path := range l.FilePaths {
		fileBytes, err := os.ReadFile(path)
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
