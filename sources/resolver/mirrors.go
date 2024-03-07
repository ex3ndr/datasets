package resolver

import (
	"os"
	"strings"
)

func ResolveMirror(link string) string {

	// Load mirrors from environment
	envMirrors := os.Getenv("DATASETS_MIRRORS")
	if envMirrors != "" {
		mirrorPairs := strings.Split(envMirrors, ",")
		for _, pair := range mirrorPairs {
			mirror := strings.Split(pair, "=")
			if len(mirror) == 2 {
				source := strings.ToLower(mirror[0])
				target := mirror[1]
				if strings.HasPrefix(strings.ToLower(link), source) {
					return target + link[len(source):]
				}
			}
		}
	}

	return link
}
