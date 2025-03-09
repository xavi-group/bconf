package bconf

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
)

// NewAppConfig creates a new application configuration struct with configuration options that allow for the
// specification of configuration sources (environment, flags, json files, etc).
func NewAppConfig(appName, appDescription string, options ...ConfigOption) *AppConfig {
	warnings := []string{}
	loaders := []Loader{}

	appVersion := "unknown"
	appID := "undefined"

	var (
		appIDGenerator      func() (any, error)
		appVersionGenerator func() (any, error)
	)

	for _, option := range options {
		switch option.ConfigOptionType() {
		case configOptionTypeLoaderEnvironment:
			if castOption, ok := option.(configOptionEnvironmentLoader); ok {
				loaders = append(loaders, castOption.Loader())
			} else {
				warnings = append(warnings, "problem casting environment loader option")
			}
		case configOptionTypeLoaderFlag:
			if castOption, ok := option.(configOptionFlagLoader); ok {
				loaders = append(loaders, castOption.Loader())
			} else {
				warnings = append(warnings, "problem casting flag loader option")
			}
		case configOptionTypeLoaderJSONFile:
			if castOption, ok := option.(*configOptionJSONFileLoader); ok {
				loaders = append(loaders, castOption.Loader())
			} else {
				warnings = append(warnings, "problem casting JSON-file loader option")
			}
		case configOptionTypeAppID:
			if castOption, ok := option.(configOptionAppID); ok {
				appID = castOption.id
			} else {
				warnings = append(warnings, "problem casting app ID option")
			}
		case configOptionTypeAppIDFunc:
			if castOption, ok := option.(configOptionAppIDFunc); ok {
				appIDGenerator = func() (any, error) {
					return castOption.idFunc(), nil
				}
			} else {
				warnings = append(warnings, "problem casting app ID func option")
			}
		case configOptionTypeAppVersion:
			if castOption, ok := option.(configOptionAppVersion); ok {
				appVersion = castOption.version
			} else {
				warnings = append(warnings, "problem casting app version option")
			}
		case configOptionTypeAppVersionFunc:
			if castOption, ok := option.(configOptionAppVersionFunc); ok {
				appVersionGenerator = func() (any, error) {
					return castOption.versionFunc(), nil
				}
			} else {
				warnings = append(warnings, "problem casting app version func option")
			}
		default:
			warnings = append(warnings, fmt.Sprintf("unsupported config option '%s'", option.ConfigOptionType()))
		}
	}

	var (
		appIDField      *Field
		appVersionField *Field
	)

	if appVersionGenerator != nil {
		appVersionField = FB("version", String).DefaultGenerator(appVersionGenerator).C()
	} else {
		appVersionField = FB("version", String).Default(appVersion).C()
	}

	if appIDGenerator != nil {
		appIDField = FB("id", String).DefaultGenerator(appIDGenerator).C()
	} else {
		appIDField = FB("id", String).Default(appID).C()
	}

	appFieldSet := FSB("app").Fields(
		FB("name", String).Default(appName).C(),
		FB("description", String).Default(appDescription).C(),
		appVersionField,
		appIDField,
	).C()

	config := &AppConfig{
		fieldSetGroups:   fieldSetGroups{},
		fieldSets:        map[string]*FieldSet{},
		fillStructs:      []any{},
		loaders:          loaders,
		warnings:         warnings,
		orderedFieldSets: FieldSets{},
	}

	config.AddFieldSet(appFieldSet)

	return config
}

// --------------------------------------------------------------------------------------------------------------------

// AppConfig manages application configuration field-sets and provides access to configuration values. It should be
// initialized with the NewAppConfig function.
type AppConfig struct {
	fieldSets        map[string]*FieldSet
	fieldSetGroups   fieldSetGroups
	loaders          []Loader
	fillStructs      []any
	warnings         []string
	orderedFieldSets FieldSets
	fieldSetLock     sync.Mutex
	loaded           bool
}

func (c *AppConfig) AppName() string {
	name, _ := c.GetString("app", "name")

	return name
}

func (c *AppConfig) AppDescription() string {
	description, _ := c.GetString("app", "description")

	return description
}

func (c *AppConfig) AppVersion() string {
	version, _ := c.GetString("app", "version")

	return version
}

func (c *AppConfig) AppID() string {
	id, _ := c.GetString("app", "id")

	return id
}

