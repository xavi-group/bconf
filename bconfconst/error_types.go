// Package bconfconst provides constants used by the bconf package.
package bconfconst

// Error message constants for field configuration validation.
const (
	ErrorFieldDefaultSetting      = "invalid settings: cannot set both Default and DefaultGenerator"
	ErrorFieldRequiredWithDefault = "invalid settings: cannot set both Required and Default/DefaultGenerator"
)
