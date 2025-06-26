package auth

import "os"

var jwtSecret = []byte(getSecret())

func getSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "replace-with-secure-secret" // fallback for tests/dev
}

func Secret() []byte {
	return jwtSecret
}