func (c *AppConfig) AddFieldSetGroup(groupName string, fieldSets FieldSets) {
	c.fieldSetGroups = append(c.fieldSetGroups, &fieldSetGroup{name: groupName, fieldSets: fieldSets})
}

func (c *AppConfig) AttachConfigStructs(configStructs ...any) {
	c.fillStructs = append(c.fillStructs, configStructs...)
}

func (c *AppConfig) AddFieldSet(fieldSet *FieldSet) {
	c.fieldSetGroups = append(c.fieldSetGroups, &fieldSetGroup{name: fieldSet.Key, fieldSets: FieldSets{fieldSet}})
}

func (c *AppConfig) GetField(fieldSetKey, fieldKey string) (*Field, error) {
	fieldSet, found := c.fieldSets[fieldSetKey]
	if !found {
		return nil, fmt.Errorf("field-set not found with key '%s'", fieldSetKey)
	}

	field, found := fieldSet.fieldMap[fieldKey]
	if !found {
		return nil, fmt.Errorf("field-set field not found with key '%s'", fieldKey)
	}

	return field, nil
}

func (c *AppConfig) SetField(fieldSetKey, fieldKey string, fieldValue any) error {
	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]
	if !fieldSetFound {
		return fmt.Errorf("field-set with key '%s' not found", fieldSetKey)
	}

	field, fieldKeyFound := fieldSet.fieldMap[fieldKey]
	if !fieldKeyFound {
		return fmt.Errorf("field with key '%s' not found", fieldKey)
	}

	if err := field.setOverride(fieldValue); err != nil {
		return fmt.Errorf("problem setting field value: %w", err)
	}

	return nil
}

func (c *AppConfig) GetString(fieldSetKey, fieldKey string) (string, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, String)
	if err != nil {
		return "", err
	}

	val, _ := fieldValue.(string)

	return val, nil
}

func (c *AppConfig) GetStrings(fieldSetKey, fieldKey string) ([]string, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Strings)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]string)

	return val, nil
}

func (c *AppConfig) GetInt(fieldSetKey, fieldKey string) (int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Int)
	if err != nil {
		return 0, err
	}

	val, _ := fieldValue.(int)

	return val, nil
}

func (c *AppConfig) GetInts(fieldSetKey, fieldKey string) ([]int, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Ints)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]int)

	return val, nil
}

func (c *AppConfig) GetBool(fieldSetKey, fieldKey string) (bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Bool)
	if err != nil {
		return false, err
	}

	val, _ := fieldValue.(bool)

	return val, nil
}

func (c *AppConfig) GetBools(fieldSetKey, fieldKey string) ([]bool, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Bools)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]bool)

	return val, nil
}

func (c *AppConfig) GetTime(fieldSetKey, fieldKey string) (time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Time)
	if err != nil {
		return time.Time{}, err
	}

	val, _ := fieldValue.(time.Time)

	return val, nil
}

func (c *AppConfig) GetTimes(fieldSetKey, fieldKey string) ([]time.Time, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Times)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]time.Time)

	return val, nil
}

func (c *AppConfig) GetDuration(fieldSetKey, fieldKey string) (time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Duration)
	if err != nil {
		return 0, err
	}

	val, _ := fieldValue.(time.Duration)

	return val, nil
}

func (c *AppConfig) GetDurations(fieldSetKey, fieldKey string) ([]time.Duration, error) {
	fieldValue, err := c.getFieldValue(fieldSetKey, fieldKey, Durations)
	if err != nil {
		return nil, err
	}

	val, _ := fieldValue.([]time.Duration)

	return val, nil
}

