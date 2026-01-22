package bconfconst

// Field type constants for configuration field definitions.
const (
	Bool      = "bool"
	Bools     = "[]bool"
	String    = "string"
	Strings   = "[]string"
	Int       = "int"
	Ints      = "[]int"
	Float     = "float64"
	Floats    = "[]float64"
	Time      = "time.Time"
	Times     = "[]time.Time"
	Duration  = "time.Duration"
	Durations = "[]time.Duration"
)

// FieldTypes returns all supported field type constants.
func FieldTypes() []string {
	return []string{
		Bool,
		Bools,
		String,
		Strings,
		Int,
		Ints,
		Float,
		Floats,
		Time,
		Times,
		Duration,
		Durations,
	}
}
