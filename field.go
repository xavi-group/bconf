package bconf

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/xavi-group/bconf/bconfconst"
)

const emptyFieldError = "empty field value"

// Fields is a slice of Field elements providing context for configuration values
type Fields []*Field

// Field is a data structure that provides context for a configuration value
type Field struct {
	// fieldValue contains a mapping of loader names to field value
	fieldValue map[string]any
	// Validator defines a function that runs during validation to check a value against validity constraints
	Validator func(value any) error
	// DefaultGenerator defines a function that creates a base value for a field
	DefaultGenerator func() (any, error)
	// Default defines a base value for a field
	Default any
	// generatedDefault tracks the value generated from the default generator function
	generatedDefault any
	// overrideValue tracks a user set field value
	overrideValue any
	// Key is a required field that defines the field lookup value
	Key string
	// Type is a required field that defines the type of value the field contains
	Type string
	// Description defines a summary of the field contents
	Description string
	// Enumeration defines a list of acceptable inputs for the field value
	Enumeration []any
	// LoadConditions defines the conditions required for a field to load values
	LoadConditions LoadConditions
	// fieldFound is a reverse priority list of where field values were found, e.g. last value has highest priority
	fieldFound []string
	// Required defines whether a field value must be set in order for the field to be valid
	Required bool
	// Sensitive identifies the field value as sensitive
	Sensitive bool
}

func (f *Field) Clone() *Field {
	clone := *f

	clone.fieldFound = slices.Clone(f.fieldFound)
	clone.Enumeration = slices.Clone(f.Enumeration)
	clone.fieldValue = maps.Clone(f.fieldValue)

	if len(f.LoadConditions) > 0 {
		clone.LoadConditions = make(LoadConditions, len(f.LoadConditions))

		for idx, value := range f.LoadConditions {
			clone.LoadConditions[idx] = value.Clone()
		}
	}

	return &clone
}

func (f *Field) generateDefault() error {
	if f.DefaultGenerator == nil {
		return nil
	}

	generatedDefault, err := f.DefaultGenerator()
	if err != nil {
		return fmt.Errorf("problem generating default field value: %w", err)
	}

	f.generatedDefault = generatedDefault

	return nil
}

func (f *Field) validate() []error {
	errs := []error{}

	if f.Key == "" {
		errs = append(errs, fmt.Errorf("invalid key value: cannot be blank"))
	}

	if f.Type == "" {
		errs = append(errs, fmt.Errorf("invalid field-type value: cannot be blank"))
	}

	if len(errs) > 0 {
		return errs
	}

	fieldTypeFound := false

	for _, fieldType := range FieldTypes() {
		if fieldType != f.Type {
			continue
		}

		fieldTypeFound = true

		if fieldErrors := f.validateNoConflictingParams(); len(fieldErrors) > 0 {
			errs = append(errs, fieldErrors...)
		}

		if err := f.validateDefaultFieldType(fieldType); err != nil {
			errs = append(errs, err)
		}

		if err := f.validateGeneratedDefaultFieldType(fieldType); err != nil {
			errs = append(errs, err)
		}

		if validationErrs := f.validateEnumerationValuesFieldType(fieldType); len(validationErrs) > 0 {
			errs = append(errs, validationErrs...)
		}

		// Return here before validating default values existing in enumeration list
		if len(errs) > 0 {
			return errs
		}

		if err := f.validateDefaultValuesInEnumeration(); err != nil {
			errs = append(errs, err)
		}

		if err := f.validateDefaultValuesPassValidatorFunc(); err != nil {
			errs = append(errs, err)
		}
	}

	if !fieldTypeFound {
		errs = append(errs, fmt.Errorf("invalid field type specified: '%s'", f.Type))
	}

	return errs
}

