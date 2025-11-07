package main

import (
	"embed"
	"errors"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/open-function-computers-llc/uptime/server"
)

//go:embed dist/*
var dist embed.FS

// version can be changed by build flags
var Version = "latest-dev"

func main() {
	err := verifyVaildEnv()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	server, err := server.Create(dist, os.Getenv("APP_URL"))
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func verifyVaildEnv() error {
	requiredKeys := []string{
		"APP_PORT",
		"APP_URL",
		"APP_ENV",
		"SMTP_HOST",
		"SMTP_USER",
		"SMTP_PASSWORD",
		"SMTP_PORT",
		"DB_TYPE",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_HOST",
		"DB_PORT",
		"INTERVAL_OK_RECHECK",
		"INTERVAL_ERROR_RECHECK",
		"INTERVAL_HOW_LONG_FOR_SITE_TIMEOUT",
		"DANGER_SECONDS",
		"EMERGENCY_SECONDS",
	}

	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			return errors.New("missing env: " + key)
		}
	}

	// supported DB engines
	if os.Getenv("DB_TYPE") != "mysql" && os.Getenv("DB_TYPE") != "mariadb" {
		return errors.New("supported db engines are mysql or mariadb")
	}

	requiredKeysThatAreInts := []string{
		"APP_PORT",
		"SMTP_PORT",
		"DB_PORT",
		"INTERVAL_OK_RECHECK",
		"INTERVAL_ERROR_RECHECK",
		"INTERVAL_HOW_LONG_FOR_SITE_TIMEOUT",
		"DANGER_SECONDS",
		"EMERGENCY_SECONDS",
	}
	for _, key := range requiredKeysThatAreInts {
		_, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			return err
		}
	}

	return nil
}
