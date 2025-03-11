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

Accessing configuration values can be done by calling lookup methods on a `bconf.AppConfig` with field-set and field
keys, but it is often easier to define a configuration value structure alongside a `bconf.FieldSet`. A
`bconf.AppConfig` can fill these configuration value structs at load time, providing easy access to precisely the
values you need, where you need them.

Check out the documentation and introductory examples below, and see if `bconf` is right for your project!

### Supported Configuration Sources

* Environment (`bconf.EnvironmentLoader`)
* Flags (`bconf.FlagLoader`)
* JSON files (`bconf.JSONFileLoader`)
* Overrides (setter functions)

### Getting Values from `bconf.AppConfig`

* `AttachConfigStructs(configStructs ...any)` (use prior to `Load(...)`)
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

### Features

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
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/segmentio/ksuid"
    "github.com/xavi-group/bconf"
    "github.com/xavi-group/bobotel"
    "github.com/xavi-group/bobzap"
    "go.uber.org/zap"
)

const (
    APIFieldSetKey   = "api"
    SessionSecretKey = "session_secret"
    ReadTimeoutKey   = "read_timeout"
    WriteTimeoutKey  = "write_timeout"
)

func APIFieldSets() bconf.FieldSets {
    checkValidSessionSecret := func(fieldValue any) error {
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
    }

    // FSB() is a shorthand function for NewFieldSetBuilder()
    // FB() is a shorthand function for NewFieldBuilder()
    // C is a shorthand method for Create()
    return bconf.FieldSets{
        bconf.FSB(APIFieldSetKey).
            Fields(
                bconf.FB(SessionSecretKey, bconf.String).Sensitive().Required().
                    Description("API secret for session management").
                    Validator(checkValidSessionSecret).C(),
                bconf.FB(ReadTimeoutKey, bconf.Duration).Default(5*time.Second).
                    Description("API read timeout").C(),
                bconf.FB(WriteTimeoutKey, bconf.Duration).Default(5*time.Second).
                    Description("API write timeout").C(),
            ).C(),
    }
}

type APIConfig struct {
    bconf.ConfigStruct `bconf:"api"`
    LogConfig          *bobzap.Config
    AppID              string        `bconf:"app.id"`
    SessionSecret      string        `bconf:"session_secret"`
    ReadTimeout        time.Duration `bconf:"read_timeout"`
    WriteTimeout       time.Duration `bconf:"write_timeout"`
}

func main() {
    config := bconf.NewAppConfig(
        "auth_service_example",
        "HTTP API for user authentication and authorization",
        bconf.WithAppIDFunc(func() string { return ksuid.New().String() }),
        bconf.WithAppVersion("1.0.0"),
        bconf.WithEnvironmentLoader("auth_svc"),
        bconf.WithFlagLoader(),
    )

    config.AddFieldSetGroup("bobzap", bobzap.FieldSets())
    config.AddFieldSetGroup("bobotel", bobotel.FieldSets())
    config.AddFieldSetGroup("api", APIFieldSets())

    apiConfig := &APIConfig{}

    config.AttachConfigStructs(
        bobzap.NewConfig(),
        bobotel.NewConfig(),
        apiConfig,
    )

    // Load when called without any options will also handle the help flag (--help or -h)
    if errs := config.Load(); len(errs) > 0 {
        fmt.Printf("problem(s) loading application configuration: %v\n", errs)
        os.Exit(1)
    }

    // -- Initialize application observability --
    if err := bobotel.InitializeTraceProvider(); err != nil {
        fmt.Printf("problem initializing application tracer: %s\n", err)
        os.Exit(1)
    }

    if err := bobzap.InitializeGlobalLogger(); err != nil {
        fmt.Printf("problem initializing application logger: %s\n", err)
        os.Exit(1)
    }

    log := bobzap.NewLogger("main")

    log.Info(
        fmt.Sprintf("%s initialized successfully", config.AppName()),
        zap.Any("app_config", config.ConfigMap()),
        zap.Strings("warnings", config.Warnings()),
    )

    // -- Configuration access examples --
}
```

In the code block above, a `bconf.AppConfig` is defined with three field-set groups (which group configuration related
to the logging, tracing, and api in this case). Two field-set groups are from separate packages (bobzap and bobotel).
One is defined right above the main function (api). The api field-set group is defined above to showcase how a
field-set and configuration value struct are implemented, but in practice it would likely be defined within an HTTP
routing package.

The APIConfig value struct is written to showcase the flexibility that exists when filling a configuration value struct.
In this case, it nests the values for the logging configuration as defined by the `bobzap` package, opts to acquire
the `app.id` configuration value from the `app` field-set, and sets a default field-set of `api`, which is used to
locate the `session_secret`, `read_timeout`, and `write_timeout` fields. Although this example is contrived, it is
useful when defining more complex service-like packages.

The easiest way to see your application configuration is to execute your applicaiton with the `-h` or `--help` flag.
The help output will take into account which loaders you have configured, and highlights which configuration values
are required, conditionally required, or optional. Additional help options are in progress.

```
Usage of 'auth_service_example':
HTTP API for user authentication and authorization

