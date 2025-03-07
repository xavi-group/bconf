package bconf

import "time"

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

type FieldValue struct {
	FieldValue  any
	FieldSetKey string
	FieldKey    string
}

type FieldValues []FieldValue
