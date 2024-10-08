package config

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	App    *StatusConfig
	Client *mongo.Client
)

// ======================EPL API Structs ======================
// ============================================================

type Data struct {
	Events       []Event       `json:"events"`
	GameSettings GameSettings  `json:"game_settings"`
	Teams        []Team        `json:"teams"`
	TotalPlayers int64         `json:"total_players"` // Change here
	Elements     []Element     `json:"elements"`
	ElementTypes []ElementType `json:"element_types"`
	GameWeek     int           `json:"game_week,omitempty"`
}

type Event struct {
	ID                     int        `json:"id"`
	Name                   string     `json:"name"`
	DeadlineTime           time.Time  `json:"deadline_time"`
	ReleaseTime            *time.Time `json:"release_time"`
	AverageEntryScore      int        `json:"average_entry_score"`
	Finished               bool       `json:"finished"`
	DataChecked            bool       `json:"data_checked"`
	HighestScoringEntry    int        `json:"highest_scoring_entry"`
	DeadlineTimeEpoch      int64      `json:"deadline_time_epoch"`
	DeadlineTimeGameOffset int        `json:"deadline_time_game_offset"`
	HighestScore           int        `json:"highest_score"`
	IsPrevious             bool       `json:"is_previous"`
	IsCurrent              bool       `json:"is_current"`
	IsNext                 bool       `json:"is_next"`
	CupLeaguesCreated      bool       `json:"cup_leagues_created"`
	H2hKoMatchesCreated    bool       `json:"h2h_ko_matches_created"`
	RankedCount            int        `json:"ranked_count"`
	ChipPlays              []ChipPlay `json:"chip_plays"`
	MostSelected           int        `json:"most_selected"`
	MostTransferredIn      int        `json:"most_transferred_in"`
	TopElement             int        `json:"top_element"`
	TopElementInfo         TopElement `json:"top_element_info"`
	TransfersMade          int        `json:"transfers_made"`
	MostCaptained          int        `json:"most_captained"`
	MostViceCaptained      int        `json:"most_vice_captained"`
}

// ==================================

type ChipPlay struct {
	ChipName  string `json:"chip_name"`
	NumPlayed int    `json:"num_played"`
}

// ==================================

type TopElement struct {
	ID     int `json:"id"`
	Points int `json:"points"`
}

// ==================================

type GameSettings struct {
	LeagueJoinPrivateMax         int      `json:"league_join_private_max"`
	LeagueJoinPublicMax          int      `json:"league_join_public_max"`
	LeagueMaxSizePublicClassic   int      `json:"league_max_size_public_classic"`
	LeagueMaxSizePublicH2h       int      `json:"league_max_size_public_h2h"`
	LeagueMaxSizePrivateH2h      int      `json:"league_max_size_private_h2h"`
	LeagueMaxKoRoundsPrivateH2h  int      `json:"league_max_ko_rounds_private_h2h"`
	LeaguePrefixPublic           string   `json:"league_prefix_public"`
	LeaguePointsH2hWin           int      `json:"league_points_h2h_win"`
	LeaguePointsH2hLose          int      `json:"league_points_h2h_lose"`
	LeaguePointsH2hDraw          int      `json:"league_points_h2h_draw"`
	LeagueKoFirstInsteadOfRandom bool     `json:"league_ko_first_instead_of_random"`
	CupStartEventID              *int     `json:"cup_start_event_id"`
	CupStopEventID               *int     `json:"cup_stop_event_id"`
	CupQualifyingMethod          *string  `json:"cup_qualifying_method"`
	CupType                      *string  `json:"cup_type"`
	FeaturedEntries              []int    `json:"featured_entries"`
	PercentileRanks              []int    `json:"percentile_ranks"`
	SquadSquadplay               int      `json:"squad_squadplay"`
	SquadSquadsize               int      `json:"squad_squadsize"`
	SquadTeamLimit               int      `json:"squad_team_limit"`
	SquadTotalSpend              int      `json:"squad_total_spend"`
	UICurrencyMultiplier         int      `json:"ui_currency_multiplier"`
	UIUseSpecialShirts           bool     `json:"ui_use_special_shirts"`
	UISpecialShirtExclusions     []int    `json:"ui_special_shirt_exclusions"`
	StatsFormDays                int      `json:"stats_form_days"`
	SysViceCaptainEnabled        bool     `json:"sys_vice_captain_enabled"`
	TransfersCap                 int      `json:"transfers_cap"`
	TransfersSellOnFee           float64  `json:"transfers_sell_on_fee"`
	MaxExtraFreeTransfers        int      `json:"max_extra_free_transfers"`
	LeagueH2hTiebreakStats       []string `json:"league_h2h_tiebreak_stats"`
	Timezone                     string   `json:"timezone"`
}

