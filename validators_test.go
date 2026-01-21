package bconf_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xavi-group/bconf"
)

func TestValidateAll(t *testing.T) {
	validator := bconf.ValidateAll(
		bconf.ValidateNonEmptyString(),
		bconf.ValidateStringMinLength(3),
	)

	// Should pass
	if err := validator("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on empty
	if err := validator(""); err == nil {
		t.Fatal("expected error for empty string")
	}

	// Should fail on too short
	if err := validator("ab"); err == nil {
		t.Fatal("expected error for short string")
	}
}

func TestValidateAny(t *testing.T) {
	validator := bconf.ValidateAny(
		bconf.ValidateStringRegex(`^\d+$`),       // all digits
		bconf.ValidateStringRegex(`^[a-z]+$`),   // all lowercase
	)

	// Should pass (digits)
	if err := validator("123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should pass (lowercase)
	if err := validator("abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail (neither)
	if err := validator("ABC123"); err == nil {
		t.Fatal("expected error for mixed string")
	}
}

func TestValidateNonEmptyString(t *testing.T) {
	validator := bconf.ValidateNonEmptyString()

	// Should pass
	if err := validator("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on empty
	if err := validator(""); err == nil {
		t.Fatal("expected error for empty string")
	}

	// Should fail on whitespace only
	if err := validator("   "); err == nil {
		t.Fatal("expected error for whitespace-only string")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateStringMinLength(t *testing.T) {
	validator := bconf.ValidateStringMinLength(5)

	// Should pass
	if err := validator("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator("hi"); err == nil {
		t.Fatal("expected error for short string")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateStringMaxLength(t *testing.T) {
	validator := bconf.ValidateStringMaxLength(5)

	// Should pass
	if err := validator("hi"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator("hello world"); err == nil {
		t.Fatal("expected error for long string")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateStringRegex(t *testing.T) {
	validator := bconf.ValidateStringRegex(`^[a-z]+$`)

	// Should pass
	if err := validator("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator("Hello"); err == nil {
		t.Fatal("expected error for uppercase string")
	}

	if err := validator("hello123"); err == nil {
		t.Fatal("expected error for alphanumeric string")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateURL(t *testing.T) {
	validator := bconf.ValidateURL()

	// Should pass
	if err := validator("https://example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator("http://localhost:8080/path"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on missing scheme
	if err := validator("example.com"); err == nil {
		t.Fatal("expected error for URL without scheme")
	}

	// Should fail on missing host
	if err := validator("https://"); err == nil {
		t.Fatal("expected error for URL without host")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateURLWithSchemes(t *testing.T) {
	validator := bconf.ValidateURLWithSchemes("https", "wss")

	// Should pass
	if err := validator("https://example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator("wss://example.com/socket"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on disallowed scheme
	if err := validator("http://example.com"); err == nil {
		t.Fatal("expected error for http scheme")
	}

	// Should fail on missing scheme
	if err := validator("example.com"); err == nil {
		t.Fatal("expected error for URL without scheme")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateIntMin(t *testing.T) {
	validator := bconf.ValidateIntMin(10)

	// Should pass
	if err := validator(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator(5); err == nil {
		t.Fatal("expected error for value below minimum")
	}

	// Should fail on wrong type
	if err := validator("10"); err == nil {
		t.Fatal("expected error for non-int type")
	}
}

func TestValidateIntMax(t *testing.T) {
	validator := bconf.ValidateIntMax(100)

	// Should pass
	if err := validator(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator(150); err == nil {
		t.Fatal("expected error for value above maximum")
	}

	// Should fail on wrong type
	if err := validator("100"); err == nil {
		t.Fatal("expected error for non-int type")
	}
}

func TestValidateIntRange(t *testing.T) {
	validator := bconf.ValidateIntRange(10, 100)

	// Should pass
	if err := validator(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator(5); err == nil {
		t.Fatal("expected error for value below range")
	}

	if err := validator(150); err == nil {
		t.Fatal("expected error for value above range")
	}

	// Should fail on wrong type
	if err := validator("50"); err == nil {
		t.Fatal("expected error for non-int type")
	}
}

func TestValidatePort(t *testing.T) {
	validator := bconf.ValidatePort()

	// Should pass
	if err := validator(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator(65535); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator(0); err == nil {
		t.Fatal("expected error for port 0")
	}

	if err := validator(65536); err == nil {
		t.Fatal("expected error for port > 65535")
	}

	if err := validator(-1); err == nil {
		t.Fatal("expected error for negative port")
	}
}

func TestValidateFileExists(t *testing.T) {
	validator := bconf.ValidateFileExists()

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "bconf_test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Should pass for existing file
	if err := validator(tmpFile.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail for non-existent file
	if err := validator("/nonexistent/path/file.txt"); err == nil {
		t.Fatal("expected error for non-existent file")
	}

	// Should fail for directory
	tmpDir, err := os.MkdirTemp("", "bconf_test_dir_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.Remove(tmpDir)

	if err := validator(tmpDir); err == nil {
		t.Fatal("expected error for directory path")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateDirExists(t *testing.T) {
	validator := bconf.ValidateDirExists()

	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "bconf_test_dir_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.Remove(tmpDir)

	// Should pass for existing directory
	if err := validator(tmpDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail for non-existent directory
	if err := validator("/nonexistent/path/dir"); err == nil {
		t.Fatal("expected error for non-existent directory")
	}

	// Should fail for file
	tmpFile, err := os.CreateTemp("", "bconf_test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if err := validator(tmpFile.Name()); err == nil {
		t.Fatal("expected error for file path")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidatePathExists(t *testing.T) {
	validator := bconf.ValidatePathExists()

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "bconf_test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Should pass for file
	if err := validator(tmpFile.Name()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "bconf_test_dir_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.Remove(tmpDir)

	// Should pass for directory
	if err := validator(tmpDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail for non-existent path
	if err := validator(filepath.Join(tmpDir, "nonexistent")); err == nil {
		t.Fatal("expected error for non-existent path")
	}

	// Should fail on wrong type
	if err := validator(123); err == nil {
		t.Fatal("expected error for non-string type")
	}
}

func TestValidateNonEmptySlice(t *testing.T) {
	validator := bconf.ValidateNonEmptySlice()

	// Should pass
	if err := validator([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator([]int{1, 2, 3}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator([]bool{true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on empty
	if err := validator([]string{}); err == nil {
		t.Fatal("expected error for empty string slice")
	}

	if err := validator([]int{}); err == nil {
		t.Fatal("expected error for empty int slice")
	}

	// Should fail on wrong type
	if err := validator("not a slice"); err == nil {
		t.Fatal("expected error for non-slice type")
	}
}

func TestValidateSliceMinLength(t *testing.T) {
	validator := bconf.ValidateSliceMinLength(2)

	// Should pass
	if err := validator([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator([]int{1, 2, 3}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator([]string{"a"}); err == nil {
		t.Fatal("expected error for slice with 1 element")
	}

	// Should fail on wrong type
	if err := validator("not a slice"); err == nil {
		t.Fatal("expected error for non-slice type")
	}
}

func TestValidateSliceMaxLength(t *testing.T) {
	validator := bconf.ValidateSliceMaxLength(2)

	// Should pass
	if err := validator([]string{"a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := validator([]int{1, 2}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail
	if err := validator([]string{"a", "b", "c"}); err == nil {
		t.Fatal("expected error for slice with 3 elements")
	}

	// Should fail on wrong type
	if err := validator("not a slice"); err == nil {
		t.Fatal("expected error for non-slice type")
	}
}

func TestValidateEachString(t *testing.T) {
	validator := bconf.ValidateEachString(bconf.ValidateNonEmptyString())

	// Should pass
	if err := validator([]string{"a", "b", "c"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on empty element
	if err := validator([]string{"a", "", "c"}); err == nil {
		t.Fatal("expected error for slice with empty element")
	}

	// Should fail on wrong type
	if err := validator([]int{1, 2, 3}); err == nil {
		t.Fatal("expected error for non-string-slice type")
	}
}

func TestValidateEachInt(t *testing.T) {
	validator := bconf.ValidateEachInt(bconf.ValidatePort())

	// Should pass
	if err := validator([]int{80, 443, 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail on invalid port
	if err := validator([]int{80, 0, 8080}); err == nil {
		t.Fatal("expected error for slice with invalid port")
	}

	// Should fail on wrong type
	if err := validator([]string{"80", "443"}); err == nil {
		t.Fatal("expected error for non-int-slice type")
	}
}
