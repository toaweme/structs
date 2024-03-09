package input

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/contentforward/structs/utils"
)

func NewBag() Bag {
	bag := Bag{
		Env: EnvToMap(),
		// Structure:    ArgsToMap(),
		Headers: map[string]string{},
		Request: map[string]string{},
	}

	// spew.Dump(bag)

	return bag
}

type Bag struct {
	Data    map[string]string
	Env     map[string]string
	Args    map[string]string
	Headers map[string]string
	Request map[string]string
}

func (bag *Bag) Collect() map[string]any {
	vars := make(map[string]any)
	for k, v := range bag.Env {
		_ = k
		_ = v
		vars[k] = v
		vars["env."+strings.ToLower(k)] = v
	}
	for k, v := range bag.Data {
		vars["data."+k] = v
	}
	for k, v := range bag.Args {
		vars[k] = v
	}
	for k, v := range bag.Headers {
		vars[k] = v
	}
	for k, v := range bag.Request {
		vars[k] = v
	}
	return vars
}

func (bag *Bag) GetBool(name string) (bool, error) {
	value, ok := bag.lookupValue(name)
	if !ok {
		return false, fmt.Errorf("variable %s not found", name)
	}

	val := utils.ParseBool(value)

	return val, nil
}

func (bag *Bag) GetInt(name string) (int, error) {
	value, ok := bag.lookupValue(name)
	if !ok {
		return 0, fmt.Errorf("variable %s not found", name)
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("variable %s is not an integer", name)
	}

	return val, nil
}

func (bag *Bag) GetString(name string) (string, error) {
	value, ok := bag.lookupValue(name)
	if !ok {
		return "", fmt.Errorf("variable %s not found", name)
	}

	return value, nil
}

func (bag *Bag) lookupValue(name string) (string, bool) {
	lookup, ok := bag.lookup(name)
	if !ok {
		normName := strings.ToLower(name)
		lookup, ok = bag.lookup(normName)
		if !ok {
			return "", true
		}
	}
	return lookup, false
}

func (bag *Bag) lookup(name string) (string, bool) {
	if val, ok := bag.Args[name]; ok {
		return val, true
	}

	if val, ok := bag.Data[name]; ok {
		return val, true
	}

	if val, ok := bag.Env[name]; ok {
		return val, true
	}

	if val, ok := bag.Headers[name]; ok {
		return val, true
	}

	if val, ok := bag.Request[name]; ok {
		return val, true
	}

	return "", false
}
