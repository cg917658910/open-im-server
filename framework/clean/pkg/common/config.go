package common

import "os"

type Config struct {
	HTTPAddr   string
	SocketAddr string
	RPCAddr    string
}

func LoadConfig() Config {
	cfg := Config{
		HTTPAddr:   getEnv("FRAMEWORK_HTTP_ADDR", ":18080"),
		SocketAddr: getEnv("FRAMEWORK_SOCKET_ADDR", ":18081"),
		RPCAddr:    getEnv("FRAMEWORK_RPC_ADDR", ":18082"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
