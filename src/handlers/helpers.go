package handlers

import (
	"epl-fantasy/src/config"
	"errors"
	"fmt"
	"sort"
)

var playerteams map[int]int
var teamswithTomanyPlayer []int

// TODO:
// 1. CALCULATE OPTIMAL TEAM WITH THE OPTION OF INCLUDING OR  EXCLUDING PLAYERS

//   https://mathematicallysafe.wordpress.com/2018/07/08/fpl-analysis-the-impact-of-fixtures-on-player-performance/
// https://jinhyuncheong.com/jekyll/update/2018/12/26/Form_over_fixture.html

// LOOKING INTO HOW MUCH WEIGHT/ CONSIDERATION TO GIVE WHEN SELECTING PLAYERS BASED ON THEIR FORM AND FIXTURES
//  BASED ON SOME READINGS, INDICATIONS SHOW FOR DEFENDERS AND AND GOALKEEPERS, FIXTURES ARE MORE IMPORTANT TH	AN FORM.
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

	topGoalies := goalies[:2]
	topDefenders := defenders[:5]
	topMidfielders := midfielders[:5]
	topForwards := forwards[:3]

	selectedTeam := append(topGoalies, topDefenders...)
	selectedTeam = append(selectedTeam, topMidfielders...)
	selectedTeam = append(selectedTeam, topForwards...)
	fmt.Println("Debug: Selected team composition:")

	positionCounts := make(map[int]int)
	for _, player := range selectedTeam {
		positionCounts[player.ElementType]++
		fmt.Printf("Player: %s, Position: %d\n", player.WebName, player.ElementType)
	}

	fmt.Println("Debug: Position counts:")
	for pos, count := range positionCounts {
		fmt.Printf("Position %d: %d players\n", pos, count)
	}

	if positionCounts[1] != 2 {
		fmt.Println("Warning: Incorrect number of goalkeepers selected!")
		// You might want to add additional logic here to handle this case
	}

	totalCost = calculateTotalCost(selectedTeam)
	playerteams = countplayersFromTeam(selectedTeam)
	for !checkTeam(playerteams) || totalCost > limitPrice {
		fmt.Println("team has more than 3 players from the same team or is over budget")
		fmt.Println("updating team composition")

		if !checkTeam(playerteams) {
			selectedTeam, err = adjustTeamComposition(selectedTeam, goalies, defenders, midfielders, forwards)
			if err != nil {
				return nil, err
			}
		}

		totalCost = calculateTotalCost(selectedTeam)
		if totalCost > limitPrice {
			selectedTeam, err = adjustTeamCompWithBudget(selectedTeam, goalies, defenders, midfielders, forwards, limitPrice)
			if err != nil {
				return nil, err
			}
		}

		playerteams = countplayersFromTeam(selectedTeam)
	}

	return selectedTeam, nil
}

// =========================================================================================================================================

func sortPlayersByAveragePoints(players []config.PlayerPerformance) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].AvgPoints > players[j].AvgPoints
	})
}

// =========================================================================================================================================

func sortPlayersByValue(players []config.PlayerPerformance) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].ValueScore > players[j].ValueScore
	})
}

// =========================================================================================================================================

func checkPlayerCount(goalies, defenders, midfielders, forwards []config.PlayerPerformance) error {
	if len(goalies) < 1 || len(defenders) < 5 || len(midfielders) < 5 || len(forwards) < 3 {
		return errors.New("not enough players to select from")
	}
	return nil
}

// =========================================================================================================================================

func calculateTotalCost(players []config.PlayerPerformance) int {
	totalCost := 0
	for _, player := range players {
		totalCost += player.NowCost
	}
	return totalCost
}

// =========================================================================================================================================

func countplayersFromTeam(selectedTeam []config.PlayerPerformance) map[int]int {
	teamCount := make(map[int]int)
	for _, player := range selectedTeam {
		teamCount[player.Team]++
	}
	return teamCount
}

// =========================================================================================================================================

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

// =========================================================================================================================================

