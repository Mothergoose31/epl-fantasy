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

	"go.mongodb.org/mongo-driver/bson"
)

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

	err = db.InsertGameWeekData(config.Client, data)
	if err != nil {
		log.Printf("Error storing game week data: %v", err)
		http.Error(w, fmt.Sprintf("Error storing game week data: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// =========================================================================================================================================

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

// =========================================================================================================================================

func GetBestPerformers(w http.ResponseWriter, r *http.Request) {
	collection := db.GetGameWeekCollection()
	if collection == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	goalkeepers, err := db.GetBestPerformersOverGameWeeks(collection, 1, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting goalkeepers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, player := range goalkeepers {
		fmt.Printf("Player:%v\n TotalValue:%v, averageScore: %v", player.WebName, player.TotalPoints, player.AvgPoints)
	}

	defenders, err := db.GetBestPerformersOverGameWeeks(collection, 2, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting defenders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	midfielders, err := db.GetBestPerformersOverGameWeeks(collection, 3, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting midfielders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	forwards, err := db.GetBestPerformersOverGameWeeks(collection, 4, 3, 6, 20)
	if err != nil {
		http.Error(w, "Error getting forwards: "+err.Error(), http.StatusInternalServerError)
		return
	}

	limitPrice := 1030
	optimalTeam, err := CalculateOptimalTeam(limitPrice, goalkeepers, defenders, midfielders, forwards)
	if err != nil {
		http.Error(w, "Error calculating optimal team: "+err.Error(), http.StatusInternalServerError)
		return
	}

	optimalTeamResponse := config.OptimalTeam{
		TotalCost:   calculateTotalCost(optimalTeam),
		Goalkeepers: optimalTeam[0:2],
		Defenders:   optimalTeam[2:7],
		Midfielders: optimalTeam[7:12],
		Forwards:    optimalTeam[12:15],
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(optimalTeamResponse); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
