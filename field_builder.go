package bconf

import "strings"

func FB(fieldKey, fieldType string) FieldBuilder {
	return NewFieldBuilder(fieldKey, fieldType)
}

func NewFieldBuilder(fieldKey, fieldType string) FieldBuilder {
	return &fieldBuilder{field: &Field{Key: fieldKey, Type: fieldType}}
}

// --------------------------------------------------------------------------------------------------------------------

type FieldBuilder interface {
	Default(value any) FieldBuilder
	Validator(validationFunc func(fieldValue any) error) FieldBuilder
	DefaultGenerator(defaultGeneratorFunc func() (any, error)) FieldBuilder
	LoadConditions(conditions ...LoadCondition) FieldBuilder
	Description(description string, concat ...string) FieldBuilder
	Enumeration(acceptedValues ...any) FieldBuilder
	Required() FieldBuilder
	Sensitive() FieldBuilder
	Create() *Field
	C() *Field
}

// --------------------------------------------------------------------------------------------------------------------

type fieldBuilder struct {
	field *Field
}

func (b *fieldBuilder) Default(value any) FieldBuilder {
	b.field.Default = value

	return b
}

func (b *fieldBuilder) Validator(value func(fieldValue any) error) FieldBuilder {
	b.field.Validator = value

	return b
}

func (b *fieldBuilder) DefaultGenerator(value func() (any, error)) FieldBuilder {
	b.field.DefaultGenerator = value

	return b
}

func (b *fieldBuilder) LoadConditions(value ...LoadCondition) FieldBuilder {
	b.field.LoadConditions = value

	return b
}

func (b *fieldBuilder) Description(value string, concat ...string) FieldBuilder {
	if len(concat) > 0 {
		builder := strings.Builder{}

		builder.WriteString(value)

		for _, concatStr := range concat {
			builder.WriteString(concatStr)
		}

		b.field.Description = builder.String()
	} else {
		b.field.Description = value
	}

	return b
}

func (b *fieldBuilder) Enumeration(value ...any) FieldBuilder {
	b.field.Enumeration = value

	return b
}

func (b *fieldBuilder) Required() FieldBuilder {
	b.field.Required = true

	return b
}

func (b *fieldBuilder) Sensitive() FieldBuilder {
	b.field.Sensitive = true

	return b
}

func (b *fieldBuilder) Create() *Field {
	return b.field.Clone()
}

func (b *fieldBuilder) C() *Field {
	return b.Create()
}
