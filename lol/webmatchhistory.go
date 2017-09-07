package lol

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type GamesInfoWebUiResponse struct {
	PlatformID  string         `json:"platformId"`
	AccountID   int64          `json:"accountId"`
	Games       GamesInfoWebUi `json:"games"`
	ShownQueues []interface{}  `json:"shownQueues"`
}

type GamesInfoWebUi struct {
	GameIndexBegin     int    `json:"gameIndexBegin"`
	GameIndexEnd       int    `json:"gameIndexEnd"`
	GameTimestampBegin int64  `json:"gameTimestampBegin"`
	GameTimestampEnd   int64  `json:"gameTimestampEnd"`
	GameCount          int    `json:"gameCount"`
	Games              []Game `json:"games"`
}

type Player struct {
	PlatformID        string `json:"platformId"`
	AccountID         int64  `json:"accountId"`
	SummonerName      string `json:"summonerName"`
	SummonerID        int64  `json:"summonerId"`
	CurrentPlatformID string `json:"currentPlatformId"`
	CurrentAccountID  int64  `json:"currentAccountId"`
	MatchHistoryURI   string `json:"matchHistoryUri"`
	ProfileIcon       int    `json:"profileIcon"`
}

type PlayerMatchStats struct {
	ParticipantID                   int  `json:"participantId"`
	Win                             bool `json:"win"`
	Item0                           int  `json:"item0"`
	Item1                           int  `json:"item1"`
	Item2                           int  `json:"item2"`
	Item3                           int  `json:"item3"`
	Item4                           int  `json:"item4"`
	Item5                           int  `json:"item5"`
	Item6                           int  `json:"item6"`
	Kills                           int  `json:"kills"`
	Deaths                          int  `json:"deaths"`
	Assists                         int  `json:"assists"`
	LargestKillingSpree             int  `json:"largestKillingSpree"`
	LargestMultiKill                int  `json:"largestMultiKill"`
	KillingSprees                   int  `json:"killingSprees"`
	LongestTimeSpentLiving          int  `json:"longestTimeSpentLiving"`
	DoubleKills                     int  `json:"doubleKills"`
	TripleKills                     int  `json:"tripleKills"`
	QuadraKills                     int  `json:"quadraKills"`
	PentaKills                      int  `json:"pentaKills"`
	UnrealKills                     int  `json:"unrealKills"`
	TotalDamageDealt                int  `json:"totalDamageDealt"`
	MagicDamageDealt                int  `json:"magicDamageDealt"`
	PhysicalDamageDealt             int  `json:"physicalDamageDealt"`
	TrueDamageDealt                 int  `json:"trueDamageDealt"`
	LargestCriticalStrike           int  `json:"largestCriticalStrike"`
	TotalDamageDealtToChampions     int  `json:"totalDamageDealtToChampions"`
	MagicDamageDealtToChampions     int  `json:"magicDamageDealtToChampions"`
	PhysicalDamageDealtToChampions  int  `json:"physicalDamageDealtToChampions"`
	TrueDamageDealtToChampions      int  `json:"trueDamageDealtToChampions"`
	TotalHeal                       int  `json:"totalHeal"`
	TotalUnitsHealed                int  `json:"totalUnitsHealed"`
	DamageSelfMitigated             int  `json:"damageSelfMitigated"`
	DamageDealtToObjectives         int  `json:"damageDealtToObjectives"`
	DamageDealtToTurrets            int  `json:"damageDealtToTurrets"`
	VisionScore                     int  `json:"visionScore"`
	TimeCCingOthers                 int  `json:"timeCCingOthers"`
	TotalDamageTaken                int  `json:"totalDamageTaken"`
	MagicalDamageTaken              int  `json:"magicalDamageTaken"`
	PhysicalDamageTaken             int  `json:"physicalDamageTaken"`
	TrueDamageTaken                 int  `json:"trueDamageTaken"`
	GoldEarned                      int  `json:"goldEarned"`
	GoldSpent                       int  `json:"goldSpent"`
	TurretKills                     int  `json:"turretKills"`
	InhibitorKills                  int  `json:"inhibitorKills"`
	TotalMinionsKilled              int  `json:"totalMinionsKilled"`
	NeutralMinionsKilled            int  `json:"neutralMinionsKilled"`
	NeutralMinionsKilledTeamJungle  int  `json:"neutralMinionsKilledTeamJungle"`
	NeutralMinionsKilledEnemyJungle int  `json:"neutralMinionsKilledEnemyJungle"`
	TotalTimeCrowdControlDealt      int  `json:"totalTimeCrowdControlDealt"`
	ChampLevel                      int  `json:"champLevel"`
	VisionWardsBoughtInGame         int  `json:"visionWardsBoughtInGame"`
	SightWardsBoughtInGame          int  `json:"sightWardsBoughtInGame"`
	WardsPlaced                     int  `json:"wardsPlaced"`
	WardsKilled                     int  `json:"wardsKilled"`
	FirstBloodKill                  bool `json:"firstBloodKill"`
	FirstBloodAssist                bool `json:"firstBloodAssist"`
	FirstTowerKill                  bool `json:"firstTowerKill"`
	FirstTowerAssist                bool `json:"firstTowerAssist"`
	FirstInhibitorKill              bool `json:"firstInhibitorKill"`
	FirstInhibitorAssist            bool `json:"firstInhibitorAssist"`
	CombatPlayerScore               int  `json:"combatPlayerScore"`
	ObjectivePlayerScore            int  `json:"objectivePlayerScore"`
	TotalPlayerScore                int  `json:"totalPlayerScore"`
	TotalScoreRank                  int  `json:"totalScoreRank"`
	PlayerScore0                    int  `json:"playerScore0"`
	PlayerScore1                    int  `json:"playerScore1"`
	PlayerScore2                    int  `json:"playerScore2"`
	PlayerScore3                    int  `json:"playerScore3"`
	PlayerScore4                    int  `json:"playerScore4"`
	PlayerScore5                    int  `json:"playerScore5"`
	PlayerScore6                    int  `json:"playerScore6"`
	PlayerScore7                    int  `json:"playerScore7"`
	PlayerScore8                    int  `json:"playerScore8"`
	PlayerScore9                    int  `json:"playerScore9"`
}

