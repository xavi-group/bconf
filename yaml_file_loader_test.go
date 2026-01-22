package bconf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xavi-group/bconf"
	"gopkg.in/yaml.v3"
)

func TestYAMLFileLoaderDefaults(t *testing.T) {
	loader := bconf.NewYAMLFileLoader()

	if loader == nil {
		t.Fatalf("unexpected nil loader")
	}

	if loader.Decoder == nil {
		t.Fatalf("expected default decoder to be set")
	}
}

func TestYAMLFileLoaderWithOptions(t *testing.T) {
	loader := bconf.NewYAMLFileLoader(
		bconf.WithYAMLDecoder(yaml.Unmarshal),
		bconf.WithYAMLFilePaths("./fixtures/yaml_config_test_fixture_01.yaml"),
	)

	if len(loader.FilePaths) != 1 {
		t.Fatalf("unexpected file-paths length '%d', expected '1'", len(loader.FilePaths))
	}
}

func TestYAMLFileLoaderClone(t *testing.T) {
	loader := yamlLoaderWithTestFixture01()
	clone := loader.Clone()

	if len(clone.FilePaths) != len(loader.FilePaths) {
		t.Fatalf("unexpected clone file-path length '%d', expected '%d'", len(clone.FilePaths), len(loader.FilePaths))
	}

	loader.FilePaths[0] = "./fixtures/yaml_config_test_fixture_02.yaml"

	if clone.FilePaths[0] == loader.FilePaths[0] {
		t.Fatalf("unexpected clone file-path value: %s", clone.FilePaths[0])
	}

	loader.FilePaths[0] = "./fixtures/yaml_config_test_fixture_01.yaml"

	loaderClone := loader.CloneLoader()

	loader.FilePaths[0] = "./fixtures/empty.yaml"

	_, found := loaderClone.Get("app", "id")
	if !found {
		t.Fatalf("unexpected issue finding app-id")
	}
}

func TestYAMLFileLoaderName(t *testing.T) {
	loader := bconf.NewYAMLFileLoader()

	if loader.Name() != "bconf_yamlfile" {
		t.Fatalf("unexpected yaml-file-loader name '%s'", loader.Name())
	}
}

//nolint:dupl // Test code intentionally follows similar patterns
func TestYAMLFileLoaderGet(t *testing.T) {
	loaderFixture01 := yamlLoaderWithTestFixture01()
	loaderNoFilePaths := yamlLoaderWithNoFilePaths()
	loaderInvalidFilePaths := yamlLoaderWithInvalidFilePaths()
	loaderBadDecoder := yamlLoaderWithBadDecoder()

	_, found := loaderFixture01.Get("strange_key", "some_field")
	if found {
		t.Fatalf("unexpected found value when looking for non-existent key")
	}

	appID, found := loaderFixture01.Get("app", "id")

	if !found {
		t.Fatalf("expected loader with fixture file to find appID value")
	}

	if appID != "test-app-id" {
		t.Fatalf("unexpected appID value '%s', expected 'test-app-id'", appID)
	}

	internalPorts, found := loaderFixture01.Get("app", "internal_ports")
	if !found {
		t.Fatalf("expected loader with fixture file to find internalPorts value")
	}

	if internalPorts != "8081,8082" {
		t.Fatalf("unexpected internalPorts value '%s', expected '8081,8082'", internalPorts)
	}

	_, found = loaderNoFilePaths.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with no file-paths")
	}

	_, found = loaderInvalidFilePaths.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with invalid file-paths")
	}

	_, found = loaderBadDecoder.Get("app", "id")
	if found {
		t.Fatalf("unexpected appID found by loader with bad decoder")
	}
}

func TestYAMLFileLoaderGetMap(t *testing.T) {
	loaderFixture01 := yamlLoaderWithTestFixture01()
	loaderNoFilePaths := yamlLoaderWithNoFilePaths()

	appMap := loaderFixture01.GetMap("app", []string{"id", "secret", "invalid_field_key"})
	if len(appMap) != 2 {
		t.Fatalf("unexpected length of app field-set map '%d', expected '2'", len(appMap))
	}

	appMap = loaderNoFilePaths.GetMap("app", []string{"id", "secret", "invalid_field_key"})
	if len(appMap) != 0 {
		t.Fatalf("unexpected length of app file-set map '%d', expected '0'", len(appMap))
	}
}

func TestYAMLFileLoaderHelpString(t *testing.T) {
	loaderFixture01 := yamlLoaderWithTestFixture01()

	helpString := loaderFixture01.HelpString("app", "id")

	if !strings.Contains(helpString, "app.id") {
		t.Fatalf("unexpected help string: '%s'", helpString)
	}

	if !strings.Contains(helpString, "YAML") {
		t.Fatalf("expected help string to contain 'YAML': '%s'", helpString)
	}
}

func TestYAMLFileLoaderArrayWithCommas(t *testing.T) {
	loaderFixture02 := yamlLoaderWithTestFixture02()

	tags, found := loaderFixture02.Get("app", "tags")
	if !found {
		t.Fatalf("expected to find tags value")
	}

	expected := "hello, world,foo,bar, baz"
	if tags != expected {
		t.Fatalf("unexpected tags value '%s', expected '%s'", tags, expected)
	}
}

func TestYAMLFileLoaderCaching(t *testing.T) {
	loader := yamlLoaderWithTestFixture01()

	// First call should load files
	appID1, found := loader.Get("app", "id")
	if !found {
		t.Fatalf("expected to find app id")
	}

	// Second call should use cached data
	appID2, found := loader.Get("app", "id")
	if !found {
		t.Fatalf("expected to find app id on second call")
	}

	if appID1 != appID2 {
		t.Fatalf("cached value mismatch: '%s' vs '%s'", appID1, appID2)
	}

	// Verify GetMap also uses cached data
	appMap := loader.GetMap("app", []string{"id", "secret"})
	if len(appMap) != 2 {
		t.Fatalf("unexpected map length: %d", len(appMap))
	}

	if appMap["id"] != appID1 {
		t.Fatalf("GetMap returned different value than Get: '%s' vs '%s'", appMap["id"], appID1)
	}
}

func yamlLoaderWithTestFixture01() *bconf.YAMLFileLoader {
	return bconf.NewYAMLFileLoader(
		bconf.WithYAMLFilePaths("./fixtures/yaml_config_test_fixture_01.yaml"),
	)
}

func yamlLoaderWithTestFixture02() *bconf.YAMLFileLoader {
	return bconf.NewYAMLFileLoader(
		bconf.WithYAMLFilePaths("./fixtures/yaml_config_test_fixture_02.yaml"),
	)
}

func yamlLoaderWithBadDecoder() *bconf.YAMLFileLoader {
	badDecoder := func(_ []byte, _ any) error {
		return fmt.Errorf("decoder error")
	}

	return bconf.NewYAMLFileLoader(
		bconf.WithYAMLDecoder(badDecoder),
		bconf.WithYAMLFilePaths("./fixtures/yaml_config_test_fixture_01.yaml"),
	)
}

func yamlLoaderWithNoFilePaths() *bconf.YAMLFileLoader {
	return bconf.NewYAMLFileLoader()
}

func yamlLoaderWithInvalidFilePaths() *bconf.YAMLFileLoader {
	return bconf.NewYAMLFileLoader(
		bconf.WithYAMLFilePaths("./fixtures/non-existent-file.yaml"),
	)
}