func (f *Field) validateNoConflictingParams() []error {
	errs := []error{}

	if f.Default != nil && f.DefaultGenerator != nil {
		errs = append(errs, fmt.Errorf(bconfconst.ErrorFieldDefaultSetting))
	}

	if f.Required && f.Default != nil || f.Required && f.DefaultGenerator != nil {
		errs = append(errs, fmt.Errorf(bconfconst.ErrorFieldRequiredWithDefault))
	}

	return errs
}

func (f *Field) validateDefaultFieldType(fieldType string) error {
	if f.Default == nil {
		return nil
	}

	if reflect.TypeOf(f.Default).String() == fieldType {
		return nil
	}

	return fmt.Errorf(
		"invalid default type: expected '%s', found '%s'",
		fieldType,
		reflect.TypeOf(f.Default).String(),
	)
}

func (f *Field) validateGeneratedDefaultFieldType(fieldType string) error {
	if f.generatedDefault == nil {
		return nil
	}

	if reflect.TypeOf(f.generatedDefault).String() == fieldType {
		return nil
	}

	return fmt.Errorf(
		"invalid generated default type: expected '%s', found '%s'",
		fieldType,
		reflect.TypeOf(f.generatedDefault).String(),
	)
}

func (f *Field) validateEnumerationValuesFieldType(fieldType string) []error {
	if len(f.Enumeration) < 1 {
		return nil
	}

	errs := []error{}

	for _, val := range f.Enumeration {
		if reflect.TypeOf(val).String() == fieldType {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf(
				"invalid enumeration value type: expected '%s', found '%s'",
				fieldType,
				reflect.TypeOf(val).String(),
			),
		)
	}

	return errs
}

func (f *Field) validateDefaultValuesInEnumeration() error {
	if f.Default != nil && !f.valueInEnumeration(f.Default) {
		return fmt.Errorf(
			"invalid default value: default value '%s' expected in enumeration list",
			f.Default,
		)
	}

	if f.generatedDefault != nil && !f.valueInEnumeration(f.generatedDefault) {
		return fmt.Errorf(
			"invalid generated default value: default value '%s' expected in enumeration list",
			f.Default,
		)
	}

	return nil
}

func (f *Field) validateDefaultValuesPassValidatorFunc() error {
	if f.Default != nil && f.Validator != nil {
		if err := f.Validator(f.Default); err != nil {
			return fmt.Errorf(
				"invalid default value: error from field validator: %w",
				err,
			)
		}
	}

	if f.generatedDefault != nil && f.Validator != nil {
		if err := f.Validator(f.generatedDefault); err != nil {
			return fmt.Errorf(
				"invalid generated default value: error from field validator: %w",
				err,
			)
		}
	}

	return nil
}

func (f *Field) getValue() (any, error) {
	if f.overrideValue != nil {
		return f.overrideValue, nil
	}

	if len(f.fieldFound) > 0 {
		value := f.fieldValue[f.fieldFound[len(f.fieldFound)-1]]

		return value, nil
	}

	if f.Default != nil {
		return f.Default, nil
	}

	if f.generatedDefault != nil {
		return f.generatedDefault, nil
	}

	return nil, errors.New(emptyFieldError)
}

// func (f *Field) getValueFrom(loader string) (any, error) {
// 	if f.fieldValue == nil {
// 		return nil, fmt.Errorf("")
// 	}

// 	value, found := f.fieldValue[loader]
// 	if !found {
// 		return nil, fmt.Errorf("")
// 	}

// 	return value, nil
// }

func (f *Field) set(loaderName, value string) error {
	parsedValue, err := f.parseString(value)
	if err != nil {
		return fmt.Errorf("problem parsing value to field-type: %w", err)
	}

	if !f.valueInEnumeration(parsedValue) {
		return fmt.Errorf("value not found in enumeration list")
	}

	if f.Validator != nil {
		if err := f.Validator(parsedValue); err != nil {
			return fmt.Errorf("value validation error: %w", err)
		}
	}

	if f.fieldValue == nil {
		f.fieldValue = map[string]any{loaderName: parsedValue}
	} else {
		f.fieldValue[loaderName] = value
	}

	if f.fieldFound == nil {
		f.fieldFound = []string{loaderName}
	} else {
		f.fieldFound = append(f.fieldFound, loaderName)
	}

	return nil
}

