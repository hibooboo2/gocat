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
	for info.Games.GameIndexEnd < info.Games.GameCount {
		games = append(games, info.Games.Games...)
		info, err = c.WebMatchHistory(34178787, platformID, info.Games.GameIndexEnd+1)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stdout, "Len Games: %d IndexStart: %d IndexEnd: %d GamesCount: %d\r", len(games), info.Games.GameIndexBegin, info.Games.GameIndexEnd, info.Games.GameCount)
	}
	return games, nil
}
