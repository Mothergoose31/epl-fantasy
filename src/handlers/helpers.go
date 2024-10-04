package handlers

import (
	"context"
	"epl-fantasy/src/config"
	"errors"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var playerteams map[int]int
var teamswithTomanyPlayer []int

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

//   https://mathematicallysafe.wordpress.com/2018/07/08/fpl-analysis-the-impact-of-fixtures-on-player-performance/
// https://jinhyuncheong.com/jekyll/update/2018/12/26/Form_over_fixture.html

// LOOKING INTO INTO HOW MUCH WEIGHT/ CONSIDERATION TO GIVE WHEN SELECTING PLAYERS BASED ON THEIR FORM AND FIXTURES
//  BASED ON SOME READINGS, INDICATIONS SHOW FOR DEFENDERS AND AND GOALKEEPERS, FIXTURES ARE MORE IMPORTANT THAN FORM.
//  FOR MIDFIELDERS AND FORWARDS, FORM IS MORE IMPORTANT THAN FIXTURES.

//  FOR FOWARDS AND MIDFIELDERS. SELECT PLAYERS WITH THE BEST FORM , FOR DEFENDERS AND GOALKEEPERS , SELECT PLAYERS WIT THE BEST VALUE
//  THAT WOULD BE  FORM  OVER PRICE AT THE MOMENT

//	ADD CHECKS TO MAKE SURE TEAM DOES NOT HAVE MORE THAN 3 PLAYERS FROM THE SAME TEAM IN THIS A  GLOBAL MAP WHERE YOU ARE ADDING THE TEAM THE PLAYER IS IN AND CHECKING IF TEAM ALREADY OVER 3
//	IF WE ARE OVER BUDGET THEN  STARTING WITH  MIDFIELDERS SECOND TO FIRST WITH REGUARDS TO FORM AND CYCLE DOWN THE 3 OTHER DEFENDERS DO A REPLACE AND CHECK IF THE TEAM  IS STILL OUT OF BUDGET IF IN BUDGET RETURN TEAM ,
//
// IF NOT THEN CYCLE DOWN TO THE NEXT MIDFIELDER AND REPEAT, IF STILL OUT OF BUDGET CYCLE THROUGH FORWARDS STARTING  AGAIN WITH SECOND BEST AND MOVING DOWN TO THE THIRD BEST FORWARD AND REPEAT THE PROCESS UNTIL A TEAM IS FOUND THAT IS WITHIN BUDGET
// IF NO TEAM IS FOUND THEN RETURN AN ERROR MESSAGE THAT NO TEAM COULD BE FOUND WITHIN THE BUDGET
func CalculateOptimalTeam(limitPrice int, goalies, defenders, midfielders, forwards []config.PlayerPerformance) ([]config.PlayerPerformance, error) {
	fmt.Println("====Running the  numbers, picking optimal team ......")
	const totalPlayerPerTeamLimmit = 3
	totalCost := 0

	err := checkPlayerCount(goalies, defenders, midfielders, forwards)
	if err != nil {
		fmt.Println("Not enough players are being returned from db")
		return nil, err
	}

	sortPlayersByValue(goalies)
	sortPlayersByValue(defenders)
	sortPlayersByAveragePoints(midfielders)
	sortPlayersByAveragePoints(forwards)

	topGoalies := goalies[:1]
	topDefenders := defenders[:5]
	topMidfielders := midfielders[:5]
	topForwards := forwards[:3]

	selectedTeam := append(topGoalies, topDefenders...)
	selectedTeam = append(selectedTeam, topMidfielders...)
	selectedTeam = append(selectedTeam, topForwards...)

	totalCost = calculateTotalCost(topGoalies, topDefenders, topMidfielders, topForwards)
	playerteams = countplayersFromTeam(topGoalies, topDefenders, topMidfielders, topForwards)
	for !checkTeam(playerteams) || totalCost > limitPrice {
		fmt.Println("team has more than 3 players from the same team or is over budget")
		fmt.Println("updating team composition")

		if !checkTeam(playerteams) {
			selectedTeam, err = adjustTeamComposition(selectedTeam, goalies, defenders, midfielders, forwards)

		}
		if totalCost > limitPrice {

		}
		return nil, errors.New("Team has more than 3 players from the same team or is over budget")

	}

	return selectedTeam, nil
}

func sortPlayersByAveragePoints(players []config.PlayerPerformance) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].AvgPoints > players[j].AvgPoints
	})
}

