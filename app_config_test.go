package bconf_test

import (
	"testing"
	"time"

	"github.com/xavi-group/bconf"
)

func TestAppConfigHelpString(t *testing.T) {
	appConfig := createBaseAppConfig()

	appConfig.AddFieldSet(bconf.FSB("example").Fields(
		bconf.FB("field", bconf.String).Description(
			"this tests a very long field description that would require line wrapping in order to be formatted ",
			"correctly according to the bconf help output logic. This should be splitting lines by word, e.g. ",
			"whenever a space is present. Areallylongwordshouldbesomewhatofaproblem but still handleable.",
		).C(),
	).C())

	appConfig.Load()

	t.Log(appConfig.HelpString())
}

// func TestAppConfig(t *testing.T) {
// 	const appName = "bconf_test_app"

// 	const appDescription = "Test-App is an HTTP server providing access to weather data"

// 	appConfig := bconf.NewAppConfig(
// 		appName,
// 		appDescription,
// 		bconf.WithEnvironmentLoader("bconf_test"),
// 	)

// 	if appConfig.AppName() != appName {
// 		t.Errorf("unexpected value returned from AppName(): '%s'", appConfig.AppName())
// 	}

// 	if appConfig.AppDescription() != appDescription {
// 		t.Errorf("unexpected value returned from AppDescription(): '%s'", appConfig.AppDescription())
// 	}

// 	const appGeneratedID = "generated-default"

// 	appFieldSet := &bconf.FieldSet{
// 		Key: "app",
// 		Fields: bconf.Fields{
// 			{
// 				Key:         "id",
// 				Type:        bconf.String,
// 				Description: "Application identifier for use in application log messages and tracing",
// 				DefaultGenerator: func() (any, error) {
// 					return appGeneratedID, nil
// 				},
// 			},
// 			{
// 				Key:         "read_timeout",
// 				Type:        bconf.Duration,
// 				Description: "Application read timeout for HTTP requests",
// 				Default:     5 * time.Second,
// 			},
// 			{
// 				Key:     "connect_sqlite",
// 				Type:    bconf.Bool,
// 				Default: true,
// 			},
// 		},
// 	}

// 	conditionalFieldSet := &bconf.FieldSet{
// 		Key: "sqlite",
// 		Fields: bconf.Fields{
// 			{
// 				Key:      "server",
// 				Type:     bconfconst.String,
// 				Required: true,
// 			},
// 		},
// 		LoadConditions: bconf.LoadConditions{
// 			&bconf.FieldCondition{
// 				FieldSetKey: "app",
// 				FieldKey:    "connect_sqlite",
// 				Condition: func(fieldValue any) (bool, error) {
// 					val, ok := fieldValue.(bool)
// 					if !ok {
// 						return false, fmt.Errorf("unexpected field-type value")
// 					}

// 					return val, nil
// 				},
// 			},
// 		},
// 	}

// 	appConfig.AddFieldSet(conditionalFieldSet)
// 	appConfig.AddFieldSet(appFieldSet)
// 	appConfig.AddFieldSet(appFieldSet)
// 	appConfig.AddFieldSet(conditionalFieldSet)

// 	if errs := appConfig.Load(); len(errs) < 1 {
// 		t.Fatalf("errors expected for unset required fields")
// 	}

// 	os.Setenv("BCONF_TEST_SQLITE_SERVER", "localhost")

// 	if errs := appConfig.Load(); len(errs) > 0 {
// 		t.Fatalf("unexpected errors registering application configuration: %v", errs)
// 	}
// }

// func TestBadAppConfigFields(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"bconf_test_app",
// 		"Test-App is an HTTP server providing access to weather data",
// 		bconf.WithEnvironmentLoader("bconf_test"),
// 	)

// 	idFieldInvalidDefaultGenerator := &bconf.Field{
// 		Key:         "id",
// 		Type:        bconfconst.Int,
// 		Description: "Application identifier for use in application log messages and tracing",
// 		DefaultGenerator: func() (any, error) {
// 			return "generated-default", nil
// 		},
// 	}
// 	readTimeoutFieldInvalidDefault := &bconf.Field{
// 		Key:         "read_timeout",
// 		Type:        bconfconst.Duration,
// 		Description: "Application read timeout for HTTP requests",
// 		Default:     5,
// 	}
// 	emptyFieldSet := &bconf.FieldSet{}

// 	if errs := appConfig.AddFieldSet(emptyFieldSet); len(errs) < 1 {
// 		t.Fatalf("expected error adding empty field set")
// 	}

// 	invalidAppFieldSet := &bconf.FieldSet{
// 		Key: "app",
// 		Fields: bconf.Fields{
// 			idFieldInvalidDefaultGenerator,
// 			readTimeoutFieldInvalidDefault,
// 		},
// 	}

// 	if errs := appConfig.AddFieldSet(invalidAppFieldSet); len(errs) < 1 {
// 		t.Fatalf("expected errors adding field set with invalid fields")
// 	}

