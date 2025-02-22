package bconf

const (
	loadOptionTypeDisableHelpFlag     = "disable_help_flag_handler"
	loadOptionTypeDisableGenerateFlag = "disable_generate_flag_handler"
)

type LoadOption interface {
	LoadOptionType() string
}

func DisableHelpFlagHandler() LoadOption {
	return loadOptionDisableHelpFlag{}
}

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