func (c *AppConfig) Load(options ...LoadOption) []error {
	// -- Add field set groups --
	groupAddErrors := []error{}

	for _, group := range c.fieldSetGroups {
		if errs := c.addFieldSets(group.fieldSets...); len(errs) > 0 {
			err := fmt.Errorf("problem(s) adding '%s' field-set group: %v", group.name, errs)

			groupAddErrors = append(groupAddErrors, err)
		}
	}

	if len(groupAddErrors) > 0 {
		return groupAddErrors
	}

	// -- Parse load options --

	handleHelpFlag := true

	for _, option := range options {
		switch option.LoadOptionType() {
		case loadOptionTypeDisableHelpFlag:
			handleHelpFlag = false
		default:
			c.warnings = append(c.warnings, fmt.Sprintf("unsupported load option '%s'", option.LoadOptionType()))
		}
	}

	// -- Output help message if conditions are satisfied --

	if handleHelpFlag && len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		c.printHelpString()
		os.Exit(0)
	}

	// -- Load field-sets --

	loadErrors := []error{}

	for _, fieldSet := range c.orderedFieldSets {
		if fieldSetErrs := c.loadFieldSet(fieldSet.Key); len(fieldSetErrs) > 0 {
			loadErrors = append(loadErrors, fieldSetErrs...)
			return loadErrors
		}
	}

	fillErrors := []error{}

	for _, fillStruct := range c.fillStructs {
		if fillErr := c.FillStruct(fillStruct); fillErr != nil {
			fillErrors = append(fillErrors, fillErr)
		}
	}

	if len(fillErrors) > 0 {
		return fillErrors
	}

	c.loaded = true

	return nil
}

func (c *AppConfig) FillStruct(configStruct any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("problem filling struct: %s", r)
		}
	}()

	configStructKind := reflect.TypeOf(configStruct).Kind()
	if configStructKind != reflect.Pointer {
		return fmt.Errorf("FillStruct expects a pointer to a struct, found '%s'", configStructKind)
	}

	configStructValue := reflect.Indirect(reflect.ValueOf(configStruct))
	configStructPointerToKind := configStructValue.Kind()
	configStructType := configStructValue.Type()

	if configStructPointerToKind != reflect.Struct {
		return fmt.Errorf(
			"FillStruct expects a pointer to a struct, found pointer to '%s'",
			configStructPointerToKind,
		)
	}

	baseFieldSetFound := false
	baseFieldSet := ""

	configStructField := configStructValue.FieldByName("ConfigStruct")

	if configStructField.IsValid() && configStructField.Type().PkgPath() == "github.com/xavi-group/bconf" {
		var configStructFieldType reflect.StructField

		configStructFieldType, baseFieldSetFound = configStructType.FieldByName("ConfigStruct")

		if baseFieldSetFound {
			baseFieldSet = configStructFieldType.Tag.Get("bconf")

			if overrideValue := configStructField.FieldByName("FieldSet"); overrideValue.String() != "" {
				baseFieldSet = overrideValue.String()
			}
		}
	}

	for i := 0; i < configStructValue.NumField(); i++ {
		field := configStructType.Field(i)

		if field.Name == "ConfigStruct" && field.Type.PkgPath() == "github.com/xavi-group/bconf" {
			continue
		}

		if field.Type.Kind() == reflect.Pointer && reflect.Indirect(reflect.ValueOf(field)).Kind() == reflect.Struct {
			fieldValue := configStructValue.Field(i)

			if fieldValue.IsNil() {
				configStructValue.Field(i).Set(reflect.New(field.Type.Elem()))
			}

			if err := c.FillStruct(configStructValue.Field(i).Interface()); err != nil {
				return fmt.Errorf("problem filling struct field: %w", err)
			}

			continue
		}

		fieldTagValue := field.Tag.Get("bconf")
		fieldKey := ""
		fieldSetKey := baseFieldSet

		switch fieldTagValue {
		case "":
			fieldKey = field.Name
		case "-":
			continue
		default:
			fieldTagParams := strings.Split(fieldTagValue, ",")
			fieldLocation := strings.Split(fieldTagParams[0], ".")

			fieldKey = fieldLocation[0]

			// NOTE: error if fieldLocation format isn't <field>.<field-name> ?
			if len(fieldLocation) > 1 {
				fieldSetKey = fieldLocation[0]
				fieldKey = fieldLocation[1]
			}
		}

		if fieldSetKey == "" {
			return fmt.Errorf("unidentified field-set for field: %s", fieldKey)
		}

		appConfigField, err := c.GetField(fieldSetKey, fieldKey)
		if err != nil {
			return fmt.Errorf("problem getting field '%s.%s': %w", fieldSetKey, fieldKey, err)
		}

		val, err := appConfigField.getValue()
		if err != nil && err.Error() == emptyFieldError {
			continue
		} else if err != nil {
			return fmt.Errorf("problem getting field '%s.%s' value: %w", fieldSetKey, fieldKey, err)
		}

		configStructValue.Field(i).Set(reflect.ValueOf(val))
	}

	return nil
}