type Participant struct {
	ParticipantID int `json:"participantId"`
	TeamID        int `json:"teamId"`
	ChampionID    int `json:"championId"`
	Spell1ID      int `json:"spell1Id"`
	Spell2ID      int `json:"spell2Id"`
	Masteries     []struct {
		MasteryID int `json:"masteryId"`
		Rank      int `json:"rank"`
	} `json:"masteries"`
	Runes []struct {
		RuneID int `json:"runeId"`
		Rank   int `json:"rank"`
	} `json:"runes"`
	HighestAchievedSeasonTier string           `json:"highestAchievedSeasonTier"`
	Stats                     PlayerMatchStats `json:"stats"`
	Timeline                  TimeLine         `json:"timeline"`
}

type TimeLine struct {
	ParticipantID      int `json:"participantId"`
	CreepsPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"creepsPerMinDeltas"`
	XpPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"xpPerMinDeltas"`
	GoldPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"goldPerMinDeltas"`
	CsDiffPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"csDiffPerMinDeltas"`
	XpDiffPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"xpDiffPerMinDeltas"`
	DamageTakenPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"damageTakenPerMinDeltas"`
	DamageTakenDiffPerMinDeltas struct {
		One020    float64 `json:"10-20"`
		Zero10    float64 `json:"0-10"`
		Three0End float64 `json:"30-end"`
		Two030    float64 `json:"20-30"`
	} `json:"damageTakenDiffPerMinDeltas"`
	Role string `json:"role"`
	Lane string `json:"lane"`
}

// WebMatchHistory circumvent riots api throttling. Or at least attepmt to. This is using the endpoint that the web ui uses. No docs for it.
func (c *Client) WebMatchHistory(accountID int64, platformID string, index int) (*GamesInfoWebUiResponse, error) {
	var games GamesInfoWebUiResponse
	query := make(url.Values)
	if index != 0 {
		query.Add("begIndex", fmt.Sprintf(`%d`, index))
		query.Add("endIndex", fmt.Sprintf(`%d`, index+20))
	}
	reqURL, err := url.Parse(fmt.Sprintf("https://acs.leagueoflegends.com/v1/stats/player_history/%s/%d", platformID, accountID))
	reqURL.RawQuery = query.Encode()
	if err != nil {
		return nil, err
	}
	resp, err := c.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Expected 200 code got: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, err
	}
	// https://acs.leagueoflegends.com/v1/stats/player_history/NA/34178787?begIndex=200`, a ...interface{}))
	return &games, nil
}