Required Configuration:
        api.session_secret string
                API secret for session management.
                Environment key: 'AUTH_SVC_API_SESSION_SECRET'
                Flag argument: '--api_session_secret'
Conditionally Required Configuration:
        otlp.host string
                Otlp host defines the host location of the trace collector.
                Environment key: 'AUTH_SVC_OTLP_HOST'
                Flag argument: '--otlp_host'
                Loading depends on field(s): 'otel.exporters'
Optional Configuration:
        api.read_timeout time.Duration
                Timeout for HTTP read operations.
                Default value: '5s'
                Environment key: 'AUTH_SVC_API_READ_TIMEOUT'
                Flag argument: '--api_read_timeout'
        api.write_timeout time.Duration
                Timeout for HTTP write operations.
                Default value: '5s'
                Environment key: 'AUTH_SVC_API_WRITE_TIMEOUT'
                Flag argument: '--api_write_timeout'
        app.id string
                Default value: <generated-at-run-time>
                Environment key: 'AUTH_SVC_APP_ID'
                Flag argument: '--app_id'
        app.version string
                Default value: '1.0.0'
                Environment key: 'AUTH_SVC_APP_VERSION'
                Flag argument: '--app_version'
        log.color bool
                Log color defines whether console formatted logs are rendered in color.
                Default value: 'true'
                Environment key: 'AUTH_SVC_LOG_COLOR'
                Flag argument: '--log_color'
        log.config string
                Log config defines whether the Zap will be configured with development or production defaults. 
                Note: `development` defaults to debug log level and console format, `production` defaults to info log 
                level and json format. 
                Accepted values: ['production', 'development']
                Default value: 'production'
                Environment key: 'AUTH_SVC_LOG_CONFIG'
                Flag argument: '--log_config'
        log.format string
                Log format defines the format logs will be emitted in (overrides log config defaults).
                Accepted values: ['console', 'json']
                Environment key: 'AUTH_SVC_LOG_FORMAT'
                Flag argument: '--log_format'
        log.level string
                Log level defines the level at which logs will be emitted (overrides log config defaults).
                Accepted values: ['debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal']
                Environment key: 'AUTH_SVC_LOG_LEVEL'
                Flag argument: '--log_level'
        otel.console_format string
                Otel console format defines the format of traces output to the console where 'pretty' is more 
                human readable (adds whitespace). 
                Accepted values: ['production', 'pretty']
                Default value: 'production'
                Environment key: 'AUTH_SVC_OTEL_CONSOLE_FORMAT'
                Flag argument: '--otel_console_format'
        otel.exporters []string
                Otel exporters defines where traces will be sent (accepted values are 'console' and 'otlp'). 
                Exporters accepts a list and can be configured to export traces to multiple destinations. 
                Default value: '[console]'
                Environment key: 'AUTH_SVC_OTEL_EXPORTERS'
                Flag argument: '--otel_exporters'
        otlp.endpoint_kind string
                Otlp endpoint kind defines the protocol used by the trace collector.
                Accepted values: ['http', 'grpc']
                Default value: 'http'
                Environment key: 'AUTH_SVC_OTLP_ENDPOINT_KIND'
                Flag argument: '--otlp_endpoint_kind'
                Loading depends on field(s): 'otel.exporters'
        otlp.port int
                Otlp port defines the port of the trace collector process. For a GRPC endpoint the default is 
                4317. 
                Default value: '4318'
                Environment key: 'AUTH_SVC_OTLP_PORT'
                Flag argument: '--otlp_port'
                Loading depends on field(s): 'otel.exporters'
```

To view more examples, including a real-world example showcasing how configuration can live alongside package code,
please visit [github.com/xavi-group/bapp-template](https://github.com/xavi-group/bapp-template).

If you are interested in quickly configuring your application with tracing and logging as showcased in the example
above, consider checking out [github.com/xavi-group/bobzap](https://github.com/xavi-group/bobzap)
and [github.com/xavi-group/bobotel](https://github.com/xavi-group/bobotel).

## Roadmap / Future Improvements

* Additional field type support (maps)
* File watching and notifications for configuration value updates
* YAML files (`bconf.YAMLFileLoader`)
* TOML files (`bconf.TOMLFileLoader`)
* Additional `-h` / `--help` options
* Provide common field validator functions
* Implement `Validators` and `Transformers` on `bconf.Field`
