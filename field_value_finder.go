package bconf

import "time"

// FieldValueFinder provides type-safe access to field values during load condition evaluation.
type FieldValueFinder interface {
	GetFieldDependencies() map[FieldLocation]any
	GetFieldValue(fieldSetKey, fieldKey string) (value any, found bool)
	GetString(fieldSetKey, fieldKey string) (val string, found bool, err error)
	GetStrings(fieldSetKey, fieldKey string) (val []string, found bool, err error)
	GetInt(fieldSetKey, fieldKey string) (val int, found bool, err error)
	GetInts(fieldSetKey, fieldKey string) (val []int, found bool, err error)
	GetBool(fieldSetKey, fieldKey string) (val bool, found bool, err error)
	GetBools(fieldSetKey, fieldKey string) (val []bool, found bool, err error)
	GetTime(fieldSetKey, fieldKey string) (val time.Time, found bool, err error)
	GetTimes(fieldSetKey, fieldKey string) (val []time.Time, found bool, err error)
	GetDuration(fieldSetKey, fieldKey string) (val time.Duration, found bool, err error)
	GetDurations(fieldSetKey, fieldKey string) (val []time.Duration, found bool, err error)
}

// FieldValue represents a configuration field value with its location.
type FieldValue struct {
	FieldValue  any
	FieldSetKey string
	FieldKey    string
}

// FieldValues is a slice of FieldValue.
type FieldValues []FieldValue
