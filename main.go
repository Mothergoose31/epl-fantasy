package main

import (
	"epl-fantasy/src/service"
	"fmt"
	"log"
)

func main() {
	fplService, err := service.NewFPLService()
	if err != nil {
		log.Fatalf("Error creating FPL service: %v", err)
	}

	data, body, err := fplService.FetchFPLData()
	if err != nil {
		log.Fatalf("Error fetching FPL data: %v", err)
	}

	fmt.Println("Raw API Response:")
	fmt.Println(string(body))

	fmt.Printf("\nFetched data for %d players\n", len(data.Elements))
	fmt.Printf("Latest GameWeek: %d\n", data.GameWeek)

	if len(data.Elements) > 0 {
		player := data.Elements[0]
		fmt.Printf("\nFirst player details:\n")
		fmt.Printf("Name: %s %s\n", player.FirstName, player.SecondName)
		fmt.Printf("Team: %d\n", player.Team)
		fmt.Printf("Position: %d\n", player.ElementType)
		fmt.Printf("Cost: %.1f\n", float64(player.NowCost)/10)
		fmt.Printf("Form: %s\n", player.Form)
		fmt.Printf("Total Points: %d\n", player.TotalPoints)
	}
}