func sortPlayersByValue(players []config.PlayerPerformance) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].ValueScore > players[j].ValueScore
	})
}

func checkPlayerCount(goalies, defenders, midfielders, forwards []config.PlayerPerformance) error {
	if len(goalies) < 1 || len(defenders) < 5 || len(midfielders) < 5 || len(forwards) < 3 {
		return errors.New("not enough players to select from")
	}
	return nil
}

func countplayersFromTeam(goalies, defenders, midfielders, forwards []config.PlayerPerformance) map[int]int {
	teamCount := make(map[int]int)

	for _, player := range goalies {
		teamCount[player.Team]++
	}
	for _, player := range defenders {
		teamCount[player.Team]++
	}
	for _, player := range midfielders {
		teamCount[player.Team]++
	}
	for _, player := range forwards {
		teamCount[player.Team]++
	}
	return teamCount
}

func calculateTotalCost(goalies, defenders, midfielders, forwards []config.PlayerPerformance) int {
	totalCost := 0
	for _, player := range goalies {
		totalCost += player.NowCost
	}
	for _, player := range defenders {
		totalCost += player.NowCost
	}
	for _, player := range midfielders {
		totalCost += player.NowCost
	}
	for _, player := range forwards {
		totalCost += player.NowCost
	}
	return totalCost
}

func checkTeam(map2 map[int]int) bool {
	for _, value := range map2 {
		if value > 3 {
			return false
		}
		if value < 3 {
			teamswithTomanyPlayer = append(teamswithTomanyPlayer, value)
			return true
		}
	}
	return true
}

func adjustTeamComposition(selectedTeam, goalies, defenders, midfielders, forwards []config.PlayerPerformance) ([]config.PlayerPerformance, error) {
	// teamswithTomanyPlayer []int
	// selectedTeam[0] = goalie
	// selectedTeam[1] = goalie
	// selectedTeam[2] = defender
	// selectedTeam[3] = defender
	// selectedTeam[4] = defender
	// selectedTeam[5] = defender
	// selectedTeam[6] = defender
	// selectedTeam[7] = midfielder
	// selectedTeam[8] = midfielder
	// selectedTeam[9] = midfielder
	// selectedTeam[10] = midfielder
	// selectedTeam[11] = midfielder
	// selectedTeam[12] = forward
	// selectedTeam[13] = forward
	// selectedTeam[14] = forward

	// var playerteams map[int]int
	// var teamswithTomanyPlayer []int

	//  start with defenders  position 3 and cycle through the other defenders
	//  if player[i].Team is  in teamswithTomanyPlayer []int then replace with a new player , replace with player found in coresponding player array
	// then suntract 1 from map and add 1 to the new player team in map
	// check if team in map is still over 3 comtinue
	//  if below 3 then remove from teamswithTomanyPlayer []int

}

// func calculateTotalCost(goalies, defenders, midfielders, forwards []config.PlayerPerformance) int {
// 	totalCost := 0
// 	for _, player := range goalies {
// 		totalCost += player.NowCost
// 	}
// 	for _, player := range defenders {
// 		totalCost += player.NowCost
// 	}
// 	for _, player := range midfielders {
// 		totalCost += player.NowCost
// 	}
// 	for _, player := range forwards {
// 		totalCost += player.NowCost
// 	}
// 	return totalCost
// }
