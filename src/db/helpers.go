package db

import (
	"context"
	"epl-fantasy/src/config"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// =========================================================================================================================================

func InsertGameWeekData(client *mongo.Client, data *config.Data) error {
	collection := GetGameWeekCollection()
	if collection == nil {
		return fmt.Errorf("error getting collection")
	}

	filter := bson.M{"game_week": data.GameWeek}
	var existingData config.GameWeekData
	err := collection.FindOne(context.Background(), filter).Decode(&existingData)
	if err == nil {

		return fmt.Errorf("data for game week %d already exists", data.GameWeek)
	} else if err != mongo.ErrNoDocuments {

		return fmt.Errorf("error checking existing data: %w", err)
	}
	gameWeekData := config.GameWeekData{
		GameWeek:  data.GameWeek,
		Season:    "2024-2025",
		Timestamp: time.Now(),
		Players:   make([]config.PlayerSnapshot, len(data.Elements)),
	}

	for i, element := range data.Elements {
		gameWeekData.Players[i] = config.PlayerSnapshot{
			ID:                   element.ID,
			GameWeek:             data.GameWeek,
			FirstName:            element.FirstName,
			SecondName:           element.SecondName,
			WebName:              element.WebName,
			Team:                 element.Team,
			ElementType:          element.ElementType,
			TotalPoints:          element.TotalPoints,
			EventPoints:          element.EventPoints,
			NowCost:              element.NowCost,
			Form:                 parseFloat(element.Form),
			SelectedByPercent:    parseFloat(element.SelectedByPercent),
			Minutes:              element.Minutes,
			GoalsScored:          element.GoalsScored,
			Assists:              element.Assists,
			CleanSheets:          element.CleanSheets,
			GoalsConceded:        element.GoalsConceded,
			OwnGoals:             element.OwnGoals,
			PenaltiesSaved:       element.PenaltiesSaved,
			PenaltiesMissed:      element.PenaltiesMissed,
			YellowCards:          element.YellowCards,
			RedCards:             element.RedCards,
			Saves:                element.Saves,
			Bonus:                element.Bonus,
			Bps:                  element.Bps,
			Influence:            parseFloat(element.Influence),
			Creativity:           parseFloat(element.Creativity),
			Threat:               parseFloat(element.Threat),
			IctIndex:             parseFloat(element.IctIndex),
			ExpectedGoals:        parseFloat(element.ExpectedGoals),
			ExpectedAssists:      parseFloat(element.ExpectedAssists),
			ExpectedGoalsPer90:   element.ExpectedGoalsPer90,
			SavesPer90:           element.SavesPer90,
			ExpectedAssistsPer90: element.ExpectedAssistsPer90,
		}
	}

	_, err = collection.InsertOne(context.Background(), gameWeekData)
	if err != nil {
		return fmt.Errorf("error inserting game week data: %w", err)
	}
	return nil
}

// =========================================================================================================================================

func GetBestPerformersOverGameWeeks(collection *mongo.Collection, position int, startGameWeek, endGameWeek, limit int) ([]config.PlayerPerformance, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"game_week":            bson.M{"$gte": startGameWeek, "$lte": endGameWeek},
			"players.element_type": position,
		}}},
		{{Key: "$unwind", Value: "$players"}},
		{{Key: "$match", Value: bson.M{"players.element_type": position}}},
		{{Key: "$group", Value: bson.M{
			"_id":                 "$players.id",
			"web_name":            bson.M{"$first": "$players.web_name"},
			"total_points":        bson.M{"$sum": "$players.event_points"},
			"avg_points":          bson.M{"$avg": "$players.event_points"},
			"team":                bson.M{"$first": "$players.team"},
			"element_type":        bson.M{"$first": "$players.element_type"},
			"goals_scored":        bson.M{"$sum": "$players.goals_scored"},
			"assists":             bson.M{"$sum": "$players.assists"},
			"clean_sheets":        bson.M{"$sum": "$players.clean_sheets"},
			"goals_conceded":      bson.M{"$sum": "$players.goals_conceded"},
			"saves":               bson.M{"$sum": "$players.saves"},
			"bonus":               bson.M{"$sum": "$players.bonus"},
			"bps":                 bson.M{"$sum": "$players.bps"},
			"influence":           bson.M{"$avg": "$players.influence"},
			"creativity":          bson.M{"$avg": "$players.creativity"},
			"threat":              bson.M{"$avg": "$players.threat"},
			"ict_index":           bson.M{"$avg": "$players.ict_index"},
			"expected_goals":      bson.M{"$sum": "$players.expected_goals"},
			"expected_assists":    bson.M{"$sum": "$players.expected_assists"},
			"now_cost":            bson.M{"$last": "$players.now_cost"},
			"selected_by_percent": bson.M{"$last": "$players.selected_by_percent"},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"base_score": bson.M{"$add": []interface{}{
				"$total_points",
				bson.M{"$multiply": []interface{}{"$avg_points", 5}},
				bson.M{"$multiply": []interface{}{"$bonus", 2}},
				bson.M{"$divide": []interface{}{"$bps", 20}},
			}},
			"ict_score": bson.M{"$add": []interface{}{
				bson.M{"$multiply": []interface{}{"$influence", 0.3}},
				bson.M{"$multiply": []interface{}{"$creativity", 0.3}},
				bson.M{"$multiply": []interface{}{"$threat", 0.3}},
				bson.M{"$multiply": []interface{}{"$ict_index", 0.1}},
			}},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"performance_score": bson.M{"$switch": bson.M{
				"branches": []interface{}{
					bson.M{"case": bson.M{"$eq": []interface{}{position, 1}}, "then": bson.M{"$add": []interface{}{
						"$base_score",
						bson.M{"$multiply": []interface{}{"$clean_sheets", 4}},
						bson.M{"$multiply": []interface{}{"$saves", 0.5}},
						bson.M{"$multiply": []interface{}{"$goals_conceded", -0.5}},
						bson.M{"$multiply": []interface{}{"$ict_score", 0.5}},
					}}},
					bson.M{"case": bson.M{"$eq": []interface{}{position, 2}}, "then": bson.M{"$add": []interface{}{
						"$base_score",
						bson.M{"$multiply": []interface{}{"$clean_sheets", 4}},
						bson.M{"$multiply": []interface{}{"$goals_scored", 6}},
						bson.M{"$multiply": []interface{}{"$assists", 3}},
						bson.M{"$multiply": []interface{}{"$goals_conceded", -0.1}},
						bson.M{"$multiply": []interface{}{"$ict_score", 0.7}},
					}}},
					bson.M{"case": bson.M{"$eq": []interface{}{position, 3}}, "then": bson.M{"$add": []interface{}{
						"$base_score",
						bson.M{"$multiply": []interface{}{"$goals_scored", 5}},
						bson.M{"$multiply": []interface{}{"$assists", 3}},
						"$clean_sheets",
						"$ict_score",
						bson.M{"$multiply": []interface{}{"$expected_goals", 2}},
						bson.M{"$multiply": []interface{}{"$expected_assists", 2}},
					}}},
					bson.M{"case": bson.M{"$eq": []interface{}{position, 4}}, "then": bson.M{"$add": []interface{}{
						"$base_score",
						bson.M{"$multiply": []interface{}{"$goals_scored", 4}},
						bson.M{"$multiply": []interface{}{"$assists", 2}},
						bson.M{"$multiply": []interface{}{"$ict_score", 1.2}},
						bson.M{"$multiply": []interface{}{"$expected_goals", 3}},
						bson.M{"$multiply": []interface{}{"$expected_assists", 1.5}},
					}}},
				},
				"default": "$base_score",
			}},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"value_score": bson.M{"$divide": []interface{}{"$total_points", "$now_cost"}},
		}}},
		{{Key: "$sort", Value: bson.M{"performance_score": -1}}},
		{{Key: "$limit", Value: limit}},
	}

	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []config.PlayerPerformance
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// =========================================================================================================================================

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
