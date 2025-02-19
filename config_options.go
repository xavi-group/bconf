package bconf

const (
	configOptionTypeEnvironmentLoader = "loader_environment"
	configOptionTypeFlagLoader        = "loader_flag"
	configOptionTypeJSONFileLoader    = "loader_json"
)

type ConfigOption interface {
	OptionType() string
}

func WithEnvironmentLoader(keyPrefix string) ConfigOption {
	return &configOptionEnvironmentLoader{
		keyPrefix: keyPrefix,
	}
}

func WithFlagLoader(keyPrefix string) ConfigOption {
	return &configOptionFlagLoader{
		keyPrefix: keyPrefix,
	}
}

func WithJSONFileLoader(decoder JSONUnmarshal, filePaths ...string) ConfigOption {
	return &configOptionJSONFileLoader{}
}

type configOptionEnvironmentLoader struct {
	keyPrefix string
}

func (o configOptionEnvironmentLoader) OptionType() string {
	return configOptionTypeEnvironmentLoader
}

func (o configOptionEnvironmentLoader) Loader() Loader {
	return NewEnvironmentLoaderWithKeyPrefix(o.keyPrefix)
}

type configOptionFlagLoader struct {
	keyPrefix string
}

func (o configOptionFlagLoader) OptionType() string {
	return configOptionTypeFlagLoader
}

func (o configOptionFlagLoader) Loader() Loader {
	return NewFlagLoaderWithKeyPrefix(o.keyPrefix)
}

type configOptionJSONFileLoader struct {
	decoder   JSONUnmarshal
	filePaths []string
}

func (o configOptionJSONFileLoader) OptionType() string {
	return configOptionTypeJSONFileLoader
}

func (o configOptionJSONFileLoader) Loader() Loader {
	return NewJSONFileLoaderWithAttributes(o.decoder, o.filePaths...)
}
