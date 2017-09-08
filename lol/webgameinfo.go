package lol

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Team struct {
	TeamID               int    `json:"teamId"`
	Win                  string `json:"win"`
	FirstBlood           bool   `json:"firstBlood"`
	FirstTower           bool   `json:"firstTower"`
	FirstInhibitor       bool   `json:"firstInhibitor"`
	FirstBaron           bool   `json:"firstBaron"`
	FirstDragon          bool   `json:"firstDragon"`
	FirstRiftHerald      bool   `json:"firstRiftHerald"`
	TowerKills           int    `json:"towerKills"`
	InhibitorKills       int    `json:"inhibitorKills"`
	BaronKills           int    `json:"baronKills"`
	DragonKills          int    `json:"dragonKills"`
	VilemawKills         int    `json:"vilemawKills"`
	RiftHeraldKills      int    `json:"riftHeraldKills"`
	DominionVictoryScore int    `json:"dominionVictoryScore"`
	Bans                 []Ban  `json:"bans"`
}

type Ban struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type Game struct {
	GameID                int64                 `json:"gameId"`
	PlatformID            string                `json:"platformId"`
	GameCreation          int64                 `json:"gameCreation"`
	GameDuration          int                   `json:"gameDuration"`
	QueueID               int                   `json:"queueId"`
	MapID                 int                   `json:"mapId"`
	SeasonID              int                   `json:"seasonId"`
	GameVersion           string                `json:"gameVersion"`
	GameMode              string                `json:"gameMode"`
	GameType              string                `json:"gameType"`
	Teams                 []Team                `json:"teams"`
	Participants          []Participant         `json:"participants"`
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
}
type ParticipantIdentity struct {
	ParticipantID int    `json:"participantId"`
	Player        Player `json:"player"`
}

// WebMatchHistory circumvent riots api throttling. Or at least attepmt to. This is using the endpoint that the web ui uses. No docs for it.
func (c *Client) WebMatch(gameID int64, currentPlatformID string) (*Game, error) {
	game, err := c.cache.GetGame(gameID, currentPlatformID)
	if err == nil {
		return &game, nil
	}
	query := make(url.Values)
	// query.Add("visibleAccountId", fmt.Sprintf(`%d`, accountID))
	query.Add("visiblePlatformId", currentPlatformID)
	//https://acs.leagueoflegends.com/v1/stats/game/NA1/2591856267?visiblePlatformId=NA1&visibleAccountId=237823602

	reqURL, err := url.Parse(fmt.Sprintf("https://acs.leagueoflegends.com/v1/stats/game/%s/%d", currentPlatformID, gameID))
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
	err = json.NewDecoder(resp.Body).Decode(&game)
	if err != nil {
		return nil, err
	}
	err = c.cache.SaveGame(game, currentPlatformID)
	if err != nil {
		// logger.Println("err: Failed to save game to db / cache", err)
	}
	return &game, nil
}
