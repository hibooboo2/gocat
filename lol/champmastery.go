package lol

import "fmt"

// /lol/champion-mastery/v3/champion-masteries/by-summoner/{summonerId}
// /lol/champion-mastery/v3/champion-masteries/by-summoner/{summonerId}/by-champion/{championId}
// /lol/champion-mastery/v3/scores/by-summoner/{summonerId}
type champMastery struct {
	c *client
}

func (cm *champMastery) All(summonerID int64) []champMasteryDTO {
	var a []champMasteryDTO
	err := cm.c.GetObjRiot(fmt.Sprintf("/lol/champion-mastery/v3/champion-masteries/by-summoner/%d", summonerID), &a)
	if err != nil {
		return nil
	}
	return a
}

func (cm *champMastery) Champ(summonerID, champID int64) *champMasteryDTO {
	var m champMasteryDTO
	err := cm.c.GetObjRiot(fmt.Sprintf(`/lol/champion-mastery/v3/champion-masteries/by-summoner/%d/by-champion/%d`, summonerID, champID), &m)
	if err != nil {
		return nil
	}
	return &m
}

func (cm *champMastery) Total(summonerID int64) int {
	var total int
	err := cm.c.GetObjRiot(fmt.Sprintf("/lol/champion-mastery/v3/scores/by-summoner/%d", summonerID), &total)
	if err != nil {
		return 0
	}
	return total
}

type champMasteryDTO struct {
	PlayerID                     int   `json:"playerId"`
	ChampionID                   int   `json:"championId"`
	ChampionLevel                int   `json:"championLevel"`
	ChampionPoints               int   `json:"championPoints"`
	LastPlayTime                 int64 `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int   `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int   `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool  `json:"chestGranted"`
	TokensEarned                 int   `json:"tokensEarned"`
}
