package config

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// * Return connection to collection
func GetCollection() *mongo.Collection {
	_, databaseName, collectionName := GetClientDetails()
	theCollection := Client.Database(databaseName).Collection(collectionName)

	return theCollection
}