// 	fieldSetWithEmptyField := &bconf.FieldSet{
// 		Key: "default",
// 		Fields: bconf.Fields{
// 			{},
// 		},
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithEmptyField); len(errs) < 2 {
// 		t.Fatalf("expected at least two errors adding a field-set with an empty field")
// 	}

// 	fieldWithDefaultAndRequiredSet := &bconf.Field{
// 		Key:      "log_level",
// 		Type:     bconfconst.String,
// 		Default:  "info",
// 		Required: true,
// 	}

// 	fieldWithDefaultNotInEnumeration := &bconf.Field{
// 		Key:         "log_level",
// 		Type:        bconfconst.String,
// 		Default:     "fatal",
// 		Enumeration: []any{"debug", "info", "warn", "error"},
// 	}

// 	fieldWithGeneratedDefaultNotInEnumeration := &bconf.Field{
// 		Key:  "log_level",
// 		Type: bconfconst.String,
// 		DefaultGenerator: func() (any, error) {
// 			return "fatal", nil
// 		},
// 		Enumeration: []any{"debug", "info", "warn", "error"},
// 	}

// 	fieldSetWithInvalidField := &bconf.FieldSet{
// 		Key:    "default",
// 		Fields: bconf.Fields{fieldWithDefaultAndRequiredSet},
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
// 		t.Fatalf("expected an error adding field with default and required set")
// 	}

// 	fieldSetWithInvalidField.Fields = bconf.Fields{fieldWithDefaultNotInEnumeration}

// 	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
// 		t.Fatalf("expected an error adding field with default value not in enumeration")
// 	}

// 	fieldSetWithInvalidField.Fields = bconf.Fields{fieldWithGeneratedDefaultNotInEnumeration}

// 	if errs := appConfig.AddFieldSet(fieldSetWithInvalidField); len(errs) < 1 {
// 		t.Fatalf("expected an error adding field with generated default value not in enumeration")
// 	}
// }

// func TestAppConfigWithLoadConditions(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	const defaultFieldSetKey = "default"

// 	const defaultFieldSetLoadAppOneKey = "load_app_one"

// 	const defaultFieldSetLoadAppTwoKey = "load_app_two"

// 	loadAppOneField := &bconf.Field{
// 		Key:      defaultFieldSetLoadAppOneKey,
// 		Type:     bconfconst.Bool,
// 		Required: true,
// 	}

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key:    "default",
// 		Fields: bconf.Fields{loadAppOneField},
// 	}

// 	fieldSetWithLoadCondition := &bconf.FieldSet{
// 		Key: "app_one",
// 		Fields: bconf.Fields{
// 			{
// 				Key:      "svc_database_host",
// 				Type:     bconfconst.String,
// 				Required: true,
// 			},
// 		},
// 		LoadConditions: bconf.LoadConditions{
// 			&bconf.FieldCondition{
// 				FieldSetKey: defaultFieldSetKey,
// 				FieldKey:    defaultFieldSetLoadAppOneKey,
// 				Condition: func(fieldValue any) (bool, error) {
// 					return true, nil
// 				},
// 			},
// 		},
// 	}

// 	fieldSetWithUnmetLoadCondition := &bconf.FieldSet{
// 		Key:    "app_two",
// 		Fields: bconf.Fields{},
// 		LoadConditions: bconf.LoadConditions{
// 			&bconf.FieldCondition{
// 				FieldSetKey: defaultFieldSetKey,
// 				FieldKey:    defaultFieldSetLoadAppTwoKey,
// 				Condition: func(fieldValue any) (bool, error) {
// 					return true, nil
// 				},
// 			},
// 		},
// 	}

// 	fieldSetWithInvalidLoadCondition := &bconf.FieldSet{
// 		Key:    "app_three",
// 		Fields: bconf.Fields{},
// 		LoadConditions: bconf.LoadConditions{
// 			&bconf.FieldCondition{
// 				FieldSetKey: defaultFieldSetKey,
// 				Condition: func(fieldValue any) (bool, error) {
// 					return true, nil
// 				},
// 			},
// 		},
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithLoadCondition); len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with unmet field-set load condition")
// 	}

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithUnmetLoadCondition); len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with unmet field load condition")
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithLoadCondition); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field set with valid load condition: %v", errs)
// 	}

// 	_ = os.Setenv("DEFAULT_LOAD_APP_ONE", "true")
// 	_ = os.Setenv("APP_ONE_SVC_DATABASE_HOST", "localhost")

// 	if errs := appConfig.Register(false); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) loading field set with valid load condition and required field: %v", errs)
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithInvalidLoadCondition); len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with invalid load condition")
// 	}
// }

// func TestAppConfigWithFieldLoadConditions(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	const fieldSetOneKey = "one"

// 	const fieldSetTwoKey = "two"

// 	const fieldSetThreeKey = "three"

// 	const fieldSetFourKey = "four"

// 	const fieldAKey = "a"

// 	const fieldBKey = "b"

// 	const fieldCKey = "c"

// 	const fieldDKey = "d"

// 	const fieldEKey = "e"

// 	const fieldFKey = "f"

// 	const fieldGKey = "g"

// 	const fieldBEnvValue = "some_str"

// 	const fieldDEnvValue = "should_not_load"

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", fieldSetOneKey, fieldBKey)), fieldBEnvValue)
// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", fieldSetOneKey, fieldDKey)), fieldDEnvValue)
// 	// test field load condition w/ invalid field-set key
// 	// test help string output with no field-set key
// 	// test help string output with field-set key
// 	// test conditionally required field cannot be unset

// 	fieldSetWithInternalFieldDependencies := bconf.FSB().Key(fieldSetOneKey).Fields(
// 		bconf.FB().Key(fieldAKey).Type(bconf.String).Default("postgres").Create(),
// 		bconf.FB().Key(fieldBKey).Type(bconf.String).LoadConditions(
// 			bconf.FCB().FieldKey(fieldAKey).Condition(func(val any) (bool, error) {
// 				return true, nil
// 			}).Create(),
// 		).Create(),
// 		bconf.FB().Key(fieldCKey).Type(bconf.String).LoadConditions(
// 			bconf.
// 				FCB().
// 				FieldKey(fieldAKey).
// 				FieldSetKey(fieldSetOneKey).
// 				Condition(func(val any) (bool, error) {
// 					return true, nil
// 				}).Create(),
// 		).Create(),
// 		bconf.FB().Key(fieldDKey).Type(bconf.String).Default("should_not_be_overridden").LoadConditions(
// 			bconf.FCB().FieldKey(fieldAKey).Condition(func(val any) (bool, error) {
// 				return false, nil
// 			}).Create(),
// 		).Create(),
// 	).Create()

// 	fieldSetWithOtherFieldSetDependencies := bconf.FSB().Key(fieldSetTwoKey).Fields(
// 		bconf.FB().Key(fieldEKey).Type(bconf.String).Create(),
// 	).Create()

// 	errs := appConfig.AddFieldSets(fieldSetWithInternalFieldDependencies, fieldSetWithOtherFieldSetDependencies)
// 	if len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-sets: %v", errs)
// 	}

// 	if errs := appConfig.Register(false); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) loading field set with valid load condition and required field: %v", errs)
// 	}

// 	foundBValue, _ := appConfig.GetString(fieldSetOneKey, fieldBKey)
// 	if foundBValue != fieldBEnvValue {
// 		t.Errorf("unexpected value found for field B: '%s'", foundBValue)
// 	}

// 	foundDValue, _ := appConfig.GetString(fieldSetOneKey, fieldDKey)
// 	if foundDValue == fieldDEnvValue {
// 		t.Errorf("unexpected value found for field D: '%s'", fieldDEnvValue)
// 	}

// 	fieldSetWithMissingInternalFieldDependencies := bconf.FSB().Key(fieldSetThreeKey).Fields(
// 		bconf.FB().Key(fieldFKey).Type(bconf.String).Create(),
// 		bconf.FB().Key(fieldGKey).Type(bconf.String).LoadConditions(
// 			bconf.FCB().FieldKey(fieldAKey).Condition(func(val any) (bool, error) {
// 				return true, nil
// 			}).Create(),
// 		).Create(),
// 	).Create()

// 	fieldSetWithMissingExternalFieldDependencies := bconf.FSB().Key(fieldSetFourKey).Fields(
// 		bconf.FB().Key(fieldAKey).Type(bconf.String).LoadConditions(
// 			bconf.FCB().FieldSetKey("missing").FieldKey(fieldBKey).Condition(func(val any) (bool, error) {
// 				return true, nil
// 			}).Create(),
// 		).Create(),
// 	).Create()

