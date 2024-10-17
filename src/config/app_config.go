package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatusConfig struct {
	Mongo MongoConfig
	Path  PathConfig
}

type MongoConfig struct {
	AuthEnabled        bool
	Host               string
	Port               int
	ConnectionInterval int
}

type PathConfig struct {
	Secrets string
}

func GetConfig(environment string) *StatusConfig {
	switch environment {
	case "LOCAL":
		return &StatusConfig{
			Mongo: MongoConfig{
				AuthEnabled:        false,
				Host:               "localhost",
				Port:               27017,
				ConnectionInterval: 5000,
			},
			Path: PathConfig{
				Secrets: "../instance",
			},
		}
	case "CLOUD":
		return &StatusConfig{
			Mongo: MongoConfig{
				AuthEnabled:        true,
				Host:               os.Getenv("MONGO_HOST"),
				Port:               27017,
				ConnectionInterval: 10000,
			},
			Path: PathConfig{
				Secrets: "/etc/secrets",
			},
		}
	case "DOCKER":
		return &StatusConfig{
			Mongo: MongoConfig{
				AuthEnabled:        false,
				Host:               "mongodb",
				Port:               27017,
				ConnectionInterval: 5000,
			},
		}
	default:
		log.Fatalf("Unknown environment: %s", environment)
		return nil
	}
}

func ReadCredential(basePath string, file string) (string, error) {
	path := filepath.Join(basePath, file)
	credential, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", file, err)
	}
	return string(credential), nil
}

func GetCredentials(secretsPath string) options.Credential {
	var cred options.Credential

	authSource, err := ReadCredential(secretsPath, "database")
	if err != nil {
		log.Printf("Failed to read database credential: %s", err.Error())
		authSource = ""
	}
	cred.AuthMechanism = "SCRAM-SHA-256"
	cred.AuthSource = authSource

	username, err := ReadCredential(secretsPath, "username")
	if err != nil {
		log.Printf("Failed to read username: %s", err.Error())
		username = ""
	}
	cred.Username = username

	password, err := ReadCredential(secretsPath, "password")
	if err != nil {
		log.Printf("Failed to read password: %s", err.Error())
		password = ""
	}
	cred.Password = password

	return cred
}
