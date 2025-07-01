package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type YamlConfig struct {
	Database DBConfig `yaml:"database"`
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func OpenDBConnection() (*sql.DB, error) {
	// Đọc file cấu hình
	yamlFile, err := ioutil.ReadFile("security/dbconfig.yml")
	if err != nil {
		return nil, err
	}

	var cfg YamlConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	// Kiểm tra kết nối
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}