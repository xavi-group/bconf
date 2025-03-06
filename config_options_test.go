package bconf_test

import (
	"testing"

	"github.com/xavi-group/bconf"
)

func TestConfigOptions(t *testing.T) {
	environmentLoaderOption := bconf.WithEnvironmentLoader("")

	t.Logf("testing: %s", environmentLoaderOption.ConfigOptionType())
}
