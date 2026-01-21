package bconf

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// FieldValidator is the function signature for field validation.
type FieldValidator func(value any) error

// ValidateAll returns a validator that passes only if all provided validators pass.
func ValidateAll(validators ...FieldValidator) FieldValidator {
	return func(value any) error {
		for _, v := range validators {
			if err := v(value); err != nil {
				return err
			}
		}

		return nil
	}
}

// ValidateAny returns a validator that passes if at least one validator passes.
// If all validators fail, returns the error from the last validator.
func ValidateAny(validators ...FieldValidator) FieldValidator {
	return func(value any) error {
		var lastErr error

		for _, v := range validators {
			if err := v(value); err == nil {
				return nil
			} else {
				lastErr = err
			}
		}

		return lastErr
	}
}

// String Validators

// ValidateNonEmptyString returns a validator that ensures a string is not empty.
func ValidateNonEmptyString() FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("string cannot be empty")
		}

		return nil
	}
}

// ValidateStringMinLength returns a validator that ensures a string has at least min characters.
func ValidateStringMinLength(min int) FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		if len(s) < min {
			return fmt.Errorf("string length %d is less than minimum %d", len(s), min)
		}

		return nil
	}
}

// ValidateStringMaxLength returns a validator that ensures a string has at most max characters.
func ValidateStringMaxLength(max int) FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		if len(s) > max {
			return fmt.Errorf("string length %d exceeds maximum %d", len(s), max)
		}

		return nil
	}
}

// ValidateStringRegex returns a validator that ensures a string matches the given regex pattern.
func ValidateStringRegex(pattern string) FieldValidator {
	re := regexp.MustCompile(pattern)

	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		if !re.MatchString(s) {
			return fmt.Errorf("string '%s' does not match pattern '%s'", s, pattern)
		}

		return nil
	}
}

// ValidateURL returns a validator that ensures a string is a valid URL.
func ValidateURL() FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		u, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("URL must have scheme and host: '%s'", s)
		}

		return nil
	}
}

// ValidateURLWithSchemes returns a validator that ensures a string is a valid URL with one of the allowed schemes.
func ValidateURLWithSchemes(schemes ...string) FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		u, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("URL must have scheme and host: '%s'", s)
		}

		for _, scheme := range schemes {
			if strings.EqualFold(u.Scheme, scheme) {
				return nil
			}
		}

		return fmt.Errorf("URL scheme '%s' not in allowed schemes %v", u.Scheme, schemes)
	}
}

// Int Validators

// ValidateIntMin returns a validator that ensures an int is at least min.
func ValidateIntMin(min int) FieldValidator {
	return func(value any) error {
		i, ok := value.(int)
		if !ok {
			return fmt.Errorf("expected int, got %T", value)
		}

		if i < min {
			return fmt.Errorf("value %d is less than minimum %d", i, min)
		}

		return nil
	}
}

// ValidateIntMax returns a validator that ensures an int is at most max.
func ValidateIntMax(max int) FieldValidator {
	return func(value any) error {
		i, ok := value.(int)
		if !ok {
			return fmt.Errorf("expected int, got %T", value)
		}

		if i > max {
			return fmt.Errorf("value %d exceeds maximum %d", i, max)
		}

		return nil
	}
}

// ValidateIntRange returns a validator that ensures an int is within [min, max].
func ValidateIntRange(min, max int) FieldValidator {
	return func(value any) error {
		i, ok := value.(int)
		if !ok {
			return fmt.Errorf("expected int, got %T", value)
		}

		if i < min || i > max {
			return fmt.Errorf("value %d is not in range [%d, %d]", i, min, max)
		}

		return nil
	}
}

// ValidatePort returns a validator that ensures an int is a valid port number (1-65535).
func ValidatePort() FieldValidator {
	return ValidateIntRange(1, 65535)
}

// File/Path Validators

// ValidateFileExists returns a validator that ensures a file exists at the given path.
func ValidateFileExists() FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		info, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: '%s'", s)
			}

			return fmt.Errorf("error checking file: %w", err)
		}

		if info.IsDir() {
			return fmt.Errorf("path is a directory, not a file: '%s'", s)
		}

		return nil
	}
}

// ValidateDirExists returns a validator that ensures a directory exists at the given path.
func ValidateDirExists() FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		info, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("directory does not exist: '%s'", s)
			}

			return fmt.Errorf("error checking directory: %w", err)
		}

		if !info.IsDir() {
			return fmt.Errorf("path is a file, not a directory: '%s'", s)
		}

		return nil
	}
}

// ValidatePathExists returns a validator that ensures a file or directory exists at the given path.
func ValidatePathExists() FieldValidator {
	return func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		if _, err := os.Stat(s); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: '%s'", s)
			}

			return fmt.Errorf("error checking path: %w", err)
		}

		return nil
	}
}

// Slice Validators

// ValidateNonEmptySlice returns a validator that ensures a slice is not empty.
// Works with []string, []int, []bool, and other slice types.
func ValidateNonEmptySlice() FieldValidator {
	return func(value any) error {
		switch v := value.(type) {
		case []string:
			if len(v) == 0 {
				return fmt.Errorf("slice cannot be empty")
			}
		case []int:
			if len(v) == 0 {
				return fmt.Errorf("slice cannot be empty")
			}
		case []bool:
			if len(v) == 0 {
				return fmt.Errorf("slice cannot be empty")
			}
		default:
			return fmt.Errorf("expected slice type, got %T", value)
		}

		return nil
	}
}

// ValidateSliceMinLength returns a validator that ensures a slice has at least min elements.
func ValidateSliceMinLength(min int) FieldValidator {
	return func(value any) error {
		var length int

		switch v := value.(type) {
		case []string:
			length = len(v)
		case []int:
			length = len(v)
		case []bool:
			length = len(v)
		default:
			return fmt.Errorf("expected slice type, got %T", value)
		}

		if length < min {
			return fmt.Errorf("slice length %d is less than minimum %d", length, min)
		}

		return nil
	}
}

// ValidateSliceMaxLength returns a validator that ensures a slice has at most max elements.
func ValidateSliceMaxLength(max int) FieldValidator {
	return func(value any) error {
		var length int

		switch v := value.(type) {
		case []string:
			length = len(v)
		case []int:
			length = len(v)
		case []bool:
			length = len(v)
		default:
			return fmt.Errorf("expected slice type, got %T", value)
		}

		if length > max {
			return fmt.Errorf("slice length %d exceeds maximum %d", length, max)
		}

		return nil
	}
}

// ValidateEachString returns a validator that applies the given validator to each element in a []string.
func ValidateEachString(validator FieldValidator) FieldValidator {
	return func(value any) error {
		slice, ok := value.([]string)
		if !ok {
			return fmt.Errorf("expected []string, got %T", value)
		}

		for i, elem := range slice {
			if err := validator(elem); err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}
		}

		return nil
	}
}

// ValidateEachInt returns a validator that applies the given validator to each element in a []int.
func ValidateEachInt(validator FieldValidator) FieldValidator {
	return func(value any) error {
		slice, ok := value.([]int)
		if !ok {
			return fmt.Errorf("expected []int, got %T", value)
		}

		for i, elem := range slice {
			if err := validator(elem); err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}
		}

		return nil
	}
}
