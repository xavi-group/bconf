package bconf

const (
	loadOptionTypeHandleHelpFlag     = "handle_help_flag"
	loadOptionTypeHandleGenerateFlag = "handle_generate_flag"
)

type LoadOption interface {
	LoadOptionType() string
}

func WithHelpFlagHandler() LoadOption {
	return loadOptionHandleHelpFlag{}
}

func WithGenerateFlagHandler() LoadOption {
	return loadOptionHandleGenerateFlag{}
}

type loadOptionHandleHelpFlag struct{}

func (o loadOptionHandleHelpFlag) LoadOptionType() string {
	return loadOptionTypeHandleHelpFlag
}

type loadOptionHandleGenerateFlag struct{}

func (o loadOptionHandleGenerateFlag) LoadOptionType() string {
	return loadOptionTypeHandleGenerateFlag
}
