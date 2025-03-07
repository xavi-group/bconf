package bconf_test

import (
	"testing"

	"github.com/xavi-group/bconf"
)

func TestFieldSetBuilderCreate(t *testing.T) {
	fieldSetKey := "field_set_key"

	fieldSet := bconf.NewFieldSetBuilder(fieldSetKey).Create()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}

	if fieldSet.Key != fieldSetKey {
		t.Errorf("unexpected field set key (expected '%s'), found: '%s'\n", fieldSetKey, fieldSet.Key)
	}

	fieldSet = bconf.FSB(fieldSetKey).Create()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}

	if fieldSet.Key != fieldSetKey {
		t.Errorf("unexpected field set key (expected '%s'), found: '%s'\n", fieldSetKey, fieldSet.Key)
	}

	fieldSet = bconf.FSB(fieldSetKey).C()
	if fieldSet == nil {
		t.Fatalf("unexpected nil field-set")
	}

	if fieldSet.Key != fieldSetKey {
		t.Errorf("unexpected field set key (expected '%s'), found: '%s'\n", fieldSetKey, fieldSet.Key)
	}
}

func TestFieldSetBuilderFields(t *testing.T) {
	fieldKey := "field_key"
	field := bconf.FB(fieldKey, bconf.String).Create()

	fieldSet := bconf.FSB("field_set_key").Fields(field).Create()
	if len(fieldSet.Fields) != 1 {
		t.Fatalf("unexpected fields length '%d', expected 1", len(fieldSet.Fields))
	}

	if fieldSet.Fields[0].Key != fieldKey {
		t.Fatalf("unexpected field key '%s', expected '%s'", fieldSet.Fields[0].Key, fieldKey)
	}
}

func TestFieldSetBuilderLoadConditions(t *testing.T) {
	fieldSet := bconf.FSB("field_set_key").LoadConditions(
		bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
			return true, nil
		}).C(),
	).C()

	if len(fieldSet.LoadConditions) != 1 {
		t.Fatalf("unexpected load-conditions length '%d', expected 1", len(fieldSet.LoadConditions))
	}
}
