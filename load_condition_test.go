package bconf_test

import (
	"testing"
	"time"

	"github.com/xavi-group/bconf"
)

func TestLoadConditionClone(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).AddFieldSetDependencies("app", "mode").Create()

	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "mode",
		FieldValue:  "production",
	})

	clone := condition.Clone()

	// Verify clone has same dependencies
	deps := clone.FieldDependencies()
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}

	// Verify clone has same field values
	val, found := clone.GetFieldValue("app", "mode")
	if !found {
		t.Fatal("expected to find field value in clone")
	}

	if val != "production" {
		t.Fatalf("expected 'production', got '%v'", val)
	}

	// Modify original, verify clone is independent
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "mode",
		FieldValue:  "development",
	})

	val, _ = clone.GetFieldValue("app", "mode")
	if val != "production" {
		t.Fatal("clone should be independent from original")
	}
}

func TestLoadConditionGetFieldDependencies(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	condition.SetFieldValues(
		bconf.FieldValue{FieldSetKey: "app", FieldKey: "mode", FieldValue: "prod"},
		bconf.FieldValue{FieldSetKey: "log", FieldKey: "level", FieldValue: "info"},
	)

	deps := condition.GetFieldDependencies()
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
}

//nolint:dupl // Test code intentionally follows similar patterns for each type
func TestLoadConditionGetString(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetString("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != "" {
		t.Fatalf("expected zero value, got '%s'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "myapp",
	})

	val, found, err = condition.GetString("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "myapp" {
		t.Fatalf("expected 'myapp', got '%s'", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "port",
		FieldValue:  8080,
	})

	_, found, err = condition.GetString("app", "port")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetStrings(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetStrings("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got '%v'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "tags",
		FieldValue:  []string{"api", "backend"},
	})

	val, found, err = condition.GetStrings("app", "tags")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(val) != 2 || val[0] != "api" {
		t.Fatalf("unexpected value: %v", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "not-a-slice",
	})

	_, found, err = condition.GetStrings("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

//nolint:dupl // Test code intentionally follows similar patterns for each type
func TestLoadConditionGetInt(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetInt("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected zero value, got '%d'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "port",
		FieldValue:  8080,
	})

	val, found, err = condition.GetInt("app", "port")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 8080 {
		t.Fatalf("expected 8080, got '%d'", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "not-an-int",
	})

	_, found, err = condition.GetInt("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetInts(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetInts("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got '%v'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "ports",
		FieldValue:  []int{8080, 8081},
	})

	val, found, err = condition.GetInts("app", "ports")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(val) != 2 || val[0] != 8080 {
		t.Fatalf("unexpected value: %v", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "port",
		FieldValue:  8080,
	})

	_, found, err = condition.GetInts("app", "port")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetBool(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetBool("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != false {
		t.Fatalf("expected zero value, got '%t'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "debug",
		FieldValue:  true,
	})

	val, found, err = condition.GetBool("app", "debug")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != true {
		t.Fatalf("expected true, got '%t'", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "not-a-bool",
	})

	_, found, err = condition.GetBool("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetBools(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetBools("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got '%v'", val)
	}

	// Test success
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "flags",
		FieldValue:  []bool{true, false, true},
	})

	val, found, err = condition.GetBools("app", "flags")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(val) != 3 || val[0] != true {
		t.Fatalf("unexpected value: %v", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "debug",
		FieldValue:  true,
	})

	_, found, err = condition.GetBools("app", "debug")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetTime(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetTime("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if !val.IsZero() {
		t.Fatalf("expected zero value, got '%v'", val)
	}

	// Test success
	now := time.Now()
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "start_time",
		FieldValue:  now,
	})

	val, found, err = condition.GetTime("app", "start_time")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !val.Equal(now) {
		t.Fatalf("expected '%v', got '%v'", now, val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "not-a-time",
	})

	_, found, err = condition.GetTime("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetTimes(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetTimes("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got '%v'", val)
	}

	// Test success
	times := []time.Time{time.Now(), time.Now().Add(time.Hour)}
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "timestamps",
		FieldValue:  times,
	})

	val, found, err = condition.GetTimes("app", "timestamps")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(val) != 2 {
		t.Fatalf("expected 2 times, got %d", len(val))
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "start_time",
		FieldValue:  time.Now(),
	})

	_, found, err = condition.GetTimes("app", "start_time")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetDuration(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetDuration("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected zero value, got '%v'", val)
	}

	// Test success
	duration := 5 * time.Second
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "timeout",
		FieldValue:  duration,
	})

	val, found, err = condition.GetDuration("app", "timeout")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != duration {
		t.Fatalf("expected '%v', got '%v'", duration, val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "name",
		FieldValue:  "not-a-duration",
	})

	_, found, err = condition.GetDuration("app", "name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestLoadConditionGetDurations(t *testing.T) {
	condition := bconf.LCB(func(_ bconf.FieldValueFinder) (bool, error) {
		return true, nil
	}).Create()

	// Test not found
	val, found, err := condition.GetDurations("app", "missing")
	if found {
		t.Fatal("expected field not to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error for missing field: %v", err)
	}
	if val != nil {
		t.Fatalf("expected nil, got '%v'", val)
	}

	// Test success
	durations := []time.Duration{5 * time.Second, 10 * time.Second}
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "timeouts",
		FieldValue:  durations,
	})

	val, found, err = condition.GetDurations("app", "timeouts")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(val) != 2 || val[0] != 5*time.Second {
		t.Fatalf("unexpected value: %v", val)
	}

	// Test type error
	condition.SetFieldValues(bconf.FieldValue{
		FieldSetKey: "app",
		FieldKey:    "timeout",
		FieldValue:  5 * time.Second,
	})

	_, found, err = condition.GetDurations("app", "timeout")
	if !found {
		t.Fatal("expected field to be found")
	}
	if err == nil {
		t.Fatal("expected type casting error")
	}
}

func TestFieldLocationHelper(t *testing.T) {
	loc := bconf.FD("app", "mode")

	if loc.FieldSetKey != "app" {
		t.Fatalf("expected FieldSetKey 'app', got '%s'", loc.FieldSetKey)
	}

	if loc.FieldKey != "mode" {
		t.Fatalf("expected FieldKey 'mode', got '%s'", loc.FieldKey)
	}
}
