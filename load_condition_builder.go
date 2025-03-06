package bconf

func NewLoadConditionBuilder(loadFunc func(c FieldValueFinder) (bool, error)) LoadConditionBuilder {
	return &loadConditionBuilder{condition: newLoadCondition(loadFunc)}
}

func LCB(loadFunc func(c FieldValueFinder) (bool, error)) LoadConditionBuilder {
	return NewLoadConditionBuilder(loadFunc)
}

// --------------------------------------------------------------------------------------------------------------------

type LoadConditionBuilder interface {
	AddFieldDependencies(dependencies ...FieldLocation) LoadConditionBuilder
	AddFieldSetDependencies(fieldSetKey string, fieldKeys ...string) LoadConditionBuilder
	Create() LoadCondition
	C() LoadCondition
}

// --------------------------------------------------------------------------------------------------------------------

type loadConditionBuilder struct {
	condition *loadCondition
}

func (b *loadConditionBuilder) AddFieldDependencies(dependencies ...FieldLocation) LoadConditionBuilder {
	b.condition.fieldDependencies = append(b.condition.fieldDependencies, dependencies...)

	return b
}

func (b *loadConditionBuilder) AddFieldSetDependencies(fieldSetKey string, fieldKeys ...string) LoadConditionBuilder {
	for _, fieldKey := range fieldKeys {
		b.condition.fieldDependencies = append(b.condition.fieldDependencies, FieldLocation{
			FieldSetKey: fieldSetKey,
			FieldKey:    fieldKey,
		})
	}

	return b
}

func (b *loadConditionBuilder) Create() LoadCondition {
	return b.condition.Clone()
}

func (b *loadConditionBuilder) C() LoadCondition {
	return b.condition.Clone()
}
