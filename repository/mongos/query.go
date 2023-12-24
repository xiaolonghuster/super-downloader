package mongos

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"super-downloader/repository/mongos/entity"
)

func (m *MongoManager) QueryAllImage(dbName, collectionName string) ([]entity.ImageInfo, error) {

	var cursor *mongo.Cursor
	defer func() {
		if cursor != nil {
			cursor.Close(context.Background())
		}
	}()

	var images []entity.ImageInfo
	var err error
	if cursor, err = m.Client.Database(dbName).Collection(collectionName).Find(context.Background(), bson.D{}); err == nil {
		if err = cursor.All(context.Background(), &images); err != nil {
			logrus.Errorf("mongo query all image exception: %v", err)
		}
	}
	return images, err
}