// 	if errs := appConfig.AddFieldSets(fieldSetWithMissingInternalFieldDependencies); len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with missing internal field dependencies")
// 	} else if !strings.Contains(errs[0].Error(), "field-set field not found") {
// 		t.Errorf("unexpected error adding field set with missing internal field dependencies: '%s'", errs[0])
// 	}

// 	if errs := appConfig.AddFieldSet(fieldSetWithMissingExternalFieldDependencies); len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with missing external field dependencies")
// 	} else if !strings.Contains(errs[0].Error(), "field-set dependency not found") {
// 		t.Errorf("unexpected error adding field set with missing external field dependencies: '%s'", errs[0])
// 	}
// }

// func TestAppConfigAddFieldSets(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	fieldSetOne := &bconf.FieldSet{
// 		Key:    "one",
// 		Fields: bconf.Fields{},
// 	}
// 	fieldSetTwo := &bconf.FieldSet{
// 		Key:    "two",
// 		Fields: bconf.Fields{},
// 	}
// 	fieldSetThree := &bconf.FieldSet{
// 		Key:    "three",
// 		Fields: bconf.Fields{},
// 	}
// 	fieldSetFour := &bconf.FieldSet{
// 		Fields: bconf.Fields{},
// 	}

// 	if errs := appConfig.AddFieldSets(fieldSetOne, fieldSetTwo); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-sets: %v", errs)
// 	}

// 	if errs := appConfig.AddFieldSets(fieldSetThree, fieldSetFour); len(errs) < 1 {
// 		t.Fatalf("expected an error adding field-set with missing key")
// 	} else if !strings.Contains(errs[0].Error(), "field-set key required") {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	if keys := appConfig.GetFieldSetKeys(); len(keys) != 2 {
// 		t.Fatalf("unexpected number of field-sets found on app config: %d", len(keys))
// 	}
// }

// func TestAppConfigAddField(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	fieldSetOne := &bconf.FieldSet{
// 		Key:    "one",
// 		Fields: bconf.Fields{},
// 	}

// 	idFieldKey := "id"
// 	idFieldGeneratedDefaultValue := "generated-default-value"
// 	idField := &bconf.Field{
// 		Key:         idFieldKey,
// 		Type:        bconf.String,
// 		Description: "Application identifier for use in application log messages and tracing",
// 		DefaultGenerator: func() (any, error) {
// 			return idFieldGeneratedDefaultValue, nil
// 		},
// 	}

// 	fieldWithGenerateDefaultError := &bconf.Field{
// 		Key:  "field_generate_default_error",
// 		Type: bconf.String,
// 		DefaultGenerator: func() (any, error) {
// 			return "", errors.New("generated error")
// 		},
// 	}

// 	fieldMissingFieldType := &bconf.Field{
// 		Key: "field_missing_field_type",
// 	}

// 	fieldMissingLoadCondition := bconf.FB().Key("field_missing_load_condition").Type(bconf.String).LoadConditions(
// 		bconf.FCB().FieldKey("missing_key").Condition(
// 			func(val any) (bool, error) {
// 				return true, nil
// 			},
// 		).Create(),
// 	).Create()

// 	if errs := appConfig.AddFieldSets(fieldSetOne); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-sets: %v", errs)
// 	}

// 	if errs := appConfig.AddField("one", idField); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field: %v", errs)
// 	}

// 	if errs := appConfig.AddField("one", idField); len(errs) < 1 {
// 		t.Fatalf("expected error trying to add duplicate field to field-set")
// 	}

// 	if errs := appConfig.AddField("undefined_field_set_key", idField); len(errs) < 1 {
// 		t.Fatalf("expected error trying to add field to undefined field-set")
// 	}

// 	if errs := appConfig.AddField("one", fieldWithGenerateDefaultError); len(errs) < 1 {
// 		t.Fatalf("expected error trying to add field with bad generated default")
// 	}

// 	if errs := appConfig.AddField("one", fieldMissingFieldType); len(errs) < 1 {
// 		t.Fatalf("expected error trying to add field with missing field-type")
// 	}

// 	if errs := appConfig.AddField("one", fieldMissingLoadCondition); len(errs) < 1 {
// 		t.Fatalf("expected error trying to add field with missing load condition")
// 	}
// }

// func TestAppConfigLoadFieldSet(t *testing.T) {
// 	appConfig := createBaseAppConfig()
// 	appConfig.Register(false)

// 	errs := appConfig.LoadFieldSet("field_set_key")

// 	if len(errs) < 1 {
// 		t.Fatalf("unexpected errors length when loading non-existent field-set: %d", len(errs))
// 	}

// 	if !strings.Contains(errs[0].Error(), "field-set with key 'field_set_key' not found") {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	errs = appConfig.AddFieldSet(bconf.FSB().Key("default").Fields(
// 		bconf.FB().Key("field_key").Type(bconf.String).Default("value").Create(),
// 	).Create())

// 	if len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	errs = appConfig.AddFieldSet(bconf.FSB().Key("bad_field_condition_conditional").Fields(
// 		bconf.FB().Key("some_key").Type(bconf.String).Default("value").Create(),
// 	).LoadConditions(
// 		bconf.NewFieldConditionBuilder().FieldSetKey("default").FieldKey("field_key").Condition(
// 			func(fieldValue any) (bool, error) {
// 				return true, fmt.Errorf("condition error")
// 			},
// 		).Create(),
// 	).Create())

// 	if len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	errs = appConfig.LoadFieldSet("bad_field_condition_conditional")

// 	if len(errs) < 1 {
// 		t.Fatalf("expected errors loading conditional field-set")
// 	}

// 	if !strings.Contains(errs[0].Error(), "problem getting load condition outcome") {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}
// }

// func TestAppConfigLoadField(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	standardFieldSetC := bconf.FSB().Key("standard_c").Fields(
// 		bconf.FB().Key("field_a").Type(bconf.String).Default("value_a").Create(),
// 		bconf.FB().Key("field_b").Type(bconf.String).Create(),
// 	).Create()

// 	standardFieldSetD := bconf.FSB().Key("standard_d").Fields(
// 		bconf.FB().Key("field_a").Type(bconf.String).Default("value_a").Create(),
// 		bconf.FB().Key("field_b").Type(bconf.String).Create(),
// 	).Create()

// 	standardFieldSetE := bconf.FSB().Key("standard_e").Fields(
// 		bconf.FB().Key("field_a").Type(bconf.String).Default("value_a").Create(),
// 		bconf.FB().Key("field_b").Type(bconf.String).Create(),
// 	).Create()

// 	appConfig.AddFieldSets(standardFieldSetC, standardFieldSetD, standardFieldSetE)

