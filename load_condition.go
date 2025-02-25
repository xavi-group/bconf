package bconf

import (
	"maps"
	"slices"
)

type LoadConditions []LoadCondition

type LoadCondition interface {
	LoadConditionValues
	Clone() LoadCondition
	FieldDependencies() FieldDependencies
	SetFieldDependencies(fieldValues ...FieldValue)
	Load(c LoadConditionValues) (bool, error)
}

type LoadConditionValues interface {
	GetFieldDependencies() map[FieldDependency]any
	GetFieldDependency(fieldSetKey, fieldKey string) (value any, found bool)
}

func FD(fieldSetKey, fieldKey string) FieldDependency {
	return FieldDependency{
		FieldSetKey: fieldSetKey,
		FieldKey:    fieldKey,
	}
}

type FieldDependency struct {
	FieldSetKey string
	FieldKey    string
}

type FieldDependencies []FieldDependency

type FieldValue struct {
	FieldSetKey string
	FieldKey    string
	FieldValue  any
}

type FieldValues []FieldValue

// --------------------------------------------------------------------------------------------------------------------

func newLoadCondition(loadFunc func(c LoadConditionValues) (bool, error)) *loadCondition {
	return &loadCondition{
		loadFunc:              loadFunc,
		fieldDependencies:     FieldDependencies{},
		fieldDependencyValues: map[FieldDependency]any{},
	}
}

// --------------------------------------------------------------------------------------------------------------------

type loadCondition struct {
	loadFunc              func(c LoadConditionValues) (bool, error)
	fieldDependencies     FieldDependencies
	fieldDependencyValues map[FieldDependency]any
}

func (c *loadCondition) Clone() LoadCondition {
	clone := *c

	clone.fieldDependencies = slices.Clone(c.fieldDependencies)
	clone.fieldDependencyValues = maps.Clone(c.fieldDependencyValues)

	return &clone
}

func (c *loadCondition) FieldDependencies() FieldDependencies {
	return slices.Clone(c.fieldDependencies)
}

func (c *loadCondition) GetFieldDependencies() map[FieldDependency]any {
	return maps.Clone(c.fieldDependencyValues)
}

func (c *loadCondition) GetFieldDependency(fieldSetKey, fieldKey string) (any, bool) {
	val, found := c.fieldDependencyValues[FieldDependency{FieldSetKey: fieldSetKey, FieldKey: fieldKey}]

	return val, found
}

func (c *loadCondition) SetFieldDependencies(fieldValues ...FieldValue) {
	for _, fieldValue := range fieldValues {
		key := FieldDependency{FieldSetKey: fieldValue.FieldSetKey, FieldKey: fieldValue.FieldKey}
		c.fieldDependencyValues[key] = fieldValue.FieldValue
	}
}

func (c *loadCondition) Load(loadConditionValues LoadConditionValues) (bool, error) {
	return c.loadFunc(loadConditionValues)
}
