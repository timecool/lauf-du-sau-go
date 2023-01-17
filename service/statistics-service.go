package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"lauf-du-sau/database"
	"lauf-du-sau/models"
)

func GetAllTimeLeaderboard() ([]models.LeaderboardUserID, error) {
	var results []models.LeaderboardUserID

	o1 := bson.M{
		"$match": bson.M{"status": models.RunActivate},
	}

	o2 := bson.M{
		"$group": bson.M{
			"_id":   "$user_id",
			"total": bson.M{"$sum": "$distance"},
		},
	}

	runCollection := database.InitRunCollection()
	cursor, err := runCollection.Aggregate(database.Ctx, []bson.M{o1, o2})

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
func GetMonthLeaderboard(month string) ([]models.LeaderboardUserID, error) {
	var results []models.LeaderboardUserID
	firstOfMonth, lastOfMonth, err := GetFirstAndLastDayFromMonth(month)

	o1 := bson.M{
		"$match": bson.M{"status": models.RunActivate},
	}
	o2 := bson.M{
		"$match": bson.M{"date": bson.M{
			"$gte": firstOfMonth,
			"$lte": lastOfMonth,
		}},
	}
	o3 := bson.M{
		"$group": bson.M{
			"_id":   "$user_id",
			"total": bson.M{"$sum": "$distance"},
		},
	}

	runCollection := database.InitRunCollection()
	cursor, err := runCollection.Aggregate(database.Ctx, []bson.M{o1, o2, o3})

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
