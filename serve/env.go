package serve

import "os"

func getEnvOrDefault(env, def string) string {
	if v := os.Getenv(env); len(v) > 0 {
		return v
	}
	return def
}
