package handlers

import (
	"context"
	"encoding/json"
	"epl-fantasy/src/config"
	"epl-fantasy/src/db"
	"epl-fantasy/src/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertGameWeekData(client *mongo.Client, data *config.Data) error {
	collection := db.GetGameWeekCollection()
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

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func FetchAndStoreGameWeekData(w http.ResponseWriter, r *http.Request) {
	fplService, err := service.NewFPLService()
	if err != nil {
		log.Printf("Error creating FPL service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data, body, err := fplService.FetchFPLData()
	if err != nil {
		log.Printf("Error fetching FPL data: %v", err)
		str := fmt.Sprintf("Error fetching FPL data: %v", err)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	err = InsertGameWeekData(config.Client, data)
	if err != nil {
		log.Printf("Error storing game week data: %v", err)
		http.Error(w, fmt.Sprintf("Error storing game week data: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// ==================================================

func GetGameData(w http.ResponseWriter, r *http.Request) {
	collection := db.GetGameWeekCollection()
	if collection == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var gameWeekData []config.GameWeekData
	err = cursor.All(context.Background(), &gameWeekData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(gameWeekData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetBestPerformers(w http.ResponseWriter, r *http.Request) {
	collection := db.GetGameWeekCollection()
	if collection == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	goalkeepers, err := GetBestPerformersOverGameWeeks(collection, 1, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting goalkeepers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defenders, err := GetBestPerformersOverGameWeeks(collection, 2, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting defenders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	midfielders, err := GetBestPerformersOverGameWeeks(collection, 3, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting midfielders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	forwards, err := GetBestPerformersOverGameWeeks(collection, 4, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting forwards: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("==========================================Goalkeepers")
	for _, player := range goalkeepers {
		fmt.Printf("Player:%v\n TotalValue:%v, averageScore: %v", player.WebName, player.TotalPoints, player.AvgPoints)

	}
	fmt.Println("==========================================Defenders")
	for _, player := range defenders {
		fmt.Printf("Player:%v\n TotalValue:%v, averageScore: %v", player.WebName, player.TotalPoints, player.AvgPoints)
	}
	fmt.Println("==========================================Midfielders")
	for _, player := range midfielders {
		fmt.Printf("Player:%v\n TotalValue:%v , averageScore: %v ", player.WebName, player.TotalPoints, player.AvgPoints)
	}
	fmt.Println("==========================================Forwards")
	for _, player := range forwards {
		fmt.Printf("Player:%v\n TotalValue:%v, averageScore: %v", player.WebName, player.TotalPoints, player.AvgPoints)
	}
	limitPrice := 1030
	optimalTeam, err := CalculateOptimalTeam(limitPrice, goalkeepers, defenders, midfielders, forwards)
	if err != nil {
		http.Error(w, "Error calculating optimal team: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		OptimalTeam []config.PlayerPerformance `json:"optimalTeam"`
	}{
		OptimalTeam: optimalTeam,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
