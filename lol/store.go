package lol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	scribble "github.com/nanobox-io/golang-scribble"
)

type lolCache struct {
	db *scribble.Driver
}

var _ lolStorer = &lolCache{}

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

func (db *lolCache) Close() {

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
// 			logger.Println("err: Errored converting to json from string in get games", err)
// 			continue
// 		}
// 		err = db.SaveGame(g, "NA1")
// 		if err != nil {
// 			logger.Println("err: Failed to save", err)
// 		}
//
// 		err = db.db.Delete("games", fmt.Sprintf("%d_%d_NA1", g.GameID, 34178787))
// 		if err != nil {
// 			logger.Println("err: failed to delete old game", g.GameID)
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
			logger.Println("err: Errored converting to json from string in get games", err)
			continue
		}
		// for _, sum := range game.ParticipantIdentities {
		// 	if sum.Player.AccountID == accountID {
		// 		games = append(games, game)
		// 		break
		// 	}
		// }
	}
	return games, nil
}

func (db *lolCache) TransferFromLocalToMongo(collection string) ([]string, error) {
	// ensure there is a collection to read
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to record location!")
	}

	//
	dir := filepath.Clean("./loldata")
	dir = filepath.Join(dir, collection)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var records []string
	mongo, err := NewLolMongo()
	if err != nil {
		return nil, err
	}

	toSave := make(chan Game)
	go func(saver chan Game) {
		var workerChans []chan Game
		for i := 0; i < 100; i++ {
			gameChan := make(chan Game)
			workerChans = append(workerChans, gameChan)
			go func(games chan Game) {
				for game := range games {
					err := mongo.SaveGame(game, game.PlatformID)
					if err != nil {
						logger.Println(err)
					}
				}
			}(gameChan)
		}
		for {
			for _, gameChan := range workerChans {
				game, ok := <-toSave
				if ok {
					gameChan <- game
				} else {
					return
				}
			}
		}
	}(toSave)

	for _, file := range files {
		gameData, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		var game Game
		err = json.Unmarshal(gameData, &game)
		if err != nil {
			logger.Println("err: Errored converting to json from string in get games", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "CurrentID: %d Done: %d\r", game.GameID, len(records))
		toSave <- game
		records = append(records, fmt.Sprint(game.GameID))
	}
	return records, nil
}

func (db *lolCache) StorePlayer(p Player, gotMatches bool) error {
	return errors.New("NOt implemented lolcache StorePlayer")
}

func (db *lolCache) GetPlayersToVisit() ([]Player, error) {
	return nil, errors.New("NOt implemented lolcache GetPlayersToVisit")
}

func (db *lolCache) VisitPlayer(p Player) error {
	return errors.New("NOt implemented lolcache UpdatePlayer")
}
