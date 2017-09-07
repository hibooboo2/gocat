package lol

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	scribble "github.com/nanobox-io/golang-scribble"
)

type lolCache struct {
	db *scribble.Driver
}

func NewLolB() *lolCache {
	d, err := scribble.New("./loldata", nil)
	if err != nil {
		panic(err)
	}
	db := lolCache{d}
	return &db
}

func (db *lolCache) GetGame(gameID, accountID int64, currentPlatformID string) (game Game, err error) {
	err = db.db.Read("games", fmt.Sprintf("%d_%d_%s", gameID, accountID, currentPlatformID), &game)
	return
}

func (db *lolCache) SaveGame(game Game, accountID int64, currentPlatformID string) error {
	return db.db.Write("games", fmt.Sprintf("%d_%d_%s", game.GameID, accountID, currentPlatformID), &game)
}

func (db *lolCache) GetGames(accountID int64, currentPlatformID string) ([]Game, error) {
	gameNames, err := db.db.ReadAll("games")
	if err != nil {
		return nil, err
	}
	var games []Game
	for _, name := range gameNames {
		vals := strings.Split(strings.TrimSuffix(name, ".json"), "_")
		accID, err := strconv.Atoi(vals[1])
		if err != nil {
			log.Println("err: Parsing accountid for game: ", err)
			continue
		}
		gameID, err := strconv.Atoi(vals[0])
		if err != nil {
			log.Println("err: Parsing gameId for game: ", err)
			continue
		}
		if int64(accID) == accountID && currentPlatformID == vals[2] {
			game, err := db.GetGame(int64(gameID), accountID, currentPlatformID)
			if err != nil {
				log.Println("err: Loading game failed for: ", gameID, accountID, currentPlatformID, err)
				continue
			}
			games = append(games, game)
		}
	}
	return games, nil
}
