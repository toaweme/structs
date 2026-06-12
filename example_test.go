package structs_test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/toaweme/structs"
)

// ServerConfig is the kind of struct you already have lying around: one set of
// fields, annotated once, that doubles as your config-file schema and your
// environment-variable contract. structs reads those tags and populates the
// struct from whatever input map you hand it, regardless of which source each
// value came from.
type ServerConfig struct {
	// Host is read from the "host" key in JSON/YAML and falls back to 0.0.0.0
	// when nobody provides it.
	Host string `json:"host" yaml:"host" default:"0.0.0.0"`
	// Port is required, defaults to 8080, and can additionally be sourced from
	// the PORT environment variable. Note it arrives as a string ("9090") and
	// lands as an int: structs coerces it for you.
	Port int `json:"port" yaml:"port" env:"PORT" default:"8080" rules:"required"`
	// LogLevel is constrained to a fixed set via the oneof rule, so a typo in a
	// config file is caught before it ever reaches your program's logic.
	LogLevel string `json:"log_level" yaml:"log_level" default:"info" rules:"oneof:debug,info,warn,error"`
	// Tags is a slice fed from a single comma-separated string ("edge,beta"),
	// the shape you get from an env var or a flat config value. The sep tag
	// picks the delimiter and defaults to a comma.
	Tags []string `json:"tags" yaml:"tags" sep:","`
	// Database is a nested struct. Its fields are addressable by their dotted
	// path ("database.dsn"), by a nested sub-map, or by their own env tag.
	Database Database `json:"database" yaml:"database"`
}

type Database struct {
	DSN string `json:"dsn" yaml:"dsn" env:"DATABASE_DSN" rules:"required"`
}

// Example wires the whole thing together the way a real program would: build a
// manager once, validate the merged inputs, then populate the struct. The
// inputs map is deliberately heterogeneous: most keys came from a decoded
// config file (including a nested "database" sub-map), while "PORT" arrived
// under its env tag. structs resolves each against the right field.
func Example() {
	cfg := &ServerConfig{}

	manager := structs.New(cfg)

	// in a real app this map is the merge of a decoded config file and os.Environ.
	inputs := map[string]any{
		"host":      "127.0.0.1",        // the "host" json/yaml tag
		"PORT":      "9090",             // the env tag, coerced string -> int
		"log_level": "debug",            // the "log_level" json tag
		"tags":      "edge,beta,canary", // split on sep into a []string
		"database": map[string]any{ // a nested config sub-section
			"dsn": "postgres://localhost/app",
		},
	}

	errs, err := manager.Validate(inputs)
	if err != nil {
		panic(err)
	}
	if len(errs) > 0 {
		fmt.Println("config is invalid:", errs)
		return
	}

	if err := manager.Set(inputs); err != nil {
		panic(err)
	}

	fmt.Printf("listen   %s:%d\n", cfg.Host, cfg.Port)
	fmt.Printf("loglevel %s\n", cfg.LogLevel)
	fmt.Printf("tags     %v\n", cfg.Tags)
	fmt.Printf("database %s\n", cfg.Database.DSN)
	// Output:
	// listen   127.0.0.1:9090
	// loglevel debug
	// tags     [edge beta canary]
	// database postgres://localhost/app
}

// Example_environmentVariables shows the same struct populated entirely from
// environment-style keys. Fields with an env tag pick their value up by that
// name, the nested DSN included, and defaults fill in whatever the environment
// did not set.
func Example_environmentVariables() {
	cfg := &ServerConfig{}
	manager := structs.New(cfg)

	// these are the names you would read out of os.Environ.
	env := map[string]any{
		"PORT":         "3000",
		"DATABASE_DSN": "postgres://db:5432/prod",
	}

	if err := manager.Set(env); err != nil {
		panic(err)
	}

	fmt.Printf("host=%s port=%d level=%s dsn=%s\n",
		cfg.Host, cfg.Port, cfg.LogLevel, cfg.Database.DSN)
	// Output:
	// host=0.0.0.0 port=3000 level=info dsn=postgres://db:5432/prod
}

// Example_defaults shows what happens when the caller provides almost nothing:
// default tags fill the gaps, so a near-empty input map still yields a fully
// populated, runnable config. The nested DSN here is addressed by its dotted
// path, the flat-key alternative to a nested sub-map.
func Example_defaults() {
	cfg := &ServerConfig{}
	manager := structs.New(cfg)

	inputs := map[string]any{
		"database.dsn": "postgres://localhost/app",
	}

	if errs, err := manager.Validate(inputs); err != nil {
		panic(err)
	} else if len(errs) > 0 {
		fmt.Println("config is invalid:", errs)
		return
	}
	if err := manager.Set(inputs); err != nil {
		panic(err)
	}

	fmt.Printf("%s:%d level=%s dsn=%s\n", cfg.Host, cfg.Port, cfg.LogLevel, cfg.Database.DSN)
	// Output:
	// 0.0.0.0:8080 level=info dsn=postgres://localhost/app
}

// Example_validation shows the validation side on its own. Validate never
// mutates the struct; it reports which fields failed which rules, keyed by the
// field's resolved tag name. Here Name is missing (required) and Format is not
// in the oneof set, so both rules report.
func Example_validation() {
	type Args struct {
		Name   string `json:"name" rules:"required"`
		Format string `json:"format" rules:"oneof:json,yaml,toml"`
	}

	manager := structs.New(&Args{}, structs.WithTags("json"))

	errs, err := manager.Validate(map[string]any{
		"format": "xml", // not one of json,yaml,toml
		// name omitted, so the required rule fires
	})
	if err != nil {
		panic(err)
	}

	// errs is a map[string][]string; sort the keys for deterministic output.
	fields := make([]string, 0, len(errs))
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		fmt.Printf("%s: %v\n", field, errs[field])
	}
	// Output:
	// format: [must be one of: json, yaml, toml]
	// name: [required]
}

// Example_cliArgs shows that structs is agnostic about where the input map
// comes from. It does not parse argv itself; you bring your own flag parser
// (here the toy cliArgsToMap), and structs matches the resulting keys against
// the arg/short tags. A bool flag with no value becomes true, and a string
// argument is coerced into the field's type.
func Example_cliArgs() {
	type Flags struct {
		Output  string `arg:"output" short:"o" default:"table" rules:"oneof:table,json,yaml"`
		Verbose bool   `arg:"verbose" short:"v"`
		Limit   int    `arg:"limit" default:"50"`
	}

	values := cliArgsToMap([]string{"--output", "json", "-v", "--limit", "10"})

	flags := &Flags{}
	manager := structs.New(flags, structs.WithTags("arg", "short"))

	if errs, err := manager.Validate(values); err != nil {
		panic(err)
	} else if len(errs) > 0 {
		fmt.Println("invalid flags:", errs)
		return
	}
	if err := manager.Set(values); err != nil {
		panic(err)
	}

	fmt.Printf("output=%s verbose=%t limit=%d\n", flags.Output, flags.Verbose, flags.Limit)
	// Output:
	// output=json verbose=true limit=10
}

// cliArgsToMap is a stand-in for whatever flag parser you already use: it turns
// raw argv into a map keyed by flag name (the leading dashes stripped). A flag
// followed by a non-flag token takes that token as its value; otherwise it is
// treated as a present boolean. structs never sees argv, only this map.
func cliArgsToMap(args []string) map[string]any {
	out := make(map[string]any)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			continue
		}
		name := strings.TrimLeft(arg, "-")
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			out[name] = args[i+1]
			i++
			continue
		}
		out[name] = true
	}
	return out
}