// ==================================

type Team struct {
	Code                int     `json:"code"`
	Draw                int     `json:"draw"`
	Form                *string `json:"form"`
	ID                  int     `json:"id"`
	Loss                int     `json:"loss"`
	Name                string  `json:"name"`
	Played              int     `json:"played"`
	Points              int     `json:"points"`
	Position            int     `json:"position"`
	ShortName           string  `json:"short_name"`
	Strength            int     `json:"strength"`
	TeamDivision        *int    `json:"team_division"`
	Unavailable         bool    `json:"unavailable"`
	Win                 int     `json:"win"`
	StrengthOverallHome int     `json:"strength_overall_home"`
	StrengthOverallAway int     `json:"strength_overall_away"`
	StrengthAttackHome  int     `json:"strength_attack_home"`
	StrengthAttackAway  int     `json:"strength_attack_away"`
	StrengthDefenceHome int     `json:"strength_defence_home"`
	StrengthDefenceAway int     `json:"strength_defence_away"`
	PulseID             int     `json:"pulse_id"`
}

// ==================================

type Element struct {
	ChanceOfPlayingNextRound         int       `json:"chance_of_playing_next_round,omitempty"`
	ChanceOfPlayingThisRound         int       `json:"chance_of_playing_this_round,omitempty"`
	Code                             int       `json:"code,omitempty"`
	CostChangeEvent                  int       `json:"cost_change_event,omitempty"`
	CostChangeEventFall              int       `json:"cost_change_event_fall,omitempty"`
	CostChangeStart                  int       `json:"cost_change_start,omitempty"`
	CostChangeStartFall              int       `json:"cost_change_start_fall,omitempty"`
	DreamteamCount                   int       `json:"dreamteam_count,omitempty"`
	ElementType                      int       `json:"element_type,omitempty"`
	EpNext                           string    `json:"ep_next,omitempty"`
	EpThis                           string    `json:"ep_this,omitempty"`
	EventPoints                      int       `json:"event_points,omitempty"`
	FirstName                        string    `json:"first_name,omitempty"`
	Form                             string    `json:"form,omitempty"`
	ID                               int       `json:"id,omitempty"`
	InDreamteam                      bool      `json:"in_dreamteam,omitempty"`
	News                             string    `json:"news,omitempty"`
	NewsAdded                        time.Time `json:"news_added,omitempty"`
	NowCost                          int       `json:"now_cost,omitempty"`
	Photo                            string    `json:"photo,omitempty"`
	PointsPerGame                    string    `json:"points_per_game,omitempty"`
	SecondName                       string    `json:"second_name,omitempty"`
	SelectedByPercent                string    `json:"selected_by_percent,omitempty"`
	Special                          bool      `json:"special,omitempty"`
	SquadNumber                      *int      `json:"squad_number,omitempty"`
	Status                           string    `json:"status,omitempty"`
	Team                             int       `json:"team,omitempty"`
	TeamCode                         int       `json:"team_code,omitempty"`
	TotalPoints                      int       `json:"total_points,omitempty"`
	TransfersIn                      int       `json:"transfers_in,omitempty"`
	TransfersInEvent                 int       `json:"transfers_in_event,omitempty"`
	TransfersOut                     int       `json:"transfers_out,omitempty"`
	TransfersOutEvent                int       `json:"transfers_out_event,omitempty"`
	ValueForm                        string    `json:"value_form,omitempty"`
	ValueSeason                      string    `json:"value_season,omitempty"`
	WebName                          string    `json:"web_name,omitempty"`
	Region                           *int      `json:"region,omitempty"`
	Minutes                          int       `json:"minutes,omitempty"`
	GoalsScored                      int       `json:"goals_scored,omitempty"`
	Assists                          int       `json:"assists,omitempty"`
	CleanSheets                      int       `json:"clean_sheets,omitempty"`
	GoalsConceded                    int       `json:"goals_conceded,omitempty"`
	OwnGoals                         int       `json:"own_goals,omitempty"`
	PenaltiesSaved                   int       `json:"penalties_saved,omitempty"`
	PenaltiesMissed                  int       `json:"penalties_missed,omitempty"`
	YellowCards                      int       `json:"yellow_cards,omitempty"`
	RedCards                         int       `json:"red_cards,omitempty"`
	Saves                            int       `json:"saves,omitempty"`
	Bonus                            int       `json:"bonus,omitempty"`
	Bps                              int       `json:"bps,omitempty"`
	Influence                        string    `json:"influence,omitempty"`
	Creativity                       string    `json:"creativity,omitempty"`
	Threat                           string    `json:"threat,omitempty"`
	IctIndex                         string    `json:"ict_index,omitempty"`
	Starts                           int       `json:"starts,omitempty"`
	ExpectedGoals                    string    `json:"expected_goals,omitempty"`
	ExpectedAssists                  string    `json:"expected_assists,omitempty"`
	ExpectedGoalInvolvements         string    `json:"expected_goal_involvements,omitempty"`
	ExpectedGoalsConceded            string    `json:"expected_goals_conceded,omitempty"`
	InfluenceRank                    int       `json:"influence_rank,omitempty"`
	InfluenceRankType                int       `json:"influence_rank_type,omitempty"`
	CreativityRank                   int       `json:"creativity_rank,omitempty"`
	CreativityRankType               int       `json:"creativity_rank_type,omitempty"`
	ThreatRank                       int       `json:"threat_rank,omitempty"`
	ThreatRankType                   int       `json:"threat_rank_type,omitempty"`
	IctIndexRank                     int       `json:"ict_index_rank,omitempty"`
	IctIndexRankType                 int       `json:"ict_index_rank_type,omitempty"`
	CornersAndIndirectFreekicksOrder *int      `json:"corners_and_indirect_freekicks_order,omitempty"`
	CornersAndIndirectFreekicksText  string    `json:"corners_and_indirect_freekicks_text,omitempty"`
	DirectFreekicksOrder             *int      `json:"direct_freekicks_order,omitempty"`
	DirectFreekicksText              string    `json:"direct_freekicks_text,omitempty"`
	PenaltiesOrder                   int       `json:"penalties_order,omitempty"`
	PenaltiesText                    string    `json:"penalties_text,omitempty"`
	ExpectedGoalsPer90               float64   `json:"expected_goals_per_90,omitempty"`
	SavesPer90                       float64   `json:"saves_per_90,omitempty"`
	ExpectedAssistsPer90             float64   `json:"expected_assists_per_90,omitempty"`
	ExpectedGoalInvolvementsPer90    float64   `json:"expected_goal_involvements_per_90,omitempty"`
	ExpectedGoalsConcededPer90       float64   `json:"expected_goals_conceded_per_90,omitempty"`
	GoalsConcededPer90               float64   `json:"goals_conceded_per_90,omitempty"`
	NowCostRank                      int       `json:"now_cost_rank,omitempty"`
	NowCostRankType                  int       `json:"now_cost_rank_type,omitempty"`
	FormRank                         int       `json:"form_rank,omitempty"`
	FormRankType                     int       `json:"form_rank_type,omitempty"`
	PointsPerGameRank                int       `json:"points_per_game_rank,omitempty"`
	PointsPerGameRankType            int       `json:"points_per_game_rank_type,omitempty"`
	SelectedRank                     int       `json:"selected_rank,omitempty"`
	SelectedRankType                 int       `json:"selected_rank_type,omitempty"`
	StartsPer90                      float64   `json:"starts_per_90,omitempty"`
	CleanSheetsPer90                 float64   `json:"clean_sheets_per_90,omitempty"`
}