func (c *AppConfig) ConfigMap() map[string]map[string]any {
	configMap := map[string]map[string]any{}

	for _, fieldSet := range c.fieldSets {
		fieldSetMap := map[string]any{}

		for _, field := range fieldSet.fieldMap {
			location := fmt.Sprintf("%s.%s", fieldSet.Key, field.Key)

			switch location {
			case "app.name":
			case "app.description":
				continue
			}

			val, err := field.getValue()

			if err != nil {
				continue
			}

			if field.Sensitive {
				fieldSetMap[field.Key] = "<sensitive-value>"
				continue
			}

			if field.Type == Duration {
				// TODO: use higher time grain if duration > 1 hour ?
				valDuration, _ := val.(time.Duration)
				valMilliseconds := valDuration.Milliseconds()
				fieldSetMap[fmt.Sprintf("%s_ms", field.Key)] = valMilliseconds

				continue
			}

			fieldSetMap[field.Key] = val
		}

		configMap[fieldSet.Key] = fieldSetMap
	}

	return configMap
}

func (c *AppConfig) Warnings() []string {
	return slices.Clone(c.warnings)
}

func (c *AppConfig) HelpString() string {
	maxCharLength := 100
	builder := strings.Builder{}

	name := c.AppName()
	description := c.AppDescription()

	if name != "" {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", name))
	} else {
		builder.WriteString(fmt.Sprintf("Usage of '%s':\n", os.Args[0]))
	}

	if description != "" && len(description) > maxCharLength {
		wrapStringForBuilder(description, &builder, maxCharLength, "")
	} else if description != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", description))
	}

	c.addFieldsToBuilder(&builder, maxCharLength)

	return builder.String()
}

// --------------------------------------------------------------------------------------------------------------------

func (c *AppConfig) addFieldSets(fieldSets ...*FieldSet) []error {
	c.fieldSetLock.Lock()
	defer c.fieldSetLock.Unlock()

	errs := []error{}
	addedFieldSets := []string{}

	for _, fieldSet := range fieldSets {
		if fieldSetErrs := c.addFieldSet(fieldSet, false); len(fieldSetErrs) > 0 {
			errs = append(errs, fieldSetErrs...)
			continue
		}

		addedFieldSets = append(addedFieldSets, fieldSet.Key)
	}

	if len(errs) > 0 {
		for _, fieldSetKey := range addedFieldSets {
			delete(c.fieldSets, fieldSetKey)
		}

		c.orderedFieldSets = c.orderedFieldSets[:len(c.orderedFieldSets)-len(addedFieldSets)]
	}

	return errs
}

func (c *AppConfig) addFieldSet(fieldSet *FieldSet, lock bool) []error {
	if lock {
		c.fieldSetLock.Lock()
		defer c.fieldSetLock.Unlock()
	}

	fieldSet = fieldSet.Clone()

	if errs := c.checkForFieldSetStructuralIntegrity(fieldSet); len(errs) > 0 {
		return errs
	}

	if _, keyFound := c.fieldSets[fieldSet.Key]; keyFound {
		return []error{fmt.Errorf("duplicate field-set key found: '%s'", fieldSet.Key)}
	}

	fieldSet.initializeFieldMap()

	if errs := c.checkForFieldSetDependencies(fieldSet); len(errs) > 0 {
		return errs
	}

	if errs := c.generateFieldSetDefaultValues(fieldSet); len(errs) > 0 {
		return errs
	}

	if errs := c.checkForFieldSetFieldsValidity(fieldSet); len(errs) > 0 {
		return errs
	}

	fieldSet.Fields = nil

	c.fieldSets[fieldSet.Key] = fieldSet
	c.orderedFieldSets = append(c.orderedFieldSets, fieldSet)

	return nil
}