// 	appConfig.Register(false)

// 	errs := appConfig.LoadField("unk_field_set_key", "unk_field_key")
// 	if len(errs) < 1 {
// 		t.Fatalf("unexpected errors length when loading non-existent field-set: %d", len(errs))
// 	}

// 	if !strings.Contains(errs[0].Error(), "field-set with key 'unk_field_set_key' not found") {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	errs = appConfig.AddFieldSet(bconf.FSB().Key("default").Fields(
// 		bconf.FB().Key("field_key").Type(bconf.String).Default("value").Create(),
// 	).Create())

// 	if len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	errs = appConfig.LoadField("default", "unk_field_key")
// 	if len(errs) < 1 {
// 		t.Fatalf("unexpected errors length when loading non-existent field-set field: %d", len(errs))
// 	}

// 	if !strings.Contains(errs[0].Error(), "field with key 'unk_field_key' not found") {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	addedFieldC := bconf.FB().Key("field_c").Type(bconf.String).LoadConditions(
// 		bconf.FCB().FieldKey("field_b").Condition(func(val any) (bool, error) {
// 			return true, nil
// 		}).Create(),
// 	).Create()

// 	errs = appConfig.AddField("standard_c", addedFieldC)
// 	if len(errs) > 0 {
// 		t.Fatalf("unexpected error adding field set to standard_c: %s", errs)
// 	}

// 	errs = appConfig.LoadField("standard_c", "field_c")
// 	if len(errs) < 1 {
// 		t.Fatalf("expected error adding field set with missing load condition field value")
// 	}

// 	if !strings.Contains(errs[0].Error(), "no value set for field") {
// 		t.Errorf("unexpected error message: %s", errs[0])
// 	}

// 	// A test case for loading a field with a truthy field load condition
// 	os.Setenv("STANDARD_D_FIELD_D", "expected_value")

// 	addedFieldD := bconf.FB().Key("field_d").Type(bconf.String).LoadConditions(
// 		bconf.FCB().FieldKey("field_a").Condition(func(val any) (bool, error) {
// 			return true, nil
// 		}).Create(),
// 	).Create()

// 	if errs := appConfig.AddField("standard_d", addedFieldD); len(errs) > 0 {
// 		t.Fatalf("unexpected error adding field set to standard_d: %s", errs)
// 	}

// 	if errs := appConfig.LoadField("standard_d", "field_d"); len(errs) > 1 {
// 		t.Fatalf("unexpected error loading field with truthy load condition: %s", errs)
// 	}

// 	foundValue, err := appConfig.GetString("standard_d", "field_d")
// 	if err != nil {
// 		t.Fatalf("unexpected error getting standard_d field_d value: %s", err)
// 	}

// 	if foundValue != "expected_value" {
// 		t.Errorf("unexpected value loaded from field, expected 'expected_value', found: '%s'", foundValue)
// 	}

// 	// A test case for loading a field with a falsy field load condition
// 	addedFieldE := bconf.FB().Key("field_e").Type(bconf.String).LoadConditions(
// 		bconf.FCB().FieldKey("field_a").Condition(func(val any) (bool, error) {
// 			return false, nil
// 		}).Create(),
// 	).Create()

// 	if errs := appConfig.AddField("standard_e", addedFieldE); len(errs) > 0 {
// 		t.Fatalf("unexpected error adding field set to standard_e: %s", errs)
// 	}

// 	if errs := appConfig.LoadField("standard_e", "field_e"); len(errs) < 1 {
// 		t.Fatalf("expected error loading field with false load condition")
// 	}
// }

// func TestAppConfigObservability(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	idFieldKey := "id"
// 	idFieldGeneratedDefaultValue := "generated-default-value"
// 	idField := &bconf.Field{
// 		Key:         idFieldKey,
// 		Type:        bconfconst.String,
// 		Description: "Application identifier for use in application log messages and tracing",
// 		DefaultGenerator: func() (any, error) {
// 			return idFieldGeneratedDefaultValue, nil
// 		},
// 	}

// 	sessionSecretFieldKey := "session_secret"
// 	sessionSecretEnvironmentValue := "environment-session-secret-value"
// 	sessionSecretField := &bconf.Field{
// 		Key:       sessionSecretFieldKey,
// 		Type:      bconfconst.String,
// 		Sensitive: true,
// 		Validator: func(fieldValue any) error {
// 			return nil
// 		},
// 	}

// 	timeoutFieldKey := "timeout"
// 	timeoutDefaultValue := 30 * time.Second
// 	timeoutField := bconf.FB().Key(timeoutFieldKey).Type(bconf.Duration).Default(timeoutDefaultValue).Create()

// 	os.Setenv("APP_SESSION_SECRET", sessionSecretEnvironmentValue)

// 	appFieldSetKey := "app"
// 	appFieldSet := &bconf.FieldSet{
// 		Key:    appFieldSetKey,
// 		Fields: bconf.Fields{idField, sessionSecretField, timeoutField},
// 	}

// 	if errs := appConfig.AddFieldSet(appFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected errors adding app field set: %v", errs)
// 	}

// 	foundFieldSetKeys := appConfig.GetFieldSetKeys()
// 	if len(foundFieldSetKeys) != 1 {
// 		t.Fatalf("unexpected length of field set keys returned from app config: %d", len(foundFieldSetKeys))
// 	}

// 	if foundFieldSetKeys[0] != appFieldSetKey {
// 		t.Fatalf("unexpected field-set key in keys returned from app config: '%s'", foundFieldSetKeys[0])
// 	}

// 	foundAppFieldSetKeys, err := appConfig.GetFieldSetFieldKeys(appFieldSetKey)
// 	if err != nil {
// 		t.Fatalf("unexpected issue getting app field-set field keys: %s", err)
// 	}

// 	if len(foundAppFieldSetKeys) < len(appFieldSet.Fields) {
// 		t.Fatalf("length of field-set field keys does not match the length of fields: %d", len(foundAppFieldSetKeys))
// 	}

// 	fieldMap := appConfig.ConfigMap()

// 	if _, found := fieldMap[appFieldSetKey]; !found {
// 		t.Fatalf("expected to find app field-set key in config map")
// 	}

// 	if _, found := fieldMap[appFieldSetKey][idFieldKey]; !found {
// 		t.Fatalf("expected to find app id key in config map")
// 	}

// 	if _, found := fieldMap[appFieldSetKey][sessionSecretFieldKey]; found {
// 		t.Fatalf("unexpected session-secret key found in config map when no value is set")
// 	}

// 	if errs := appConfig.Register(false); errs != nil {
// 		t.Fatalf("unexpected errors registering app config: %v", errs)
// 	}

// 	fieldMap = appConfig.ConfigMap()

