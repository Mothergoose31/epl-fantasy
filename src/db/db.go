package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"epl-fantasy/src/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeMongoDB(configs *config.StatusConfig) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%d", config.App.Mongo.Host, config.App.Mongo.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	if config.App.Mongo.AuthEnabled {
		clientOptions = clientOptions.SetAuth(config.GetCredentials(config.App.Path.Secrets))
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ===============================================================
func getDatabaseName() string {
	if config.App.Mongo.AuthEnabled {
		databaseName, err := config.ReadCredential(config.App.Path.Secrets, "database")
		if err != nil {
			log.Printf("Failed to read database name: %v", err)
		}
		return databaseName
	}
	return "fantasy_football"
}

// ===============================================================

func GetCollection(collectionName string) *mongo.Collection {
	if config.Client != nil {
		return config.Client.Database(getDatabaseName()).Collection(collectionName)
	}
	return nil
}

// ===============================================================

func GetGameWeekCollection() *mongo.Collection {
	return GetCollection("gameweek_data")

}
