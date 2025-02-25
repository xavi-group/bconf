package bconf_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rheisen/bconf"
)

func TestFieldBuilderCreate(t *testing.T) {
	var field *bconf.Field

	fieldKey := "field_key"
	fieldType := bconf.String

	field = bconf.NewFieldBuilder(fieldKey, fieldType).Create()
	if field == nil {
		t.Fatal("unexpected nil field from builder create")
	}

	field = bconf.FB(fieldKey, fieldType).Create()
	if field == nil {
		t.Fatal("unexpected nil field from builder create")
	}

	field = bconf.FB(fieldKey, fieldType).C()
	if field == nil {
		t.Fatal("unexpected nil field from builder create")
	}

	if field.Type != fieldType {
		t.Errorf("unexpected field type (expected '%s'), found: '%s'\n", fieldType, field.Type)
	}

	if field.Key != fieldKey {
		t.Errorf("unexpected field key (expected '%s'), found: '%s'\n", fieldKey, field.Key)
	}
}

func TestFieldBuilderDefault(t *testing.T) {
	const fieldDefault = 30 * time.Second

	field := bconf.FB("field_key", bconf.Duration).Default(fieldDefault).Create()
	if field.Default != fieldDefault {
		t.Fatalf("unexpected field default '%v', expected '%v'\n", field.Default, fieldDefault)
	}
}

func TestFieldBuilderValidator(t *testing.T) {
	validator := func(fieldValue any) error {
		return fmt.Errorf("validator error")
	}
	field := bconf.FB("field_key", bconf.String).Validator(validator).Create()

	if err := field.Validator(nil); err.Error() != "validator error" {
		t.Fatalf("unexpected validator error value: %s\n", err)
	}
}

func TestFieldBuilderDefaultGenerator(t *testing.T) {
	defaultGenerator := func() (any, error) {
		return "default", nil
	}
	field := bconf.FB("field_key", bconf.String).DefaultGenerator(defaultGenerator).Create()

	generatedDefault, _ := field.DefaultGenerator()
	if generatedDefault != "default" {
		t.Fatalf("unexpected generated default value '%s', expected 'default'\n", generatedDefault)
	}
}

func TestFieldBuilderLoadConditions(t *testing.T) {
	field := bconf.FB("field_key", bconf.String).LoadConditions(
		bconf.LCB(func(c bconf.LoadConditionValues) (bool, error) {
			return true, nil
		}).AddFieldDependencies(
			bconf.FD("field_set_key", "field_key"),
		).C(),
	).C()

	if len(field.LoadConditions) != 1 {
		t.Fatalf("unexpected length of field load conditions (expected 1), found: %d\n", len(field.LoadConditions))
	}
}

func TestFieldBuilderDescription(t *testing.T) {
	const fieldDescription = "field description test"

	field := bconf.FB("field_key", bconf.String).Description(fieldDescription).Create()
	if field.Description != fieldDescription {
		t.Fatalf("unexpected field description (expected '%s'), found: '%s'\n", fieldDescription, field.Description)
	}
}

func TestFieldBuilderEnumeration(t *testing.T) {
	fieldEnumeration := []any{"one", "two", "three"}
	field := bconf.FB("field_key", bconf.String).Enumeration(fieldEnumeration...).Create()

	if len(field.Enumeration) != len(fieldEnumeration) {
		t.Fatalf(
			"unexpected field enumeration length '%d', expected '%d'\n",
			len(field.Enumeration), len(fieldEnumeration),
		)
	}
}

func TestFieldBuilderRequired(t *testing.T) {
	field := bconf.FB("field_key", bconf.String).Required().Create()
	if field.Required == false {
		t.Fatal("expected field to be required")
	}
}

func TestFieldBuilderSensitive(t *testing.T) {
	field := bconf.FB("field_key", bconf.String).Sensitive().Create()

	if field.Sensitive == false {
		t.Fatal("expected field to be sensitive")
	}
}
