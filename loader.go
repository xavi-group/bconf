package bconf

// Loader is the interface that configuration loaders must implement.
type Loader interface {
	CloneLoader() Loader
	Name() string
	Get(fieldSetKey, fieldKey string) (value any, found bool)
	GetMap(fieldSetKey string, fieldKeys []string) (fieldValues map[string]any)
	HelpString(fieldSetKey, fieldKey string) string
}

// LoaderKeyOverride allows overriding the key used by a specific loader for a field.
type LoaderKeyOverride struct {
	LoaderName     string
	KeyOverride    string
	IgnorePrefixes bool
}