// 	fieldMapValue, found := fieldMap[appFieldSetKey][sessionSecretFieldKey]
// 	if !found {
// 		t.Fatalf("expected to find session-secret key in config map")
// 	}

// 	if fieldMapValue == sessionSecretEnvironmentValue {
// 		t.Fatalf(
// 			"unexpected sensitive value (%s) output in config map values: '%s'",
// 			sessionSecretFieldKey,
// 			fieldMapValue,
// 		)
// 	}

// 	fieldMapValue, found = fieldMap[appFieldSetKey][fmt.Sprintf("%s_ms", timeoutFieldKey)]
// 	if !found {
// 		t.Fatalf("expected to find timeout_ms key in config map")
// 	}

// 	if fieldMapValue != int64(30000) {
// 		t.Fatalf("unexpected timeout_ms value: %v", fieldMapValue)
// 	}
// }

// func TestAppConfigSetField(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	const stringFieldKey = "string"

// 	const stringFieldValue = "string_one"

// 	stringField := &bconf.Field{
// 		Key:         stringFieldKey,
// 		Type:        bconfconst.String,
// 		Default:     stringFieldValue,
// 		Enumeration: []any{"string_one", "string_two", "string_three"},
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			stringField,
// 		},
// 	}

// 	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, "some_val"); err == nil {
// 		t.Fatalf("expected error setting field when field-set is not present")
// 	} else if !strings.Contains(err.Error(), fmt.Sprintf("field-set with key '%s' not found", defaultFieldSetKey)) {
// 		t.Fatalf("unexpected error message: %s", err.Error())
// 	}

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, 3928482); err == nil {
// 		t.Fatalf("expected error setting field to mismatched type")
// 	} else if !strings.Contains(err.Error(), "invalid value field-type") {
// 		t.Fatalf("unexpected error message when setting field to mismatched field-type: %s", err)
// 	}

// 	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, "string_zero"); err == nil {
// 		t.Fatalf("expected error setting field to value not in enumeration list")
// 	} else if !strings.Contains(err.Error(), "value not found in enumeration list") {
// 		t.Fatalf("unexpected error message when setting field to value not in enumeraiton list: %s", err)
// 	}

// 	if err := appConfig.SetField(defaultFieldSetKey, "some_key", "some_val"); err == nil {
// 		t.Fatalf("expected error setting field when field is not present")
// 	} else if !strings.Contains(err.Error(), "field with key") {
// 		t.Fatalf("unexpected error message: %s", err.Error())
// 	}
// }

// func TestAppConfigReloadingFields(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	const stringFieldKey = "string"

// 	const stringFieldValue = "string_one"

// 	stringField := &bconf.Field{
// 		Key:     stringFieldKey,
// 		Type:    bconfconst.String,
// 		Default: stringFieldValue,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			stringField,
// 		},
// 	}

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	if errs := appConfig.LoadFieldSet(defaultFieldSetKey); len(errs) != 1 {
// 		t.Fatalf("expected error loading field-set before the app-config is registered")
// 	} else if !strings.Contains(errs[0].Error(), "cannot be called before the app-config has been registered") {
// 		t.Fatalf("unexpected error message when loading field-set before app-config is registered: %s", errs[0])
// 	}

// 	if errs := appConfig.LoadField(defaultFieldSetKey, stringFieldKey); len(errs) != 1 {
// 		t.Fatalf("expected error loading field before the app-config is registered")
// 	} else if !strings.Contains(errs[0].Error(), "cannot be called before the app-config has been registered") {
// 		t.Fatalf("unexpected error message when loading field before app-config is registered: %s", errs[0])
// 	}

// 	if errs := appConfig.Register(false); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app-config: %v", errs)
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringFieldKey)), "string_two")

// 	if errs := appConfig.LoadFieldSet(defaultFieldSetKey); len(errs) > 0 {
// 		t.Fatalf("unexpected errors loading field-set: %v", errs)
// 	}

// 	if val, err := appConfig.GetString(defaultFieldSetKey, stringFieldKey); err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if val != "string_two" {
// 		t.Fatalf("unexpected field value: '%s'", val)
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringFieldKey)), "string_three")

// 	if errs := appConfig.LoadField(defaultFieldSetKey, stringFieldKey); len(errs) > 0 {
// 		t.Fatalf("unexpected errors loading field: %v", errs)
// 	}

// 	if val, err := appConfig.GetString(defaultFieldSetKey, stringFieldKey); err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if val != "string_three" {
// 		t.Fatalf("unexpected field value: '%s'", val)
// 	}
// }

// func TestAppConfigFieldValidators(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	const stringFieldKey = "string"

// 	stringFieldValue := "string_one"
// 	validatorExpectedValue := "string_two"
// 	validatorErrorString := fmt.Sprintf("expected value to be '%s'", validatorExpectedValue)
// 	stringField := &bconf.Field{
// 		Key:     stringFieldKey,
// 		Type:    bconfconst.String,
// 		Default: stringFieldValue,
// 		Validator: func(fieldValue any) error {
// 			val, _ := fieldValue.(string)

// 			if val != validatorExpectedValue {
// 				return errors.New(validatorErrorString)
// 			}

// 			return nil
// 		},
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			stringField,
// 		},
// 	}

// 	expectContains := fmt.Sprintf(
// 		"invalid default value: error from field validator: %s",
// 		validatorErrorString,
// 	)

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
// 		t.Fatalf("expected 1 error adding default field-set with default value not passing validator: %v", errs)
// 	} else if !strings.Contains(errs[0].Error(), expectContains) {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	stringField.Default = nil
// 	stringField.DefaultGenerator = func() (any, error) {
// 		return stringFieldValue, nil
// 	}

// 	expectContains = fmt.Sprintf(
// 		"invalid generated default value: error from field validator: %s",
// 		validatorErrorString,
// 	)

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
// 		t.Fatalf(
// 			"expected 1 error adding default field-set with generated default value not passing validator: %v",
// 			errs,
// 		)
// 	} else if !strings.Contains(errs[0].Error(), expectContains) {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}

// 	stringField.Default = validatorExpectedValue
// 	stringField.DefaultGenerator = nil

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding field-set: %v", errs)
// 	}

// 	if err := appConfig.SetField(defaultFieldSetKey, stringFieldKey, stringFieldValue); err == nil {
// 		t.Fatalf("expected error setting field value violating validator func")
// 	}
// }

// func TestAppConfigFieldDefaultGenerators(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"app",
// 		"description",
// 	)

// 	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

// 	const stringFieldKey = "string"

// 	defaultGeneratorError := "problem generating default"
// 	stringField := &bconf.Field{
// 		Key:  stringFieldKey,
// 		Type: bconfconst.String,
// 		DefaultGenerator: func() (any, error) {
// 			return nil, errors.New(defaultGeneratorError)
// 		},
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			stringField,
// 		},
// 	}

