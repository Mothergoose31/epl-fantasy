package handlers

import (
	"context"
	"epl-fantasy/src/config"
	"errors"
	"fmt"
	"math"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// hhamilton.typepad.com/files/pythag_mit_sa_2010.pdf
// https://blogs.salford.ac.uk/business-school/wp-content/uploads/sites/7/2016/09/paper.pdf
//

//  using Weibull distrubition first paper is describes win perentage as a function of goals scored and goals allowed

func PythagoreanExpectation(goalsScored, goalsAllowed float64, pythagoreanExponent float64) float64 {
	gamma := pythagoreanExponent
	kappa := math.Gamma(1 + 1/gamma)
	N := 10

	winProbability := math.Pow(goalsScored, gamma) / (math.Pow(goalsScored, gamma) + math.Pow(goalsAllowed, gamma))

	drawProbability := 0.0
	for c := 0; c <= N; c++ {
		cFloat := float64(c)
		gsProb := math.Exp(-math.Pow(kappa*(cFloat+1)/goalsScored, gamma)) - math.Exp(-math.Pow(kappa*cFloat/goalsScored, gamma))
		gaProb := math.Exp(-math.Pow(kappa*(cFloat+1)/goalsAllowed, gamma)) - math.Exp(-math.Pow(kappa*cFloat/goalsAllowed, gamma))
		drawProbability += gsProb * gaProb
	}

	expectedPoints := 3*winProbability + drawProbability

	return expectedPoints
}

func GetBestPerformersOverGameWeeks(collection *mongo.Collection, position int, startGameWeek, endGameWeek, limit int) ([]config.PlayerPerformance, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"game_week": bson.M{"$gte": startGameWeek, "$lte": endGameWeek}}}},
		{{Key: "$unwind", Value: "$players"}},
		{{Key: "$match", Value: bson.M{"players.element_type": position}}},
		{{Key: "$group", Value: bson.M{
			"_id":                 "$players.id",
			"web_name":            bson.M{"$first": "$players.web_name"},
			"total_points":        bson.M{"$sum": "$players.event_points"},
			"avg_points":          bson.M{"$avg": "$players.event_points"},
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
			"performance_score": bson.M{
				"$add": []interface{}{
					"$total_points",
					bson.M{"$multiply": []interface{}{"$goals_scored", 5}},
					bson.M{"$multiply": []interface{}{"$assists", 3}},
					"$clean_sheets",
					bson.M{"$multiply": []interface{}{"$saves", 0.5}},
					"$bonus",
					bson.M{"$divide": []interface{}{"$bps", 10}},
				},
			},
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

func calculatePlayerScore(player config.PlayerPerformance) float64 {
	// This is a simplified scoring function. You may want to adjust the weights based on your specific requirements.
	return float64(player.TotalPoints) +
		player.AvgPoints*5 +
		float64(player.GoalsScored)*4 +
		float64(player.Assists)*3 +
		float64(player.CleanSheets)*4 +
		float64(player.Saves)*0.5 +
		float64(player.Bonus)*2 +
		player.Influence*0.5 +
		player.Creativity*0.5 +
		player.Threat*0.5 +
		player.ICTIndex*2 +
		player.ExpectedGoals*3 +
		player.ExpectedAssists*2 +
		player.PerformanceScore*2 +
		player.ValueScore*2
}

// Greedy aproach
func CalculateOptimalTeam(LimitPrice int, goalies, defenders, midfielders, forwards []config.PlayerPerformance) ([]config.PlayerPerformance, error) {
	fmt.Println("================HITTING CALCULATE OPTIMAL=====================")
	if len(goalies) < 2 {
		return nil, fmt.Errorf("not enough goalkeepers: need 2, have %d", len(goalies))
	}
	if len(defenders) < 4 {
		return nil, fmt.Errorf("not enough defenders: need 4, have %d", len(defenders))
	}
	if len(midfielders) < 4 {
		return nil, fmt.Errorf("not enough midfielders: need 4, have %d", len(midfielders))
	}
	if len(forwards) < 3 {
		return nil, fmt.Errorf("not enough forwards: need 3, have %d", len(forwards))
	}

	// Sort players by score/price ratio
	sort.Slice(goalies, func(i, j int) bool {
		fmt.Println(goalies[i].NowCost)
		return calculatePlayerScore(goalies[i])/float64(goalies[i].NowCost) > calculatePlayerScore(goalies[j])/float64(goalies[j].NowCost)
	})
	sort.Slice(defenders, func(i, j int) bool {
		return calculatePlayerScore(defenders[i])/float64(defenders[i].NowCost) > calculatePlayerScore(defenders[j])/float64(defenders[j].NowCost)
	})
	sort.Slice(midfielders, func(i, j int) bool {
		return calculatePlayerScore(midfielders[i])/float64(midfielders[i].NowCost) > calculatePlayerScore(midfielders[j])/float64(midfielders[j].NowCost)
	})
	sort.Slice(forwards, func(i, j int) bool {
		return calculatePlayerScore(forwards[i])/float64(forwards[i].NowCost) > calculatePlayerScore(forwards[j])/float64(forwards[j].NowCost)
	})

	// Initialize dynamic programming table
	dp := make([][]float64, LimitPrice+1)
	for i := range dp {
		dp[i] = make([]float64, 14)
	}

	// // Fill the dynamic programming table
	for i := 1; i <= LimitPrice; i++ {
		for j := 1; j <= 13; j++ {
			dp[i][j] = dp[i][j-1]
			if j <= 2 && j-1 < len(goalies) && i >= goalies[j-1].NowCost {
				dp[i][j] = max(dp[i][j], dp[i-goalies[j-1].NowCost][j-1]+calculatePlayerScore(goalies[j-1]))
			} else if j > 2 && j <= 6 && j-3 < len(defenders) && i >= defenders[j-3].NowCost {
				dp[i][j] = max(dp[i][j], dp[i-defenders[j-3].NowCost][j-1]+calculatePlayerScore(defenders[j-3]))
			} else if j > 6 && j <= 10 && j-7 < len(midfielders) && i >= midfielders[j-7].NowCost {
				dp[i][j] = max(dp[i][j], dp[i-midfielders[j-7].NowCost][j-1]+calculatePlayerScore(midfielders[j-7]))
			} else if j > 10 && j <= 13 && j-11 < len(forwards) && i >= forwards[j-11].NowCost {
				dp[i][j] = max(dp[i][j], dp[i-forwards[j-11].NowCost][j-1]+calculatePlayerScore(forwards[j-11]))
			}
		}
	}

	// // Backtrack to find the selected players
	selected := make([]config.PlayerPerformance, 0, 13)
	i, j := LimitPrice, 13
	for j > 0 {
		if j <= 2 && j <= len(goalies) {
			if i >= goalies[j-1].NowCost && dp[i][j] == dp[i-goalies[j-1].NowCost][j-1]+calculatePlayerScore(goalies[j-1]) {
				selected = append(selected, goalies[j-1])
				i -= goalies[j-1].NowCost
			}
		} else if j <= 6 && j-3 < len(defenders) {
			if i >= defenders[j-3].NowCost && dp[i][j] == dp[i-defenders[j-3].NowCost][j-1]+calculatePlayerScore(defenders[j-3]) {
				selected = append(selected, defenders[j-3])
				i -= defenders[j-3].NowCost
			}
		} else if j <= 10 && j-7 < len(midfielders) {
			if i >= midfielders[j-7].NowCost && dp[i][j] == dp[i-midfielders[j-7].NowCost][j-1]+calculatePlayerScore(midfielders[j-7]) {
				selected = append(selected, midfielders[j-7])
				i -= midfielders[j-7].NowCost
			}
		} else if j <= 13 && j-11 < len(forwards) {
			if i >= forwards[j-11].NowCost && dp[i][j] == dp[i-forwards[j-11].NowCost][j-1]+calculatePlayerScore(forwards[j-11]) {
				selected = append(selected, forwards[j-11])
				i -= forwards[j-11].NowCost
			}
		}
		j--
	}
	//  print the selected players

	// // Check if we have the correct number of players in each position
	goaliesCount := 0
	defendersCount := 0
	midfieldersCount := 0
	forwardsCount := 0
	for _, player := range selected {
		switch {
		case player.ID == goalies[0].ID || player.ID == goalies[1].ID:
			goaliesCount++
		case player.ID == defenders[0].ID || player.ID == defenders[1].ID || player.ID == defenders[2].ID || player.ID == defenders[3].ID:
			defendersCount++
		case player.ID == midfielders[0].ID || player.ID == midfielders[1].ID || player.ID == midfielders[2].ID || player.ID == midfielders[3].ID:
			midfieldersCount++
		case player.ID == forwards[0].ID || player.ID == forwards[1].ID || player.ID == forwards[2].ID:
			forwardsCount++
		}
	}

	if goaliesCount != 2 || defendersCount != 4 || midfieldersCount != 4 || forwardsCount != 3 {
		return nil, errors.New("could not select the required number of players for each position within the given budget")
	}

	// Reverse the selected slice to get the correct order
	for i, j := 0, len(selected)-1; i < j; i, j = i+1, j-1 {
		selected[i], selected[j] = selected[j], selected[i]
	}

	return selected, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

type Node struct {
	players         []config.PlayerPerformance
	selectedPlayers []bool
	bound           float64
	value           float64
	cost            int
	positionCounts  map[string]int
}

// TODO Look into Branch and Bound Algo implementation
