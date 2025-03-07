package bconf

const (
	configOptionTypeLoaderEnvironment = "loader_environment"
	configOptionTypeLoaderFlag        = "loader_flag"
	configOptionTypeLoaderJSONFile    = "loader_json"
	configOptionTypeAppVersionFunc    = "app_version_func"
	configOptionTypeAppVersion        = "app_version"
	configOptionTypeAppIDFunc         = "app_id_func"
	configOptionTypeAppID             = "app_id"
)

type JSONLoaderConfigOption interface {
	ConfigOption
	WithDecoder(decoder JSONUnmarshal)
}

type ConfigOption interface {
	ConfigOptionType() string
}

// WithEnvironmentLoader enables the Environment loader. Only the first value in the keyPrefix parameter will be
// accepted as a key prefix.
func WithEnvironmentLoader(keyPrefix ...string) ConfigOption {
	prefix := ""
	if len(keyPrefix) > 0 {
		prefix = keyPrefix[0]
	}

	return configOptionEnvironmentLoader{
		keyPrefix: prefix,
	}
}

// WithFlagLoader enables the Flag loader. Only the first value in the keyPrefix parameter will be
// accepted as a key prefix.
func WithFlagLoader(keyPrefix ...string) ConfigOption {
	prefix := ""
	if len(keyPrefix) > 0 {
		prefix = keyPrefix[0]
	}

	return configOptionFlagLoader{
		keyPrefix: prefix,
	}
}

func WithJSONFileLoader(filePaths ...string) JSONLoaderConfigOption {
	return &configOptionJSONFileLoader{filePaths: filePaths}
}

func WithAppID(appID string) ConfigOption {
	return configOptionAppID{id: appID}
}

func WithAppIDFunc(appIDFunc func() string) ConfigOption {
	return configOptionAppIDFunc{idFunc: appIDFunc}
}

func WithAppVersion(appVersion string) ConfigOption {
	return configOptionAppVersion{version: appVersion}
}

func WithAppVersionFunc(appVersionFunc func() string) ConfigOption {
	return configOptionAppVersionFunc{versionFunc: appVersionFunc}
}

type configOptionEnvironmentLoader struct {
	keyPrefix string
}

func (o configOptionEnvironmentLoader) ConfigOptionType() string {
	return configOptionTypeLoaderEnvironment
}

func (o configOptionEnvironmentLoader) Loader() Loader {
	return NewEnvironmentLoaderWithKeyPrefix(o.keyPrefix)
}

type configOptionFlagLoader struct {
	keyPrefix string
}

func (o configOptionFlagLoader) ConfigOptionType() string {
	return configOptionTypeLoaderFlag
}

func (o configOptionFlagLoader) Loader() Loader {
	return NewFlagLoaderWithKeyPrefix(o.keyPrefix)
}

type configOptionJSONFileLoader struct {
	decoder   JSONUnmarshal
	filePaths []string
}

func (o *configOptionJSONFileLoader) WithDecoder(decoder JSONUnmarshal) {
	o.decoder = decoder
}

func (o *configOptionJSONFileLoader) ConfigOptionType() string {
	return configOptionTypeLoaderJSONFile
}

func (o configOptionJSONFileLoader) Loader() Loader {
	return NewJSONFileLoaderWithAttributes(o.decoder, o.filePaths...)
}

type configOptionAppVersion struct {
	version string
}

func (o configOptionAppVersion) ConfigOptionType() string {
	return configOptionTypeAppVersion
}

type configOptionAppVersionFunc struct {
	versionFunc func() string
}

func (o configOptionAppVersionFunc) ConfigOptionType() string {
	return configOptionTypeAppVersionFunc
}

type configOptionAppID struct {
	id string
}

func (o configOptionAppID) ConfigOptionType() string {
	return configOptionTypeAppID
}

type configOptionAppIDFunc struct {
	idFunc func() string
}

func (o configOptionAppIDFunc) ConfigOptionType() string {
	return configOptionTypeAppIDFunc
}