func (c *AppConfig) checkForFieldSetStructuralIntegrity(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.validate(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(errs, fmt.Errorf("field-set '%s' validation error: %w", fieldSet.Key, err))
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldSetDependencies(fieldSet *FieldSet) []error {
	errs := []error{}

	for _, loadCondition := range fieldSet.LoadConditions {
		dependencies := loadCondition.FieldDependencies()
		if len(dependencies) < 1 {
			continue
		}

		for _, dependency := range dependencies {
			fieldSetDependency, found := c.fieldSets[dependency.FieldSetKey]
			if !found {
				errs = append(
					errs,
					fmt.Errorf(
						"field-set '%s' (for field-set '%s' load condition) not found in config",
						dependency.FieldSetKey, fieldSet.Key,
					),
				)

				continue
			}

			_, found = fieldSetDependency.fieldMap[dependency.FieldKey]
			if !found {
				errs = append(
					errs,
					fmt.Errorf(
						"field-set '%s' field '%s' (for field-set '%s' load condition) not found in config",
						dependency.FieldSetKey, dependency.FieldKey, fieldSet.Key,
					),
				)
			}
		}
	}

	for _, field := range fieldSet.Fields {
		if err := c.checkForFieldDependencies(field, fieldSet); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldDependencies(field *Field, parent *FieldSet) error {
	for _, loadCondition := range field.LoadConditions {
		dependencies := loadCondition.FieldDependencies()

		if len(dependencies) < 1 {
			continue
		}

		for _, dependency := range dependencies {
			var (
				fieldSetDependency *FieldSet
				found              bool
			)

			if dependency.FieldSetKey == "" || dependency.FieldSetKey == parent.Key {
				fieldSetDependency = parent
			} else {
				fieldSetDependency, found = c.fieldSets[dependency.FieldSetKey]

				if !found {
					return fmt.Errorf(
						"field-set '%s' (for field-set '%s' field '%s' load condition) not found in config",
						dependency.FieldSetKey, parent.Key, field.Key,
					)
				}
			}

			if _, found = fieldSetDependency.fieldMap[dependency.FieldKey]; !found {
				return fmt.Errorf(
					"field-set '%s' field '%s' (for field-set '%s' field '%s' load condition) not found in config",
					fieldSetDependency.Key, dependency.FieldKey, parent.Key, field.Key,
				)
			}
		}
	}

	return nil
}

func (c *AppConfig) generateFieldSetDefaultValues(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.generateFieldDefaults(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field default value generation error: %w", fieldSet.Key, err),
			)
		}
	}

	return errs
}

func (c *AppConfig) checkForFieldSetFieldsValidity(fieldSet *FieldSet) []error {
	errs := []error{}

	if fieldSetErrs := fieldSet.validateFields(); len(fieldSetErrs) > 0 {
		for _, err := range fieldSetErrs {
			errs = append(
				errs,
				fmt.Errorf("field-set '%s' field validation error: %w", fieldSet.Key, err),
			)
		}
	}

	return errs
}

func (c *AppConfig) loadFieldSet(fieldSetKey string) []error {
	errs := []error{}

	fieldSet, fieldSetFound := c.fieldSets[fieldSetKey]
	if !fieldSetFound {
		errs = append(errs, fmt.Errorf("field-set with key '%s' not found", fieldSetKey))
		return errs
	}

	if load, err := c.shouldLoadFieldSet(fieldSet); err != nil {
		return append(errs, err)
	} else if !load {
		return errs
	}

	for _, loader := range c.loaders {
		values := loader.GetMap(fieldSetKey, c.fieldSets[fieldSetKey].fieldKeys())
		for key, value := range values {
			field := c.fieldSets[fieldSetKey].fieldMap[key]

			if load, err := c.shouldLoadField(field, fieldSetKey); err != nil {
				errs = append(errs, err)
				continue
			} else if !load {
				continue
			}

			if err := c.fieldSets[fieldSetKey].fieldMap[key].set(loader.Name(), value); err != nil {
				errs = append(errs, fmt.Errorf("field '%s' load error: %w", key, err))
			}
		}
	}

	for _, field := range fieldSet.fieldMap {
		if field.Required && len(field.LoadConditions) < 1 {
			if _, err := field.getValue(); err != nil {
				errs = append(errs, fmt.Errorf("required field '%s_%s' not set", fieldSet.Key, field.Key))
			}
		} else if field.Required {
			if load, _ := c.shouldLoadField(field, fieldSet.Key); load {
				if _, err := field.getValue(); err != nil {
					errs = append(errs, fmt.Errorf(
						"conditionally required field '%s_%s' load condition met, but field value not set",
						fieldSet.Key,
						field.Key,
					))
				}
			}
		}
	}

	return errs
}

func (c *AppConfig) shouldLoadFieldSet(fieldSet *FieldSet) (loadFieldSet bool, err error) {
	loadFieldSet = true

	for _, loadCondition := range fieldSet.LoadConditions {
		if !loadFieldSet {
			break
		}

		dependencies := loadCondition.FieldDependencies()
		if len(dependencies) < 1 {
			loadFieldSet, err = loadCondition.Load(loadCondition)
			if err != nil {
				return false, fmt.Errorf("problem getting field-set '%s' load condition outcome: %w", fieldSet.Key, err)
			}

			continue
		}

		for _, dependency := range dependencies {
			fieldValue, err := c.getFieldValue(dependency.FieldSetKey, dependency.FieldKey, "any")
			if err != nil {
				return false, fmt.Errorf(
					"problem getting field value for field-set '%s' load condition: %w", fieldSet.Key, err,
				)
			}

			loadCondition.SetFieldValues(FieldValue{
				FieldSetKey: dependency.FieldSetKey, FieldKey: dependency.FieldKey, FieldValue: fieldValue,
			})
		}

		loadFieldSet, err = loadCondition.Load(loadCondition)
		if err != nil {
			return false, fmt.Errorf("problem getting field-set '%s' load condition outcome: %w", fieldSet.Key, err)
		}
	}

	return loadFieldSet, nil
}

func (c *AppConfig) shouldLoadField(field *Field, fieldSetKey string) (loadField bool, err error) {
	loadField = true

	for _, loadCondition := range field.LoadConditions {
		if !loadField {
			break
		}

		dependencies := loadCondition.FieldDependencies()
		if len(dependencies) < 1 {
			loadField, err = loadCondition.Load(loadCondition)
			if err != nil {
				return false, fmt.Errorf(
					"problem getting field-set '%s' field '%s' load condition outcome: %w",
					fieldSetKey, field.Key, err,
				)
			}

			continue
		}

		for _, dependency := range dependencies {
			if dependency.FieldSetKey == "" {
				dependency.FieldSetKey = fieldSetKey
			}

			fieldValue, err := c.getFieldValue(dependency.FieldSetKey, dependency.FieldKey, "any")
			if err != nil {
				return false, fmt.Errorf(
					"problem getting field value for field-set '%s' field '%s' load condition: %w",
					fieldSetKey, field.Key, err,
				)
			}

			loadCondition.SetFieldValues(FieldValue{
				FieldSetKey: dependency.FieldSetKey, FieldKey: dependency.FieldKey, FieldValue: fieldValue,
			})
		}

		loadField, err = loadCondition.Load(loadCondition)
		if err != nil {
			return false, fmt.Errorf(
				"problem getting field-set '%s' field '%s' load condition outcome: %w",
				fieldSetKey, field.Key, err,
			)
		}
	}

	return loadField, nil
}

func (c *AppConfig) getFieldValue(fieldSetKey, fieldKey, expectedType string) (any, error) {
	field, err := c.GetField(fieldSetKey, fieldKey)
	if err != nil {
		return nil, err
	}

	if expectedType != "" && expectedType != "any" && field.Type != expectedType {
		return nil, fmt.Errorf("incorrect field-type for field '%s', found '%s'", fieldKey, field.Type)
	}

	fieldValue, err := field.getValue()
	if err != nil {
		return nil, fmt.Errorf("no value set for field '%s'", fieldKey)
	}

	return fieldValue, nil
}

func (c *AppConfig) printHelpString() {
	fmt.Printf("%s", c.HelpString())
}

type fieldEntry struct {
	fieldSetKey    string
	field          *Field
	loadConditions LoadConditions
}

func (c *AppConfig) fields() map[string]*fieldEntry {
	fields := map[string]*fieldEntry{}

	for fieldSetKey, fieldSet := range c.fieldSets {
		for _, field := range fieldSet.fieldMap {
			entry := fieldEntry{field: field, fieldSetKey: fieldSetKey}

			if len(fieldSet.LoadConditions) > 0 {
				entry.loadConditions = fieldSet.LoadConditions
			}

			if len(field.LoadConditions) > 0 {
				entry.loadConditions = append(entry.loadConditions, field.LoadConditions...)
			}

			fields[fmt.Sprintf("%s.%s", fieldSetKey, field.Key)] = &entry
		}
	}

	return fields
}

func (c *AppConfig) addFieldsToBuilder(builder *strings.Builder, maxCharLength int) {
	fields := c.fields()
	if len(fields) > 0 {
		keys := make([]string, len(fields))
		idx := 0

		for key := range fields {
			keys[idx] = key
			idx++
		}

		sort.Strings(keys)

		conditionallyRequiredFields := []string{}
		requiredFields := []string{}
		optionalFields := []string{}

		for _, key := range keys {
			fieldEntry := fields[key]

			switch {
			case fieldEntry.field.Required && fieldEntry.loadConditions == nil:
				requiredFields = append(requiredFields, key)
			case fieldEntry.field.Required && fieldEntry.loadConditions != nil:
				conditionallyRequiredFields = append(conditionallyRequiredFields, key)
			default:
				optionalFields = append(optionalFields, key)
			}
		}

		if len(requiredFields) > 0 {
			builder.WriteString("Required Configuration:\n")

			for _, key := range requiredFields {
				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key, maxCharLength))
			}
		}

		if len(conditionallyRequiredFields) > 0 {
			builder.WriteString("Conditionally Required Configuration:\n")

			for _, key := range conditionallyRequiredFields {
				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key, maxCharLength))
			}
		}

		if len(optionalFields) > 0 {
			builder.WriteString("Optional Configuration:\n")

			for _, key := range optionalFields {
				if key == "app.name" || key == "app.description" {
					continue
				}

				fmt.Fprintf(builder, "\t%s", c.fieldHelpString(fields, key, maxCharLength))
			}
		}
	}
}

