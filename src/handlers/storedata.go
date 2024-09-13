package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"epl-fantasy/src/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StoreGameWeekData(client *mongo.Client, data *config.Data) error {
	collection := client.Database("fantasy_football").Collection("gameweek_data")

	filter := bson.M{"gameweek": data.GameWeek}
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

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