// ==================================

type ElementType struct {
	ID                 int    `json:"id"`
	PluralName         string `json:"plural_name"`
	PluralNameShort    string `json:"plural_name_short"`
	SingularName       string `json:"singular_name"`
	SingularNameShort  string `json:"singular_name_short"`
	SquadSelect        int    `json:"squad_select"`
	SquadMinSelect     *int   `json:"squad_min_select"`
	SquadMaxSelect     *int   `json:"squad_max_select"`
	SquadMinPlay       int    `json:"squad_min_play"`
	SquadMaxPlay       int    `json:"squad_max_play"`
	UIShirtSpecific    bool   `json:"ui_shirt_specific"`
	SubPositionsLocked []int  `json:"sub_positions_locked"`
	ElementCount       int    `json:"element_count"`
}

// ================ DB Structs ================
// ===========================================+

type PlayerSnapshot struct {
	ID                   int     `bson:"id"`
	GameWeek             int     `bson:"game_week"`
	FirstName            string  `bson:"first_name"`
	SecondName           string  `bson:"second_name"`
	WebName              string  `bson:"web_name"`
	Team                 int     `bson:"team"`
	ElementType          int     `bson:"element_type"`
	TotalPoints          int     `bson:"total_points"`
	EventPoints          int     `bson:"event_points"`
	NowCost              int     `bson:"now_cost"`
	Form                 float64 `bson:"form"`
	SelectedByPercent    float64 `bson:"selected_by_percent"`
	Minutes              int     `bson:"minutes"`
	GoalsScored          int     `bson:"goals_scored"`
	Assists              int     `bson:"assists"`
	CleanSheets          int     `bson:"clean_sheets"`
	GoalsConceded        int     `bson:"goals_conceded"`
	OwnGoals             int     `bson:"own_goals"`
	PenaltiesSaved       int     `bson:"penalties_saved"`
	PenaltiesMissed      int     `bson:"penalties_missed"`
	YellowCards          int     `bson:"yellow_cards"`
	RedCards             int     `bson:"red_cards"`
	Saves                int     `bson:"saves"`
	Bonus                int     `bson:"bonus"`
	Bps                  int     `bson:"bps"`
	Influence            float64 `bson:"influence"`
	Creativity           float64 `bson:"creativity"`
	Threat               float64 `bson:"threat"`
	IctIndex             float64 `bson:"ict_index"`
	ExpectedGoals        float64 `bson:"expected_goals"`
	ExpectedAssists      float64 `bson:"expected_assists"`
	ExpectedGoalsPer90   float64 `bson:"expected_goals_per_90"`
	SavesPer90           float64 `bson:"saves_per_90"`
	ExpectedAssistsPer90 float64 `bson:"expected_assists_per_90"`
}

