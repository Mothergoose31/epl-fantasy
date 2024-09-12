package db

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func StoreGameWeekData(client *mongo.Client, data *Data) error {
	collection := client.Database("fantasy_football").Collection("gameweek_data")

	gameWeekData := GameWeekData{
		GameWeek:  data.GameWeek,
		Season:    "2024-2025",
		Timestamp: time.Now(),
		Players:   make([]PlayerSnapshot, len(data.Elements)),
	}

	for i, element := range data.Elements {
		gameWeekData.Players[i] = PlayerSnapshot{
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

	_, err := collection.InsertOne(context.Background(), gameWeekData)
	return err
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
