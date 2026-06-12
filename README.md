# structs

[![Quality](https://github.com/toaweme/structs/actions/workflows/tests.yml/badge.svg)](https://github.com/toaweme/structs/actions/workflows/tests.yml)
[![Go Reference](https://img.shields.io/badge/Docs-pkg.go.dev-blue)](https://pkg.go.dev/github.com/toaweme/structs)
[![GitHub Tag](https://img.shields.io/github/v/tag/toaweme/structs?label=Tag&color=green)](https://github.com/toaweme/structs/releases)
[![License](https://img.shields.io/badge/License-MIT-blue)](/LICENSE)

## Fill and read Go's structs

`github.com/toaweme/structs` gives you runtime tools to work with Go's structs, its fields, tags and values.

This module was originally built as a fun way to solve the CLI app arg parsing problem.
I'm a big fan of simplicity and the stdlib while powerful, doesn't make CLI flag/arg parsing simple, there's a lot of boilerplate.
`structs` abstracts the complicated bits and can magically set struct field values (however nested) from a simple `map[string]any`.

## Install

```sh
go get github.com/toaweme/structs
```

## Overview

### Struct embedding and nesting

#### Nesting (a named struct field)

Define your structs:

```go
type Server struct {
	Database Database `json:"database" env:"DATABASE"`
}

type Database struct {
	URL string `json:"url" env:"URL"`
}
```

A nested field can be reached three ways, all equivalent:

```go
// 1. dotted path: the field's tag glued to its parent's with "."
map[string]any{
	"database.url": "mysql://127.0.0.1:3306/beep",
}

// 2. nested map: a sub-section keyed by the parent's tag
map[string]any{
	"database": map[string]any{
		"url": "mysql://127.0.0.1:3306/beep",
	},
}

// 3. env tag: the env tags glued with "_"
map[string]any{
	"DATABASE_URL": "mysql://127.0.0.1:3306/beep",
}
```

Nesting goes arbitrarily deep (`a.b.c`, or maps within maps). This is how a
decoded JSON/YAML config drops straight in.

#### Embedding (an anonymous struct field)

An untagged embedded struct has its fields promoted to the parent level, exactly
as Go (and `encoding/json`) promote them: no wrapper, no prefix. The embedded
type may be exported or unexported.

```go
type Network struct {
	Host string `json:"host" env:"HOST"`
	Port int    `json:"port" env:"PORT"`
}

type Server struct {
	Network        // embedded: Host and Port are promoted
	Name string `json:"name"`
}
```

Set the promoted fields by their own tag or name, with no parent prefix:

```go
map[string]any{
	"host": "127.0.0.1", // -> Server.Host
	"port": 8080,        // -> Server.Port
	"name": "edge",      // -> Server.Name
}
```

A *tagged* anonymous field is not promoted; it nests under its tag instead, just
like `encoding/json`, so it behaves like the named nesting above.

#### Limitations

- Nested maps must be `map[string]any` at every level (the form JSON/YAML
  decoders produce). A value whose concrete type is a typed map such as
  `map[string]map[string]any` is only descended into where its element type is
  `map[string]any`; a deeper typed-map intermediate is not traversed, so the
  leaf stays unset. Use the dotted path or a `map[string]any` sub-section.


## Module

- `structs.New` a small abstraction to Validate and Set.
    - `structs.WithTags` a priority list of tags for `Set` (default: `["json", "yaml"]`).
    - `structs.WithEncodingTags` a list of tags in which commas are treated as encoding configuration (e.g. `json:"field,omitempty"`).
    - `structs.WithRules` extend or replace the built-in validation rules.
    - `structs.WithValidationTag` tag used to define the validation rules (default: `rules`)
- `structs.GetStructFields` reads the entire nested struct field tree.
- `structs.SetStructFields` takes a `map[string]any` and fills the struct fields.
- `structs.ValidateStructFields` uses a rule map to validate your `map[string]any` against selected fields.

## Features

- **Validate without mutating** - check inputs against each field's rules and get
  back a map of field names with the validation messages
- **Populate from a single map** - fill a struct from one map of values,
  matching each field and converting the value into the field's type.
- **Type coercion** - string, int, float, bool, slice, map, and interface fields
  are all set from loosely typed inputs, so a port given as the string "9090"
  lands in an int field.
- **Tag priority** - decide which struct tag names a field by giving an ordered
  list; the first tag a field carries wins. Defaults to json then yaml, and is overridable.
- **Defaults** - a field left empty is seeded from its declared default value,
  and a default never overrides a value that is already present.
- **Built-in validation rules** - required and one-of out of the box, with the
  ability to add your own named rules or replace the built-in set.
- **Slice splitting** - a single string handed to a scalar slice field is split
  into elements (comma by default, or a custom separator per field) and each
  element is converted; already-structured inputs pass through untouched.

- **Nested structs** - reach a field inside a nested struct by dotted path, by a
  nested map, or by an env-style key, to any depth.
- **Embedded structs** - fields of an anonymous embedded struct are promoted and
  set directly, the way Go does it, whether the embedded type is exported or not.

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
