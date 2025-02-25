package bconf

func NewFieldSetBuilder(fieldSetKey string) FieldSetBuilder {
	return &fieldSetBuilder{fieldSet: &FieldSet{Key: fieldSetKey}}
}

func FSB(fieldSetKey string) FieldSetBuilder {
	return NewFieldSetBuilder(fieldSetKey)
}

// --------------------------------------------------------------------------------------------------------------------

type FieldSetBuilder interface {
	Fields(fields ...*Field) FieldSetBuilder
	LoadConditions(conditions ...LoadCondition) FieldSetBuilder
	Create() *FieldSet
	C() *FieldSet
}

// --------------------------------------------------------------------------------------------------------------------

type fieldSetBuilder struct {
	fieldSet *FieldSet
}

func (b *fieldSetBuilder) Fields(fields ...*Field) FieldSetBuilder {
	b.fieldSet.Fields = fields

	return b
}

func (b *fieldSetBuilder) LoadConditions(conditions ...LoadCondition) FieldSetBuilder {
	b.fieldSet.LoadConditions = conditions

	return b
}

func (b *fieldSetBuilder) Create() *FieldSet {
	return b.fieldSet.Clone()
}

func (b *fieldSetBuilder) C() *FieldSet {
	return b.Create()
}
