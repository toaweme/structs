# structs

[![Quality](https://github.com/toaweme/structs/actions/workflows/tests.yml/badge.svg)](https://github.com/toaweme/structs/actions/workflows/tests.yml)
[![Go Reference](https://img.shields.io/badge/Docs-pkg.go.dev-blue)](https://pkg.go.dev/github.com/toaweme/structs)
[![GitHub Tag](https://img.shields.io/github/v/tag/toaweme/structs?label=Tag&color=green)](https://github.com/toaweme/structs/releases)
[![License](https://img.shields.io/badge/License-MIT-blue)](/LICENSE)

## Use Go's structs to innovate

`github.com/toaweme/structs` gives you the tools to work with Go's structs, its fields, tags and values.

This module was originally built as a fun way to solve the CLI app arg parsing problem.
I'm a big fan of simplicity and the stdlib while powerful, doesn't make CLI flag/arg parsing simple, there's a lot of boilerplate.
`structs` abstracts the complicated bits and can magically set struct field values (however nested) from a simple `map[string]any`.

## Module

- `structs.New` a small abstraction to Validate and Set.
  - `structs.WithTags` a priority list of tags for `Set` (default: `["json", "yaml"]`).
  - `structs.WithEncodingTags` a list of tags in which commas are treated as encoding configuration (e.g. `json:"field,omitempty"`).
  - `structs.WithRules` extend or replace the built-in validation rules.
  - `structs.WithValidationTag` tag used to define the validation rules (default: `rules`)
- `structs.GetStructFields` reads the entire nested struct field tree.
- `structs.SetStructFields` takes a `map[string]any` and fills the struct fields.
- `structs.ValidateStructFields` uses a rule map to validate your `map[string]any` against selected fields.

## Quickstart

### Nested structs

Define your structs:

```go
type Server struct {
	Database Database `json:"database" env:"DATABASE"`
}

type Database struct {
	URL string `json:"url" env:"URL"`
}
```

Then use a map to set the fields:

```go
map[string]any{
	"database.url": "mysql://127.0.0.1:3306/beep"
}
```

or 

```go
map[string]any{
	"DATABASE_URL": "mysql://127.0.0.1:3306/beep"
}
```


## Install

```sh
go get github.com/toaweme/structs
```

## Features

- **Validate without mutating** - `Validate(inputs)` runs each field's `rules:`
  and returns a `map[field][]messages`; an empty map means everything passed.
- **Populate from a `map[string]any`** - `Set(inputs)` applies `default:` values
  then matches each field by tag, coercing the value into the field's type.
- **Tag priority** - matches by the first tag a field carries (default`["json", "yaml"]`), overridable with
  `structs.WithTags(...)`.
- **Defaults** - `default:"..."` seeds empty fields.
- **Built-in rules** - `required` and `oneof:a,b,c`, extend or replace them with
  `structs.WithRules(...)`.
- **Slice splitting** - a scalar slice fed one string is split on the field's
  `sep` tag (default `,`) and each element coerced; structured inputs pass through.
- **Nested and embedded structs** - reach nested fields by dotted path, nested
  map, or `env` tag; embedded structs promote their fields like Go does.

> This package does not read the env or any other value source. That's your responsibility.

---

## Quickstart

```go
type ServerConfig struct {
    Host     string   `json:"host" yaml:"host" default:"0.0.0.0"`
    Port     int      `json:"port" yaml:"port" env:"PORT" default:"8080" rules:"required"`
    LogLevel string   `json:"log_level" yaml:"log_level" default:"info" rules:"oneof:debug,info,warn,error"`
    Tags     []string `json:"tags" yaml:"tags" sep:","`
    Database Database `json:"database" yaml:"database"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn" env:"DATABASE_DSN" rules:"required"`
}

cfg := &ServerConfig{}
structManager := structs.New(cfg)

// it's your responsibility to collect the values
// inputs := merge(env(), config())
inputs := map[string]any{
	"host":      "127.0.0.1",        // matched by the json/yaml "host" tag
	"PORT":      "9090",             // matched by the env tag, coerced to int
	"log_level": "debug",            // matched by the "log_level" json tag
	"tags":      "edge,beta,canary", // split on sep into []string
	"database":  map[string]any{     // nested sub-section, matched by dotted path
		"dsn": "postgres://localhost/app",
	},
}

if errs, err := structManager.Validate(inputs); err != nil {
	log.Fatal(err)
} else if len(errs) > 0 {
	log.Fatalf("config is invalid: %v", errs)
}

if err := structManager.Set(inputs); err != nil {
	log.Fatal(err)
}
// cfg.Port == 9090, cfg.Tags == ["edge","beta","canary"], cfg.Database.DSN set.
```

## Runnable examples

See [`example_test.go`](./example_test.go) for the full, runnable versions of everything mentioned above.

```sh
go test -run Example -v
```
