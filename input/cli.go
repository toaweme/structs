package input

import (
	"os"
	"strings"
)

func ArgsToMap() map[string]any {
	args := os.Args[1:]
	argsMap := make(map[string]any, len(args))
	for i, arg := range args {
		if strings.Contains(arg, "=") {
			pair := strings.Split(arg, "=")
			argsMap[pair[0]] = pair[1]
		} else if i+1 < len(args) {
			argsMap[arg] = args[i+1]
		}
	}
	return argsMap
}