func adjustTeamComposition(selectedTeam, goalies, defenders, midfielders, forwards []config.PlayerPerformance) ([]config.PlayerPerformance, error) {
	playerTeams := countplayersFromTeam(selectedTeam)

	teamsWithTooManyPlayers := []int{}

	for team, count := range playerTeams {
		fmt.Println("team", team, "count", count)
		if count > 3 {
			teamsWithTooManyPlayers = append(teamsWithTooManyPlayers, team)
		}
	}
	// 0-1 goalies
	// 2-6 defenders
	// 7-11 midfielders
	// 12-14 forwards

	replacePlayer := func(position int, availablePlayers []config.PlayerPerformance) bool {
		currentPlayer := selectedTeam[position]
		for _, newPlayer := range availablePlayers {
			if !contains(teamsWithTooManyPlayers, newPlayer.Team) && newPlayer.ID != currentPlayer.ID {
				playerTeams[currentPlayer.Team]--
				playerTeams[newPlayer.Team]++
				selectedTeam[position] = newPlayer

				if playerTeams[currentPlayer.Team] <= 3 {
					teamsWithTooManyPlayers = remove(teamsWithTooManyPlayers, currentPlayer.Team)

				}
				if playerTeams[newPlayer.Team] > 3 {
					teamsWithTooManyPlayers = append(teamsWithTooManyPlayers, newPlayer.Team)
				}
				return true
			}
		}
		return false
	}

	for i := 2; i <= 6; i++ {
		if contains(teamsWithTooManyPlayers, selectedTeam[i].Team) {
			if !replacePlayer(i, defenders) {
				return nil, fmt.Errorf("unable to adjust team composition for defenders")
			}
		}
	}

	for i := 7; i <= 11; i++ {
		if contains(teamsWithTooManyPlayers, selectedTeam[i].Team) {
			if !replacePlayer(i, midfielders) {
				return nil, fmt.Errorf("unable to adjust team composition for midfielders")
			}
		}
	}

	for i := 12; i <= len(selectedTeam)-1; i++ {
		if contains(teamsWithTooManyPlayers, selectedTeam[i].Team) {
			if !replacePlayer(i, forwards) {
				return nil, fmt.Errorf("unable to adjust team composition for forwards")
			}
		}
	}

	for i := 0; i <= 1; i++ {
		if contains(teamsWithTooManyPlayers, selectedTeam[i].Team) {
			if !replacePlayer(i, goalies) {
				return nil, fmt.Errorf("unable to adjust team composition for goalies")
			}
		}
	}

	if len(teamsWithTooManyPlayers) > 0 {
		return nil, fmt.Errorf("unable to adjust team composition to meet requirements")
	}

	return selectedTeam, nil

}

// =========================================================================================================================================

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// =========================================================================================================================================

func remove(slice []int, item int) []int {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// =========================================================================================================================================

func adjustTeamCompWithBudget(selectedTeam, goalies, defenders, midfielders, forwards []config.PlayerPerformance, limitPrice int) ([]config.PlayerPerformance, error) {

	currentBudget := calculateTotalCost(selectedTeam)
	budgetDeficit := currentBudget - limitPrice

	if budgetDeficit <= 0 {
		return selectedTeam, nil
	}

	countTeamPlayers := func() map[int]int {
		teamCount := make(map[int]int)
		for _, player := range selectedTeam {
			teamCount[player.Team]++
		}
		return teamCount
	}

	tryReplace := func(index int, candidates []config.PlayerPerformance) bool {
		teamCount := countTeamPlayers()
		currentPlayer := selectedTeam[index]
		for _, candidate := range candidates {
			if candidate.NowCost < currentPlayer.NowCost &&
				teamCount[candidate.Team] < 3 &&
				candidate.NowCost-currentPlayer.NowCost+budgetDeficit <= 0 {
				selectedTeam[index] = candidate
				budgetDeficit += candidate.NowCost - currentPlayer.NowCost
				return true
			}
		}
		return false
	}

	positionRanges := []struct {
		start, end int
		candidates []config.PlayerPerformance
	}{
		{2, 6, defenders},
		{7, 11, midfielders},
		{0, 1, goalies},
		{12, len(selectedTeam) - 1, forwards},
	}

	for _, pr := range positionRanges {
		for i := pr.start; i <= pr.end; i++ {
			if tryReplace(i, pr.candidates) {
				if budgetDeficit <= 0 {
					return selectedTeam, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("unable to adjust team composition to meet budget requirements")
}
