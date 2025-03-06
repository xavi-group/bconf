package bconf_test

import (
	"testing"

	"github.com/xavi-group/bconf"
)

func TestLoadConditionBuilderCreate(t *testing.T) {
	condition := bconf.NewLoadConditionBuilder(
		func(_ bconf.LoadConditionValues) (bool, error) {
			return true, nil
		},
	).Create()

	if condition == nil {
		t.Fatalf("unexpected nil condition")
	}

	condition = bconf.LCB(
		func(_ bconf.LoadConditionValues) (bool, error) {
			return true, nil
		},
	).Create()

	if condition == nil {
		t.Fatalf("unexpected nil condition")
	}

	ok, err := condition.Load(condition)
	if !ok {
		t.Errorf("unexpected load value: %v\n", ok)
	}

	if err != nil {
		t.Errorf("unexpected load error: %s\n", err)
	}
}

func TestLoadConditionBuilderFieldDependencies(t *testing.T) {
	const fieldSetKey = "test_field_set_key"
	const fieldKey = "test_field_key"
	const fieldVal = "test_field_value"

	condition := bconf.LCB(func(c bconf.LoadConditionValues) (bool, error) {
		val, found := c.GetFieldDependency(fieldSetKey, fieldKey)
		if !found {
			t.Fatal("expected to find field dependency")
		}

		valString, ok := val.(string)
		if !ok {
			t.Fatal("expected field value to be string")
		}

		return valString == fieldVal, nil
	}).AddFieldSetDependencies(fieldSetKey, fieldKey).Create()

	condition.SetFieldDependencies(bconf.FieldValue{
		FieldSetKey: fieldSetKey,
		FieldKey:    fieldKey,
		FieldValue:  fieldVal,
	})

	if load, err := condition.Load(condition); err != nil {
		t.Errorf("unexpected error checking load condition: %s\n", err)
	} else if !load {
		t.Error("expected condition.Load(...) to be true")
	}
}