// 	expectContains := fmt.Sprintf(
// 		"default value generation error: problem generating default field value: %s",
// 		defaultGeneratorError,
// 	)

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) != 1 {
// 		t.Fatalf(
// 			"expected 1 error adding default field-set with generated default value function error: %v",
// 			errs,
// 		)
// 	} else if !strings.Contains(errs[0].Error(), expectContains) {
// 		t.Fatalf("unexpected error message: %s", errs[0])
// 	}
// }

// func TestAppConfigFieldEnumeration(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	errs := appConfig.AddFieldSets(
// 		bconf.FSB().Key("default").Fields(
// 			bconf.FB().Key("field").Type(bconf.String).Enumeration(1, 2, 3).Create(),
// 		).Create(),
// 	)

// 	if len(errs) != 3 {
// 		t.Fatalf("expected three errors adding a field-set with a field containing invalid enumeration values")
// 	}

// 	if !strings.Contains(errs[0].Error(), "invalid enumeration value type") {
// 		t.Fatalf("unexpected error message: %s", errs[0].Error())
// 	}
// }

// func TestAppConfigStringFieldTypes(t *testing.T) {
// 	appConfig := createBaseAppConfig()

// 	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

// 	stringsFieldKey := "strings"
// 	stringsFieldValue := []string{"string_one", "string_two"}
// 	stringsEnvValue := "string_three, string_four"
// 	stringsParsedEnvValue := []string{"string_three", "string_four"}
// 	stringsField := &bconf.Field{
// 		Key:     stringsFieldKey,
// 		Type:    bconfconst.Strings,
// 		Default: stringsFieldValue,
// 	}

// 	herringFieldKey := "ints"
// 	herringField := &bconf.Field{
// 		Key:  herringFieldKey,
// 		Type: bconfconst.Ints,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			stringsField,
// 			herringField,
// 		},
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, stringsFieldKey)), stringsEnvValue)

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
// 	}

// 	if _, err := appConfig.GetStrings(defaultFieldSetKey, herringFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundStringVals, err := appConfig.GetStrings(defaultFieldSetKey, stringsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundStringVals {
// 		if stringsFieldValue[idx] != val {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val, stringsFieldValue[idx])
// 		}
// 	}

// 	if errs := appConfig.Register(false); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app config: %v", errs)
// 	}

// 	foundStringVals, err = appConfig.GetStrings(defaultFieldSetKey, stringsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundStringVals {
// 		if stringsParsedEnvValue[idx] != val {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val, stringsParsedEnvValue[idx])
// 		}
// 	}
// }

// func TestAppConfigIntFieldTypes(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"app",
// 		"description",
// 	)

// 	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

// 	intFieldKey := "int"
// 	intFieldValue := 1
// 	intEnvValue := "2"
// 	intParsedEnvValue := 2
// 	intField := &bconf.Field{
// 		Key:     intFieldKey,
// 		Type:    bconfconst.Int,
// 		Default: intFieldValue,
// 	}

// 	intsFieldKey := "ints"
// 	intsFieldValue := []int{1, 2}
// 	intsEnvValue := "3, 4"
// 	intsParsedEnvValue := []int{3, 4}
// 	intsField := &bconf.Field{
// 		Key:     intsFieldKey,
// 		Type:    bconfconst.Ints,
// 		Default: intsFieldValue,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			intField,
// 			intsField,
// 		},
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, intFieldKey)), intEnvValue)
// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, intsFieldKey)), intsEnvValue)

// 	if errs := appConfig.AddFieldSet(defaultFieldSet); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) adding default field-set: %v", errs)
// 	}

// 	if _, err := appConfig.GetInt(defaultFieldSetKey, intsFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundIntVal, err := appConfig.GetInt(defaultFieldSetKey, intFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundIntVal != intFieldValue {
// 		t.Errorf("unexpected value found: '%d', expected '%d", foundIntVal, intFieldValue)
// 	}

// 	if _, err = appConfig.GetInts(defaultFieldSetKey, intFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundIntVals, err := appConfig.GetInts(defaultFieldSetKey, intsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundIntVals {
// 		if intsFieldValue[idx] != val {
// 			t.Errorf("unexpected value found: '%d', expected '%d", val, intsFieldValue[idx])
// 		}
// 	}

// 	if errs := appConfig.Register(false); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app config: %v", errs)
// 	}

// 	foundIntVal, err = appConfig.GetInt(defaultFieldSetKey, intFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundIntVal != intParsedEnvValue {
// 		t.Errorf("unexpected value found: '%d', expected '%d", foundIntVal, intParsedEnvValue)
// 	}

// 	foundIntVals, err = appConfig.GetInts(defaultFieldSetKey, intsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundIntVals {
// 		if intsParsedEnvValue[idx] != val {
// 			t.Errorf("unexpected value found: '%d', expected '%d", val, intsParsedEnvValue[idx])
// 		}
// 	}
// }

// func TestAppConfigBoolFieldTypes(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"app",
// 		"description",
// 	)

// 	_ = appConfig.SetLoaders(&bconf.EnvironmentLoader{})

// 	boolFieldKey := "bool"
// 	boolFieldValue := true
// 	boolEnvValue := "false"
// 	boolParsedEnvValue := false
// 	boolField := &bconf.Field{
// 		Key:     boolFieldKey,
// 		Type:    bconfconst.Bool,
// 		Default: boolFieldValue,
// 	}

// 	boolsFieldKey := "bools"
// 	boolsFieldValue := []bool{true, false}
// 	boolsEnvValue := "false, true"
// 	boolsParsedEnvValue := []bool{false, true}
// 	boolsField := &bconf.Field{
// 		Key:     boolsFieldKey,
// 		Type:    bconfconst.Bools,
// 		Default: boolsFieldValue,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			boolField,
// 			boolsField,
// 		},
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, boolFieldKey)), boolEnvValue)
// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, boolsFieldKey)), boolsEnvValue)

// 	appConfig.AddFieldSet(defaultFieldSet)

// 	if _, err := appConfig.GetBool(defaultFieldSetKey, boolsFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundBoolVal, err := appConfig.GetBool(defaultFieldSetKey, boolFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundBoolVal != boolFieldValue {
// 		t.Errorf("unexpected value found: '%v', expected '%v", foundBoolVal, boolFieldValue)
// 	}

// 	if _, err = appConfig.GetBools(defaultFieldSetKey, boolFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundBoolVals, err := appConfig.GetBools(defaultFieldSetKey, boolsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundBoolVals {
// 		if boolsFieldValue[idx] != val {
// 			t.Errorf("unexpected value found: '%v', expected '%v", val, boolsFieldValue[idx])
// 		}
// 	}

