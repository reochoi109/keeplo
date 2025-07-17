package logdata_repo

import (
	"context"
	"keeplo/internal/domain/logdata"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoLogRepo(col *mongo.Collection) logdata.Repository {
	return &mongoRepo{collection: col}
}

func (r *mongoRepo) GetHealthLogs(ctx context.Context, monitorID string, from, to time.Time, statusCode *int, isSuccess *bool) ([]logdata.MonitorHealthLog, error) {
	filter := bson.M{
		"monitor_id": monitorID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	if statusCode != nil {
		filter["status_code"] = *statusCode
	}
	if isSuccess != nil {
		filter["is_success"] = *isSuccess
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []logdata.MonitorHealthLog
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
