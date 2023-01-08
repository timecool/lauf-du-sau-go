package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
)

func GetAllTimeLeaderboard(userCollection *mongo.Collection) ([]models.LeaderboardUser, error) {
	var results []models.LeaderboardUser
	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$group": bson.M{
			"_id":       "$_id",
			"username":  bson.M{"$first": "$username"},
			"email":     bson.M{"$first": "$email"},
			"image_url": bson.M{"$first": "$image_url"},
			"total":     bson.M{"$sum": "$runs.distance"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2})

	if err != nil {
		return results, err

	}

	if err = cursor.All(database.Ctx, &results); err != nil {
		return results, err

	}

	if err := cursor.Close(database.Ctx); err != nil {
		return results, err
	}
	return results, nil

}
func GetMonthLeaderboard(userCollection *mongo.Collection, month string) ([]models.LeaderboardUser, error) {
	var results []models.LeaderboardUser
	firstOfMonth, lastOfMonth, err := GetFirstAndLastDayFromMonth(month)

	o1 := bson.M{
		"$unwind": "$runs",
	}
	o2 := bson.M{
		"$match": bson.M{"runs.date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o3 := bson.M{
		"$group": bson.M{
			"_id":       "$_id",
			"username":  bson.M{"$first": "$username"},
			"email":     bson.M{"$first": "$email"},
			"image_url": bson.M{"$first": "$image_url"},
			"total":     bson.M{"$sum": "$runs.distance"},
		},
	}

	cursor, err := userCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3})

	if err != nil {
		return results, err

	}

	if err = cursor.All(database.Ctx, &results); err != nil {
		return results, err

	}
	if err := cursor.Close(database.Ctx); err != nil {
		return results, err
	}

	return results, nil

}
