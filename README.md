# `bconf`: builder configuration for go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/xavi-group/bconf?status.svg)](https://pkg.go.dev/github.com/xavi-group/bconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/xavi-group/bconf)](https://goreportcard.com/report/github.com/xavi-group/bconf)
[![Build Status](https://github.com/xavi-group/bconf/actions/workflows/golang-test.yml/badge.svg?branch=main)](https://github.com/xavi-group/bconf/actions/workflows/golang-test.yml)
[![codecov.io](https://codecov.io/github/xavi-group/bconf/coverage.svg?branch=main)](https://codecov.io/github/xavi-group/bconf?branch=main)

`bconf` is a configuration framework that makes it easy to define, load, and validate application configuration values.

```sh
go get github.com/xavi-group/bconf
```

### Why `bconf`

`bconf` provides tooling to write your configuration package by package. With `bconf`, configuration lives right
alongside the code that needs it. This also makes it so that configuration is more easily re-used and composible by
multiple applications (just like your packages should be).

`bconf` accomplishes this with `bconf.FieldSets`, which provide a namespace and logical grouping for related
configuration. Independent packages define their `bconf.FieldSets`, and then application executables can attach them
to a `bconf.AppConfig`, which provides a unified structure for loading and retrieving configuration values.

Within `bconf.FieldSets`, you define `bconf.Fields`, with each field defining the expected format and behavior of a
configuration value.

Check out the documentation and introductory examples below, and see if `bconf` is right for your project!

### Supported Configuration Sources

* Environment (`bconf.EnvironmentLoader`)
* Flags (`bconf.FlagLoader`)
* JSON files (`bconf.JSONFileLoader`)
* Overrides (setter functions)

### Getting Values from `bconf.AppConfig`

* `FillStruct(configStruct any) error`
* `GetField(fieldSetKey, fieldKey string) (*bconf.Field, error)`
* `GetString(fieldSetKey, fieldKey string) (string, error)`
* `GetStrings(fieldSetKey, fieldKey string) ([]string, error)`
* `GetInt(fieldSetKey, fieldKey string) (int, error)`
* `GetInts(fieldSetKey, fieldKey string) ([]int, error)`
* `GetBool(fieldSetKey, fieldKey string) (bool, error)`
* `GetBools(fieldSetKey, fieldKey string) ([]bool, error)`
* `GetTime(fieldSetKey, fieldKey string) (time.Time, error)`
* `GetTimes(fieldSetKey, fieldKey string) ([]time.Time, error)`
* `GetDuration(fieldSetKey, fieldKey string) (time.Duration, error)`
* `GetDurations(fieldSetKey, fieldKey string) ([]time.Duration, error)`

### Additional Features

* Ability to generate default configuration values with the `bconf.Field` `DefaultGenerator` parameter
* Ability to define custom configuration value validation with the `bconf.Field` `Validator` parameter
* Ability to conditionally load a `bconf.FieldSet` by defining `bconf.LoadConditions`
* Ability to conditionally load a `bconf.Field` by defining `bconf.LoadConditions`
* Ability to get a safe map of configuration values from the `bconf.AppConfig` `ConfigMap()` function
  * (the configuration map will obfuscate values from fields with `Sensitive` parameter set to `true`)
* Ability to reload field-sets and individual fields via the `bconf.AppConfig`
* Ability to fill configuration structures with values from a `bconf.AppConfig` using the `FillStruct(...)` method

### Limitations

* No support for watching / automatically updating configuration values

### Example

Below is an example of a `bconf.AppConfig` defined with the builder pattern. Below this code block the behavior of the
example is discussed.

```go
configuration := bconf.NewAppConfig(
    "external_http_api",
    "HTTP API for user authentication and authorization",
    bconf.WithAppIDFunc(),
    bconf.WithEnvironmentLoader("ext_http_api"),
    bconf.WithFlagLoader(""),
)

configuration.AddFieldSetGroup(
    bconf.FSB("api").Fields( // FSB() is a shorthand function for NewFieldSetBuilder()
        bconf.FB("session_secret", bconf.String). // FB() is a shorthand function for NewFieldBuilder()
            Description("API secret for session management").Sensitive().Required().
            Validator(
                func(fieldValue any) error {
                    secret, _ := fieldValue.(string)

                    minLength := 20
                    if len(secret) < minLength {
                        return fmt.Errorf(
                            "expected string of minimum %d characters (len=%d)",
                            minLength,
                            len(secret),
                        )
                    }

                    return nil
                },
            ).C(), // C is a shorthand method for Create()
    ).C(),
    bconf.FSB("log").Fields(
        bconf.FB("level", bconf.String).
            .Default("info").Description("Logging level").Enumeration("debug", "info", "warn", "error").C(),
        bconf.FB("format", bconf.String).
            .Default("json").Description("Logging format").Enumeration("console", "json").C(),
        bconf.FB("color_enabled", bconf.Bool).
            .Default(true).Description("Colored logs when format is 'console'").C(),
    ).C(),
)

// Load when called without any options will also handle the help flag (--help or -h)
if errs := configuration.Load(); len(errs) > 0 {
    // handle configuration load errors
}

// returns the log level found in order of: user override -> flag -> environment -> default
// (based on the loaders set above).
logLevel, err := configuration.GetString("log", "level")
if err != nil {
    // handle error
}

fmt.Printf("log-level: %s\n", logLevel)

type loggerConfig struct {
    bconf.ConfigStruct `bconf:"log"`
    Level string `bconf:"level"`
    Format string `bconf:"format"`
    ColorEnabled bool `bconf:"color_enabled"`
}

logConfig := &loggerConfig{}
if err := configuration.FillStruct(logConfig); err != nil {
    // handle error
}

fmt.Printf("log config: %v\n", logConfig)
```

In the code blocks above, a `bconf.AppConfig` is defined with two field-sets (which group configuration related to the
application and logging in this case), and registered with help flag parsing.

If this code was executed in a `main()` function, it would print the log level picked up by the configuration from the
flags or run-time environment before falling back on the defined default value of "info". It would then fill the
`logConfig` struct with multiple values from the log field-set fields, and print those values as well.

If this code was executed inside the `main()` function and passed a `--help` or `-h` flag, it would print the following
output:

```
Usage of 'external_http_api':
HTTP API for user authentication and authorization

Required Configuration:
        api_session_secret string
                Application secret for session management
                Environment key: 'EXT_HTTP_API_APP_SESSION_SECRET'
                Flag argument: '--app_session_secret'
Optional Configuration:
        app_id string
                Application identifier
                Default value: <generated-at-run-time>
                Environment key: 'EXT_HTTP_API_APP_ID'
                Flag argument: '--app_id'
        log_color_enabled bool
                Colored logs when format is 'console'
                Default value: 'true'
                Environment key: 'EXT_HTTP_API_LOG_COLOR_ENABLED'
                Flag argument: '--log_color_enabled'
        log_format string
                Logging format
                Accepted values: ['console', 'json']
                Default value: 'json'
                Environment key: 'EXT_HTTP_API_LOG_FORMAT'
                Flag argument: '--log_format'
        log_level string
                Logging level
                Accepted values: ['debug', 'info', 'warn', 'error']
                Default value: 'info'
                Environment key: 'EXT_HTTP_API_LOG_LEVEL'
                Flag argument: '--log_level'
```

This is a simple example where all the configuration code is in one place, but it doesn't need to be!

To view more examples, including a real-world example showcasing how configuration can live alongside package code,
please visit [github.com/xavi-group/bapp-template](https://github.com/xavi-group/bapp-template).

## Roadmap / Future Improvements

* Additional field types (maps)
* Support for file watching and notifications for configuration value updates
* YAML files (`bconf.YAMLFileLoader`)
* TOML files (`bconf.TOMLFileLoader`)
* Additional `-h` / `--help` options
