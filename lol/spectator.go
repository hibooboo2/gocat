package lol

import "fmt"

type spectator struct {
	// /lol/spectator/v3/active-games/by-summoner/{summonerId}
	// /lol/spectator/v3/featured-games
	c *client
}

func (s *spectator) Game(summonerID int64) *CurrentGameInfo {
	var g CurrentGameInfo
	err := s.c.GetObjRiot(fmt.Sprintf("/lol/spectator/v3/active-games/by-summoner/%d", summonerID), &g)
	if err != nil {
		return nil
	}
	return &g
}

func (s *spectator) GameSummonerName(summonerName string) *CurrentGameInfo {
	sum := s.c.Summoners().ByName(summonerName)
	g := s.Game(sum.ID)
	logger.Println("trace: game: ", g)
	return g
}

func (s *spectator) Featured() []CurrentGameInfo {
	var games []CurrentGameInfo
	err := s.c.GetObjRiot("/lol/spectator/v3/featured-games", &games)
	if err != nil {
		return nil
	}
	return games
}

type CurrentGameInfo struct {
	GameID            int64  `json:"gameId"`
	GameStartTime     int64  `json:"gameStartTime"`
	PlatformID        string `json:"platformId"`
	GameMode          string `json:"gameMode"`
	MapID             int    `json:"mapId"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int    `json:"gameQueueConfigId"`
	Observers         struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	Participants []struct {
		ProfileIconID int    `json:"profileIconId"`
		ChampionID    int    `json:"championId"`
		SummonerName  string `json:"summonerName"`
		Runes         []struct {
			Count  int `json:"count"`
			RuneID int `json:"runeId"`
		} `json:"runes"`
		Bot       bool `json:"bot"`
		Masteries []struct {
			MasteryID int `json:"masteryId"`
			Rank      int `json:"rank"`
		} `json:"masteries"`
		Spell2ID   int   `json:"spell2Id"`
		TeamID     int   `json:"teamId"`
		Spell1ID   int   `json:"spell1Id"`
		SummonerID int64 `json:"summonerId"`
	} `json:"participants"`
	GameLength      int `json:"gameLength"`
	BannedChampions []struct {
		TeamID     int `json:"teamId"`
		ChampionID int `json:"championId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
}
