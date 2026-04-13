package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	MongoURI     string
	Port         string
	JWTSecret    string
	FrontendRoot string
}

// ResolveFrontendRoot returns the directory that contains "static" and "templates"
// (FRONTEND_ROOT env, or auto-detect frontend/ vs ../frontend).
func ResolveFrontendRoot() string {
	if v := os.Getenv("FRONTEND_ROOT"); v != "" {
		return v
	}
	for _, candidate := range []string{"frontend", filepath.Join("..", "frontend")} {
		st, err := os.Stat(filepath.Join(candidate, "static"))
		if err == nil && st.IsDir() {
			return candidate
		}
	}
	return "frontend"
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	return &Config{
		MongoURI:     os.Getenv("MONGODB_URI"),
		Port:         port,
		JWTSecret:    secret,
		FrontendRoot: ResolveFrontendRoot(),
	}
}
