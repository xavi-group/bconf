package bconf

const (
	loadOptionTypeHandleHelpFlag     = "handle_help_flag"
	loadOptionTypeHandleGenerateFlag = "handle_generate_flag"
)

type LoadOption interface {
	OptionType() string
}

func WithHelpFlagHandler() LoadOption {
	return &loadOptionHandleHelpFlag{}
}

func WithGenerateFlagHandler() LoadOption {
	return &loadOptionHandleGenerateFlag{}
}

type loadOptionHandleHelpFlag struct{}

func (o loadOptionHandleHelpFlag) OptionType() string {
	return loadOptionTypeHandleHelpFlag
}

type loadOptionHandleGenerateFlag struct{}

func (o loadOptionHandleGenerateFlag) OptionType() string {
	return loadOptionTypeHandleGenerateFlag
}
