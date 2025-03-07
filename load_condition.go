package bconf

import (
	"fmt"
	"maps"
	"slices"
	"time"
)

type LoadConditions []LoadCondition

type LoadCondition interface {
	FieldValueFinder
	Clone() LoadCondition
	FieldDependencies() FieldLocations
	SetFieldValues(fieldValues ...FieldValue)
	Load(c FieldValueFinder) (bool, error)
}

func FD(fieldSetKey, fieldKey string) FieldLocation {
	return FieldLocation{
		FieldSetKey: fieldSetKey,
		FieldKey:    fieldKey,
	}
}

type FieldLocation struct {
	FieldSetKey string
	FieldKey    string
}

type FieldLocations []FieldLocation

// --------------------------------------------------------------------------------------------------------------------

func newLoadCondition(loadFunc func(f FieldValueFinder) (bool, error)) *loadCondition {
	return &loadCondition{
		loadFunc:              loadFunc,
		fieldDependencyValues: map[FieldLocation]any{},
		fieldDependencies:     FieldLocations{},
	}
}

// --------------------------------------------------------------------------------------------------------------------

type loadCondition struct {
	loadFunc              func(f FieldValueFinder) (bool, error)
	fieldDependencyValues map[FieldLocation]any
	fieldDependencies     FieldLocations
}

func (c *loadCondition) Clone() LoadCondition {
	clone := *c

	clone.fieldDependencies = slices.Clone(c.fieldDependencies)
	clone.fieldDependencyValues = maps.Clone(c.fieldDependencyValues)

	return &clone
}

func (c *loadCondition) FieldDependencies() FieldLocations {
	return slices.Clone(c.fieldDependencies)
}

func (c *loadCondition) SetFieldValues(fieldValues ...FieldValue) {
	for _, fieldValue := range fieldValues {
		key := FieldLocation{FieldSetKey: fieldValue.FieldSetKey, FieldKey: fieldValue.FieldKey}
		c.fieldDependencyValues[key] = fieldValue.FieldValue
	}
}

func (c *loadCondition) Load(loadConditionValues FieldValueFinder) (bool, error) {
	return c.loadFunc(loadConditionValues)
}

const FieldNotFoundError = "field not found"

// FieldValueFinder implementation

func (c *loadCondition) GetFieldDependencies() map[FieldLocation]any {
	return maps.Clone(c.fieldDependencyValues)
}

func (c *loadCondition) GetFieldValue(fieldSetKey, fieldKey string) (any, bool) {
	val, found := c.fieldDependencyValues[FieldLocation{FieldSetKey: fieldSetKey, FieldKey: fieldKey}]

	return val, found
}

func (c *loadCondition) GetString(fieldSetKey, fieldKey string) (val string, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.(string)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to string", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetStrings(fieldSetKey, fieldKey string) (val []string, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.([]string)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to []string", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetInt(fieldSetKey, fieldKey string) (val int, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.(int)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to int", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetInts(fieldSetKey, fieldKey string) (val []int, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.([]int)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to []int", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetBool(fieldSetKey, fieldKey string) (val, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.(bool)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to bool", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetBools(fieldSetKey, fieldKey string) (val []bool, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.([]bool)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to []bool", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetTime(fieldSetKey, fieldKey string) (val time.Time, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.(time.Time)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to time.Time", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetTimes(fieldSetKey, fieldKey string) (val []time.Time, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.([]time.Time)
	if !ok {
		err = fmt.Errorf("problem casting field (%s.%s) value '%v' to []time.Time", fieldSetKey, fieldKey, fieldValue)

		return
	}

	return
}

func (c *loadCondition) GetDuration(fieldSetKey, fieldKey string) (val time.Duration, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.(time.Duration)
	if !ok {
		err = fmt.Errorf(
			"problem casting field (%s.%s) value '%v' to time.Duration",
			fieldSetKey, fieldKey, fieldValue,
		)

		return
	}

	return
}

func (c *loadCondition) GetDurations(fieldSetKey, fieldKey string) (val []time.Duration, found bool, err error) {
	fieldValue, found := c.GetFieldValue(fieldSetKey, fieldKey)
	if !found {
		return
	}

	val, ok := fieldValue.([]time.Duration)
	if !ok {
		err = fmt.Errorf(
			"problem casting field (%s.%s) value '%v' to []time.Duration",
			fieldSetKey, fieldKey, fieldValue,
		)

		return
	}

	return
}
