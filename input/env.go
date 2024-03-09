package input

import (
	"os"
	"strings"
)

func EnvToMap() map[string]string {
	env := os.Environ()
	envMap := make(map[string]string)
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		envMap[pair[0]] = pair[1]
	}
	return envMap
}
