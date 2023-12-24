package mongos

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//type Handler struct {
//
//}

func Init(mongoURI, dbName string) *MongoManager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	return &MongoManager{
		Client: client,
	}
}

func (m *MongoManager) Disconnect() {
	m.Client.Disconnect(context.Background())
}
