package config

import (
	"fmt"
	"os"
)

type Config struct {
	Driver     string // "odbc" (mặc định cho Windows Shared Memory) hoặc "sqlserver"
	DBServer   string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	return &Config{
		Driver:     getEnv("DB_DRIVER", "odbc"),
		DBServer:   getEnv("DB_SERVER", ".\\sqlexpress"),
		DBUser:     getEnv("DB_USER", "sa"),
		DBPassword: getEnv("DB_PASSWORD", "123@456"),
		DBName:     getEnv("DB_NAME", "SOFTZ_SHOPHOA"),
	}
}

// ConnectionString tạo chuỗi kết nối phù hợp với driver
func (c *Config) ConnectionString() string {
	if c.Driver == "odbc" {
		// Kết nối qua Windows Shared Memory (LPC) siêu tốc, không phụ thuộc cổng mạng TCP
		return fmt.Sprintf("Driver={ODBC Driver 17 for SQL Server};Server=%s;Database=%s;Uid=%s;Pwd=%s;TrustServerCertificate=yes;",
			c.DBServer, c.DBName, c.DBUser, c.DBPassword)
	}
	// Fallback pure-Go TDS
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;trustservercertificate=true",
		c.DBServer, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