func (f *Field) setOverride(value any) error {
	if reflect.TypeOf(value).String() != f.Type {
		return fmt.Errorf(
			"invalid value field-type: expected '%s', found '%s'",
			f.Type,
			reflect.TypeOf(value).String(),
		)
	}

	if !f.valueInEnumeration(value) {
		return fmt.Errorf("value not found in enumeration list")
	}

	if f.Validator != nil {
		if err := f.Validator(value); err != nil {
			return fmt.Errorf("value validation error: %w", err)
		}
	}

	f.overrideValue = value

	return nil
}

func (f *Field) parseString(value string) (any, error) {
	switch f.Type {
	case String:
		return value, nil
	case Strings:
		return f.parseToStrings(value), nil
	case Bool:
		return strconv.ParseBool(value)
	case Bools:
		return f.parseToBools(value)
	case Int:
		return strconv.Atoi(value)
	case Ints:
		return f.parseToInts(value)
	case Time:
		return time.Parse(time.RFC3339, value)
	case Times:
		return f.parseToTimes(value)
	case Duration:
		return time.ParseDuration(value)
	case Durations:
		return f.parseToDurations(value)
	default:
		return "", fmt.Errorf("unsupported field type: %s", f.Type)
	}
}

func (f *Field) parseToStrings(value string) []string {
	if value == "" {
		return []string{}
	}

	list := strings.Split(value, ",")
	values := make([]string, len(list))

	for idx, elem := range list {
		parsedValue := strings.Trim(elem, " ")
		values[idx] = parsedValue
	}

	return values
}

func (f *Field) parseToBools(value string) ([]bool, error) {
	list := strings.Split(value, ",")
	values := make([]bool, len(list))

	for idx, elem := range list {
		parsedValue, err := strconv.ParseBool(strings.Trim(elem, " "))
		if err != nil {
			return nil, err
		}

		values[idx] = parsedValue
	}

	return values, nil
}

func (f *Field) parseToInts(value string) ([]int, error) {
	list := strings.Split(value, ",")
	values := make([]int, len(list))

	for idx, elem := range list {
		parsedValue, err := strconv.Atoi(strings.Trim(elem, " "))
		if err != nil {
			return nil, err
		}

		values[idx] = parsedValue
	}

	return values, nil
}

func (f *Field) parseToTimes(value string) ([]time.Time, error) {
	list := strings.Split(value, ",")
	values := make([]time.Time, len(list))

	for idx, elem := range list {
		parsedValue, err := time.Parse(time.RFC3339, strings.Trim(elem, " "))
		if err != nil {
			return nil, err
		}

		values[idx] = parsedValue
	}

	return values, nil
}

func (f *Field) parseToDurations(value string) ([]time.Duration, error) {
	list := strings.Split(value, ",")
	values := make([]time.Duration, len(list))

	for idx, elem := range list {
		parsedValue, err := time.ParseDuration(strings.Trim(elem, " "))
		if err != nil {
			return nil, err
		}

		values[idx] = parsedValue
	}

	return values, nil
}

func (f *Field) valueInEnumeration(value any) bool {
	if len(f.Enumeration) < 1 {
		return true
	}

	for _, acceptedValue := range f.Enumeration {
		if value == acceptedValue {
			return true
		}
	}

	return false
}

func (f *Field) enumerationString() string {
	builder := strings.Builder{}

	if len(f.Enumeration) > 0 {
		builder.WriteString("[")

		for index, value := range f.Enumeration {
			if index != 0 {
				builder.WriteString(", ")
			}

			builder.WriteString(fmt.Sprintf("'%s'", value))
		}

		builder.WriteString("]")
	}

	return builder.String()
}