// 	if errs := appConfig.Load(); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app config: %v", errs)
// 	}

// 	foundBoolVal, err = appConfig.GetBool(defaultFieldSetKey, boolFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundBoolVal != boolParsedEnvValue {
// 		t.Errorf("unexpected value found: '%v', expected '%v", foundBoolVal, boolParsedEnvValue)
// 	}

// 	foundBoolVals, err = appConfig.GetBools(defaultFieldSetKey, boolsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundBoolVals {
// 		if boolsParsedEnvValue[idx] != val {
// 			t.Errorf("unexpected value found: '%v', expected '%v", val, boolsParsedEnvValue[idx])
// 		}
// 	}
// }

// func TestAppConfigDurationFieldTypes(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"app",
// 		"description",
// 		bconf.WithEnvironmentLoader(""),
// 	)

// 	durationFieldKey := "duration"
// 	durationFieldValue := 1 * time.Minute
// 	durationEnvValue := "1h"
// 	durationParsedEnvValue := 1 * time.Hour
// 	durationField := &bconf.Field{
// 		Key:     durationFieldKey,
// 		Type:    bconfconst.Duration,
// 		Default: durationFieldValue,
// 	}

// 	durationsFieldKey := "durations"
// 	durationsFieldValue := []time.Duration{1 * time.Minute, 1 * time.Hour}
// 	durationsEnvValue := "1h, 1m"
// 	durationsParsedEnvValue := []time.Duration{1 * time.Hour, 1 * time.Minute}
// 	durationsField := &bconf.Field{
// 		Key:     durationsFieldKey,
// 		Type:    bconfconst.Durations,
// 		Default: durationsFieldValue,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			durationField,
// 			durationsField,
// 		},
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, durationFieldKey)), durationEnvValue)
// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, durationsFieldKey)), durationsEnvValue)

// 	appConfig.AddFieldSet(defaultFieldSet)

// 	if _, err := appConfig.GetDuration(defaultFieldSetKey, durationsFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundDurationVal, err := appConfig.GetDuration(defaultFieldSetKey, durationFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundDurationVal != durationFieldValue {
// 		t.Errorf("unexpected value found: '%s', expected '%s", foundDurationVal.String(), durationFieldValue.String())
// 	}

// 	if _, err = appConfig.GetDurations(defaultFieldSetKey, durationFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundDurationVals, err := appConfig.GetDurations(defaultFieldSetKey, durationsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundDurationVals {
// 		if durationsFieldValue[idx] != val {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), durationsFieldValue[idx].String())
// 		}
// 	}

// 	if errs := appConfig.Load(); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app config: %v", errs)
// 	}

// 	foundDurationVal, err = appConfig.GetDuration(defaultFieldSetKey, durationFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundDurationVal != durationParsedEnvValue {
// 		t.Errorf(
// 			"unexpected value found: '%s', expected '%s",
// 			foundDurationVal.String(),
// 			durationParsedEnvValue.String(),
// 		)
// 	}

// 	foundDurationVals, err = appConfig.GetDurations(defaultFieldSetKey, durationsFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundDurationVals {
// 		if durationsParsedEnvValue[idx] != val {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), durationsParsedEnvValue[idx].String())
// 		}
// 	}
// }

// func TestAppConfigTimeFieldTypes(t *testing.T) {
// 	appConfig := bconf.NewAppConfig(
// 		"app",
// 		"description",
// 		bconf.WithEnvironmentLoader(""),
// 	)

// 	baseTime := time.Now()

// 	timeFieldKey := "time"
// 	timeFieldValue := baseTime
// 	timeEnvValue := baseTime.Add(-1 * time.Hour).Format(time.RFC3339)
// 	timeParsedEnvValue := baseTime.Add(-1 * time.Hour)
// 	timeField := &bconf.Field{
// 		Key:     timeFieldKey,
// 		Type:    bconfconst.Time,
// 		Default: timeFieldValue,
// 	}

// 	timesFieldKey := "times"
// 	timesFieldValue := []time.Time{baseTime, baseTime.Add(-1 * time.Hour)}
// 	timesEnvValue := fmt.Sprintf(
// 		"%s, %s",
// 		baseTime.Add(-1*time.Hour).Format(time.RFC3339),
// 		baseTime.Format(time.RFC3339),
// 	)
// 	timesParsedEnvValue := []time.Time{baseTime.Add(-1 * time.Hour), baseTime}
// 	timesField := &bconf.Field{
// 		Key:     timesFieldKey,
// 		Type:    bconfconst.Times,
// 		Default: timesFieldValue,
// 	}

// 	const defaultFieldSetKey = "default"

// 	defaultFieldSet := &bconf.FieldSet{
// 		Key: defaultFieldSetKey,
// 		Fields: bconf.Fields{
// 			timeField,
// 			timesField,
// 		},
// 	}

// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, timeFieldKey)), timeEnvValue)
// 	os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", defaultFieldSetKey, timesFieldKey)), timesEnvValue)

// 	appConfig.AddFieldSet(defaultFieldSet)

// 	if _, err := appConfig.GetTime(defaultFieldSetKey, timesFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundTimeVal, err := appConfig.GetTime(defaultFieldSetKey, timeFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundTimeVal != timeFieldValue {
// 		t.Errorf("unexpected value found: '%s', expected '%s", foundTimeVal.String(), timeFieldValue.String())
// 	}

// 	if _, err = appConfig.GetTimes(defaultFieldSetKey, timeFieldKey); err == nil {
// 		t.Fatalf("expected error getting mismatched field type")
// 	}

// 	foundTimeVals, err := appConfig.GetTimes(defaultFieldSetKey, timesFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundTimeVals {
// 		if !timesFieldValue[idx].Equal(val) {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), timesFieldValue[idx].String())
// 		}
// 	}

// 	if errs := appConfig.Load(); len(errs) > 0 {
// 		t.Fatalf("unexpected error(s) registering app config: %v", errs)
// 	}

// 	foundTimeVal, err = appConfig.GetTime(defaultFieldSetKey, timeFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	} else if foundTimeVal.Format(time.RFC3339) != timeParsedEnvValue.Format(time.RFC3339) {
// 		t.Errorf("unexpected value found: '%s', expected '%s", foundTimeVal.String(), timeParsedEnvValue.String())
// 	}

// 	foundTimeVals, err = appConfig.GetTimes(defaultFieldSetKey, timesFieldKey)
// 	if err != nil {
// 		t.Fatalf("unexpected error getting field value: %s", err)
// 	}

// 	for idx, val := range foundTimeVals {
// 		if timesParsedEnvValue[idx].Format(time.RFC3339) != val.Format(time.RFC3339) {
// 			t.Errorf("unexpected value found: '%s', expected '%s", val.String(), timesParsedEnvValue[idx].String())
// 		}
// 	}
// }

