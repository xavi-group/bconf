package bconf

func NewFieldSetBuilder(fieldSetKey string) *FieldSetBuilder {
	return &FieldSetBuilder{fieldSet: &FieldSet{Key: fieldSetKey}}
}

func FSB(fieldSetKey string) *FieldSetBuilder {
	return NewFieldSetBuilder(fieldSetKey)
}

type FieldSetBuilder struct {
	fieldSet *FieldSet
}

func (b *FieldSetBuilder) Fields(value ...*Field) *FieldSetBuilder {
	b.init()
	b.fieldSet.Fields = value

	return b
}

func (b *FieldSetBuilder) LoadConditions(value ...LoadCondition) *FieldSetBuilder {
	b.init()
	b.fieldSet.LoadConditions = value

	return b
}

func (b *FieldSetBuilder) Create() *FieldSet {
	b.init()
	return b.fieldSet.Clone()
}

func (b *FieldSetBuilder) C() *FieldSet {
	return b.Create()
}

func (b *FieldSetBuilder) init() {
	if b.fieldSet == nil {
		b.fieldSet = &FieldSet{}
	}
}