// ==================================

type GameWeekData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	GameWeek  int                `bson:"game_week"`
	Season    string             `bson:"season"`
	Timestamp time.Time          `bson:"timestamp"`
	Players   []PlayerSnapshot   `bson:"players"`
}

// ==================================

type PlayerPerformance struct {
	ID                int     `bson:"_id"`
	WebName           string  `bson:"web_name"`
	TotalPoints       int     `bson:"total_points"`
	AvgPoints         float64 `bson:"avg_points"`
	Team              int     `bson:"team"`
	ElementType       int     `bson:"element_type"`
	GoalsScored       int     `bson:"goals_scored"`
	Assists           int     `bson:"assists"`
	CleanSheets       int     `bson:"clean_sheets"`
	GoalsConceded     int     `bson:"goals_conceded"`
	Saves             int     `bson:"saves"`
	Bonus             int     `bson:"bonus"`
	BPS               int     `bson:"bps"`
	Influence         float64 `bson:"influence"`
	Creativity        float64 `bson:"creativity"`
	Threat            float64 `bson:"threat"`
	ICTIndex          float64 `bson:"ict_index"`
	ExpectedGoals     float64 `bson:"expected_goals"`
	ExpectedAssists   float64 `bson:"expected_assists"`
	NowCost           int     `bson:"now_cost"`
	SelectedByPercent float64 `bson:"selected_by_percent"`
	PerformanceScore  float64 `bson:"performance_score"`
	ValueScore        float64 `bson:"value_score"`
}

// ========= RESPONSE STRUCTS =========
// ====================================

type OptimalTeam struct {
	TotalCost   int                 `json:"total_cost"`
	Goalkeepers []PlayerPerformance `json:"goalkeepers"`
	Defenders   []PlayerPerformance `json:"defenders"`
	Midfielders []PlayerPerformance `json:"midfielders"`
	Forwards    []PlayerPerformance `json:"forwards"`
}
