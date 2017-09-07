package lol

import (
	"encoding/json"
	"fmt"
	"log"

	scribble "github.com/nanobox-io/golang-scribble"
)

type lolCache struct {
	db *scribble.Driver
}

func NewLolCache() *lolCache {
	d, err := scribble.New("./loldata", nil)
	if err != nil {
		panic(err)
	}
	db := lolCache{d}
	return &db
}

func (db *lolCache) GetGame(gameID int64, currentPlatformID string) (Game, error) {
	var game Game
	err := db.db.Read("games", fmt.Sprintf("%d_%s", gameID, currentPlatformID), &game)
	return game, err
}

//
// func (db *lolCache) UpdateToNewName() error {
// 	games, err := db.db.ReadAll("games")
// 	if err != nil {
// 		return err
// 	}
// 	for _, gameData := range games {
// 		var g Game
// 		err := json.Unmarshal([]byte(gameData), &g)
// 		if err != nil {
// 			log.Println("err: Errored converting to json from string in get games", err)
// 			continue
// 		}
// 		err = db.SaveGame(g, "NA1")
// 		if err != nil {
// 			log.Println("err: Failed to save", err)
// 		}
//
// 		err = db.db.Delete("games", fmt.Sprintf("%d_%d_NA1", g.GameID, 34178787))
// 		if err != nil {
// 			log.Println("err: failed to delete old game", g.GameID)
// 			continue
// 		}
// 	}
// 	return nil
// }

func (db *lolCache) SaveGame(game Game, currentPlatformID string) error {
	return db.db.Write("games", fmt.Sprintf("%d_%s", game.GameID, currentPlatformID), &game)
}

func (db *lolCache) GetGames(accountID int64, currentPlatformID string) ([]Game, error) {
	gameNames, err := db.db.ReadAll("games")
	if err != nil {
		return nil, err
	}
	var games []Game
	for _, gameData := range gameNames {
		var game Game
		err := json.Unmarshal([]byte(gameData), &game)
		if err != nil {
			log.Println("err: Errored converting to json from string in get games", err)
			continue
		}
		for _, sum := range game.ParticipantIdentities {
			if sum.Player.AccountID == accountID {
				games = append(games, game)
				break
			}
		}
	}
	return games, nil
}
