package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var Database *mongo.Database

func databaseFromConnectionString(connectionString string) string {
	tokens := strings.Split(connectionString, "/")
	return tokens[len(tokens)-1]
}

func Disconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func GetConnectionString() string {
	return os.Getenv("DB_CONNECTION_STRING")
}

func Connect(ctx context.Context) error {
	connectionString := GetConnectionString()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	} else {
		if err := client.Ping(ctx, nil); err != nil {
			return err
		}
		fmt.Println("successfully connected to the database!")
	}
	Database = client.Database(databaseFromConnectionString(connectionString))
	return nil
}
