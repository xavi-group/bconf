package bconf

const (
	loadOptionTypeDisableHelpFlag     = "disable_help_flag_handler"
	loadOptionTypeDisableGenerateFlag = "disable_generate_flag_handler"
)

// LoadOption is the interface for options passed to the Load method.
type LoadOption interface {
	LoadOptionType() string
}

// DisableHelpFlagHandler returns a LoadOption that disables the built-in help flag handler.
func DisableHelpFlagHandler() LoadOption {
	return loadOptionDisableHelpFlag{}
}

// DisableGenerateFlagHandler returns a LoadOption that disables the built-in generate flag handler.
func DisableGenerateFlagHandler() LoadOption {
	return loadOptionDisableGenerateFlag{}
}

type loadOptionDisableHelpFlag struct{}

func (o loadOptionDisableHelpFlag) LoadOptionType() string {
	return loadOptionTypeDisableHelpFlag
}

type loadOptionDisableGenerateFlag struct{}

func (o loadOptionDisableGenerateFlag) LoadOptionType() string {
	return loadOptionTypeDisableGenerateFlag
}
