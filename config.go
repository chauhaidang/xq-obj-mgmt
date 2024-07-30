package main

import "os"

type Config struct {
	Addr       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	JWTSecret  string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Addr:       getEnv("XADDR", ":8080"),
		DBUser:     getEnv("XDBUSER", "admin"),
		DBPassword: getEnv("XDBPASS", "admin"),
		DBHost:     getEnv("XDBHOST", "localhost"),
		DBName:     getEnv("XDBNAME", "postgres"),
		JWTSecret:  getEnv("JWT_SECRET", "randomjwtforxbank"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