//nolint:govet // doesn't need to be optimal for tests
type ValidConfigA struct {
	bconf.ConfigStruct `bconf:"api"`
	DBSwitchTime       time.Time     `bconf:"db_switch_time"`
	Host               string        `bconf:"host"`
	ReadTimeout        time.Duration `bconf:"read_timeout"`
	Port               int           `bconf:"port"`
	DebugMode          bool          `bconf:"api.debug_mode"`
	LogPrefix          string        `bconf:"log_prefix"`
}

//nolint:govet // doesn't need to be optimal for tests
type ValidConfigB struct {
	DBSwitchTime time.Time     `bconf:"api.db_switch_time"`
	Host         string        `bconf:"api.host"`
	ReadTimeout  time.Duration `bconf:"api.read_timeout"`
	Port         int           `bconf:"api.port"`
	DebugMode    bool          `bconf:"api.debug_mode"`
	LogPrefix    string        `bconf:"api.log_prefix"`
}

type ValidConfigC struct {
	ConfigA *ValidConfigA
	ConfigB *ValidConfigB
}

func TestValidAppConfigFillStruct(t *testing.T) {
	const (
		host        = "localhost"
		port        = 8080
		readTimeout = 5 * time.Second
		debugMode   = true
	)

	var (
		dbSwitchTime = time.Now().Add(-100 * time.Hour)
	)

	appConfig := createBaseAppConfig()

	appConfig.AddFieldSet(
		bconf.FSB("api").Fields(
			bconf.FB("host", bconf.String).Default(host).C(),
			bconf.FB("port", bconf.Int).Default(port).C(),
			bconf.FB("read_timeout", bconf.Duration).Default(readTimeout).C(),
			bconf.FB("db_switch_time", bconf.Time).Default(dbSwitchTime).C(),
			bconf.FB("debug_mode", bconf.Bool).Default(debugMode).C(),
			bconf.FB("log_prefix", bconf.String).C(),
		).C(),
	)

	if errs := appConfig.Load(); len(errs) > 0 {
		t.Fatalf("unexpected error(s) loading app config: %v\n", errs)
	}

	configStructA := &ValidConfigA{}
	if err := appConfig.FillStruct(configStructA); err != nil {
		t.Fatalf("unexpected error when filling struct 'configStructA': %s\n", err)
	}

	testHelperCheckValidConfigStructA(t, configStructA, dbSwitchTime)

	configStructB := &ValidConfigB{}
	if err := appConfig.FillStruct(configStructB); err != nil {
		t.Fatalf("unexpected error when filling struct 'configStructB': %s\n", err)
	}

	testHelperCheckValidConfigStructB(t, configStructB, dbSwitchTime)

	configStructC := &ValidConfigC{}
	if err := appConfig.FillStruct(configStructC); err != nil {
		t.Fatalf("unexpected error when filling struct 'configStructC': %s\n", err)
	}

	testHelperCheckValidConfigStructA(t, configStructC.ConfigA, dbSwitchTime)
	testHelperCheckValidConfigStructB(t, configStructC.ConfigB, dbSwitchTime)
}

func TestValidAppConfigAttachConfigStructs(t *testing.T) {
	const (
		host        = "localhost"
		port        = 8080
		readTimeout = 5 * time.Second
		debugMode   = true
	)

	var (
		dbSwitchTime = time.Now().Add(-100 * time.Hour)
	)

	appConfig := createBaseAppConfig()

	appConfig.AddFieldSet(
		bconf.FSB("api").Fields(
			bconf.FB("host", bconf.String).Default(host).C(),
			bconf.FB("port", bconf.Int).Default(port).C(),
			bconf.FB("read_timeout", bconf.Duration).Default(readTimeout).C(),
			bconf.FB("db_switch_time", bconf.Time).Default(dbSwitchTime).C(),
			bconf.FB("debug_mode", bconf.Bool).Default(debugMode).C(),
			bconf.FB("log_prefix", bconf.String).C(),
		).C(),
	)

	configStructA := &ValidConfigA{}
	configStructB := &ValidConfigB{}
	configStructC := &ValidConfigC{}

	appConfig.AttachConfigStructs(
		configStructA,
		configStructB,
		configStructC,
	)

	if errs := appConfig.Load(); len(errs) > 0 {
		t.Fatalf("unexpected error(s) loading app config: %v\n", errs)
	}

	testHelperCheckValidConfigStructA(t, configStructA, dbSwitchTime)
	testHelperCheckValidConfigStructB(t, configStructB, dbSwitchTime)
	testHelperCheckValidConfigStructA(t, configStructC.ConfigA, dbSwitchTime)
	testHelperCheckValidConfigStructB(t, configStructC.ConfigB, dbSwitchTime)
}

func testHelperCheckValidConfigStructA(t *testing.T, configStructA *ValidConfigA, dbSwitchTime time.Time) {
	t.Helper()

	const (
		host        = "localhost"
		port        = 8080
		readTimeout = 5 * time.Second
		debugMode   = true
	)

	if configStructA.Host != host {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", host, configStructA.Host)
	}

	if configStructA.Port != port {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", port, configStructA.Port)
	}

	if configStructA.ReadTimeout != readTimeout {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", readTimeout, configStructA.ReadTimeout)
	}

	if !configStructA.DBSwitchTime.Equal(dbSwitchTime) {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", dbSwitchTime, configStructA.DBSwitchTime)
	}

	if configStructA.DebugMode != debugMode {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", debugMode, configStructA.DebugMode)
	}

	if configStructA.LogPrefix != "" {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", "", configStructA.LogPrefix)
	}
}

func testHelperCheckValidConfigStructB(t *testing.T, configStructB *ValidConfigB, dbSwitchTime time.Time) {
	t.Helper()

	const (
		host        = "localhost"
		port        = 8080
		readTimeout = 5 * time.Second
		debugMode   = true
	)

	if configStructB.Host != host {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", host, configStructB.Host)
	}

	if configStructB.Port != port {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", port, configStructB.Port)
	}

	if configStructB.ReadTimeout != readTimeout {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", readTimeout, configStructB.ReadTimeout)
	}

	if !configStructB.DBSwitchTime.Equal(dbSwitchTime) {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", dbSwitchTime, configStructB.DBSwitchTime)
	}

	if configStructB.DebugMode != debugMode {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", debugMode, configStructB.DebugMode)
	}

	if configStructB.LogPrefix != "" {
		t.Errorf("unexpected config value (expected '%v'), found: '%v'\n", "", configStructB.LogPrefix)
	}
}

func createBaseAppConfig() *bconf.AppConfig {
	appConfig := bconf.NewAppConfig(
		"testapp",
		"testapp description",
		bconf.WithEnvironmentLoader(""),
	)

	return appConfig
}
