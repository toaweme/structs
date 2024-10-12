package main

import (
	"fmt"

	"github.com/contentforward/structs"
)

type Example struct {
	Name string `json:"name" rules:"required"`
	Age  int    `json:"age" rules:"required"`
}

func main() {
	example := &Example{
		Name: "John Doe",
		Age:  30,
	}

	err := validateStruct(example, map[string]any{
		"name": "Jane Doe",
		"age":  25, // or "25", works either way
	})
	if err != nil {
		panic(fmt.Errorf("error validating struct: %w", err))
	}

	fmt.Printf("Name: %s\n", example.Name)
	fmt.Printf("Age: %d\n", example.Age)
}

func validateStruct(structure any, inputs map[string]any) error {
	manager := structs.New(structure, structs.DefaultRules, structs.DefaultTags...)
	errors, err := manager.Validate(inputs)
	if err != nil {
		return fmt.Errorf("error validating cli command structure: %w", err)
	}

	if len(errors) > 0 {
		for field, rules := range errors {
			for _, rule := range rules {
				fmt.Printf("validation error: %s(%s)", field, rule)
			}
		}

		return fmt.Errorf("validation failed: %v", errors)
	}

	err = manager.Set(inputs)
	if err != nil {
		return fmt.Errorf("failed to set fields: %w", err)
	}

	return nil
}