func (c *AppConfig) fieldHelpString(fields map[string]*fieldEntry, key string, maxCharLength int) string {
	entry := fields[key]
	field := entry.field
	loadConditions := entry.loadConditions

	if field == nil {
		return "no field matching key"
	}

	builder := strings.Builder{}
	spaceBuffer := "\t\t"

	builder.WriteString(fmt.Sprintf("%s %s\n", key, field.Type))

	if field.Description != "" {
		builder.WriteString(spaceBuffer)
		if len(spaceBuffer)+len(field.Description) > maxCharLength {
			wrapStringForBuilder(field.Description, &builder, maxCharLength, spaceBuffer)
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", field.Description))
		}
	}

	if len(field.Enumeration) > 0 {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Accepted values: %s\n", field.enumerationString()))
	}

	if field.Default != nil && field.Sensitive {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: '<sensitive-value>'\n")
	} else if field.Default != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString(fmt.Sprintf("Default value: '%v'\n", field.Default))
	}

	if field.DefaultGenerator != nil {
		builder.WriteString(spaceBuffer)
		builder.WriteString("Default value: <generated-at-run-time>\n")
	}

	for _, loader := range c.loaders {
		helpString := loader.HelpString(entry.fieldSetKey, entry.field.Key)
		if helpString != "" {
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s\n", helpString))
		}
	}

	for _, condition := range loadConditions {
		dependencies := condition.FieldDependencies()

		if len(dependencies) > 0 {
			builder.WriteString(spaceBuffer)
			builder.WriteString("Loading depends on field(s): ")

			for idx, dependency := range dependencies {
				if idx == 0 {
					builder.WriteString(fmt.Sprintf("'%s.%s'", dependency.FieldSetKey, dependency.FieldKey))
				} else {
					builder.WriteString(fmt.Sprintf(", '%s.%s'", dependency.FieldSetKey, dependency.FieldKey))
				}
			}

			builder.WriteString("\n")
		} else {
			builder.WriteString(spaceBuffer)
			builder.WriteString("Loading depends on: <custom-load-condition-function>\n")
		}
	}

	return builder.String()
}

func wrapStringForBuilder(content string, builder *strings.Builder, maxCharLength int, spaceBuffer string) {
	maxCharLength = maxCharLength - len(spaceBuffer)

	words := strings.Split(content, " ")
	chunkLen := 0

	for _, word := range words {
		wordLen := len(word) + 1

		if chunkLen+wordLen > maxCharLength {
			builder.WriteString("\n")
			builder.WriteString(spaceBuffer)
			builder.WriteString(fmt.Sprintf("%s ", word))

			chunkLen = 0

			continue
		}

		builder.WriteString(fmt.Sprintf("%s ", word))
		chunkLen += wordLen
	}

	builder.WriteString("\n")
}
