package lol

import (
	"fmt"
	"os"
)

func (c *Client) GetAllGames(accountID int64, platformID string) ([]Game, error) {
	var games []Game
	var info *GamesInfoWebUiResponse
	var err error
	info, err = c.WebMatchHistory(accountID, platformID, 0)
	if err != nil {
		return nil, err
	}
	for info.Games.GameIndexEnd < info.Games.GameCount-1 {
		games = append(games, info.Games.Games...)
		info, err = c.WebMatchHistory(accountID, platformID, info.Games.GameIndexEnd)
		if err != nil {
			return nil, err
		}
		var player string
		for _, sum := range info.Games.Games[0].ParticipantIdentities {
			if sum.Player.AccountID == accountID {
				player = sum.Player.SummonerName
				break
			}
		}
		fmt.Fprintf(os.Stdout, "Len Games: %d IndexStart: %d IndexEnd: %d GamesCount: %d Player: %s\r", len(games), info.Games.GameIndexBegin, info.Games.GameIndexEnd, info.Games.GameCount, player)
	}
	return games, nil
}
