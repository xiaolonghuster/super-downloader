package mongos

import (
	"context"
	"github.com/sirupsen/logrus"
)

func (m *MongoManager) SaveDocuments(dbName, collectionName string, documents []interface{}) ([]interface{}, error) {
	var err error
	if result, err := m.Client.Database(dbName).Collection(collectionName).InsertMany(context.Background(), documents); err == nil {
		return result.InsertedIDs, nil
	}
	logrus.Errorf("mongodb save document failure. exception:%v", err)
	return nil, err
}
