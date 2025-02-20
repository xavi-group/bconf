package bconf_test

import (
	"testing"

	"github.com/rheisen/bconf"
)

func TestConfigOptions(t *testing.T) {
	environmentLoaderOption := bconf.WithEnvironmentLoader("")

	t.Logf("testing: %s", environmentLoaderOption.ConfigOptionType())
}
