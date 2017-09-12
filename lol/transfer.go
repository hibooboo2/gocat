package lol

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *lolMongo) TransferToAnother(host string, port int) error {
	db2, err := NewLolMongoWAccess(host, port)
	if err != nil {
		return err
	}
	logger.Println("Starting transfer")
	batchSize := 100
	totalGames, _ := db.games.Find(nil).Count()
	var count int
	for count < totalGames {
		var games []Game
		err = db.games.Find(nil).Skip(count).Limit(batchSize).All(&games)
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		count += len(games)
		logger.Println("Got batch", count)
		b := db2.games.Bulk()
		for _, game := range games {
			b.Insert(game)
		}
		res, err := b.Run()
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		logger.Println("Inserted batch", count)
		if res.Matched+res.Modified != len(games) {
			logger.Printf("May have skipped games. Games: %d Matched+Mod: %d", len(games), res.Matched+res.Modified)
		}
		games = nil
	}
	logger.Println("Moved ", count, "Games")

	var players []Player
	totalPlayers, _ := db.players.Find(nil).Count()
	count = 0
	for count < totalPlayers {
		err = db.players.Find(nil).Skip(count).Batch(batchSize).All(&players)
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		count += len(players)
		b := db2.players.Bulk()
		for _, player := range players {
			b.Insert(player)
		}
		res, err := b.Run()
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		if res.Matched+res.Modified != len(players) {
			logger.Printf("May have skipped players. Players: %d Matched+Mod: %d", len(players), res.Matched+res.Modified)
		}
		players = nil
	}

	totalPlayers, _ = db.playersVisited.Find(nil).Count()
	count = 0
	for count < totalPlayers {
		err = db.playersVisited.Find(nil).Skip(count).Batch(batchSize).All(&players)
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		count += len(players)
		b := db2.playersVisited.Bulk()
		for _, player := range players {
			b.Insert(player)
		}
		res, err := b.Run()
		if err != nil {
			logger.Println("err:", err)
			return err
		}
		if res.Matched+res.Modified != len(players) {
			logger.Printf("May have skipped players. Players: %d Matched+Mod: %d", len(players), res.Matched+res.Modified)
		}
		players = nil
	}

	return nil
}

func (db *lolMongo) GameIDSToIDTable() {
	// log.Println(db.session.DB("lol").C("gamesid").DropCollection())
	// start := time.Now()
	// log.Println("Building index")
	// db.games.Pipe()
	// log.Println(db.games.EnsureIndex(mgo.Index{
	// 	Key:      []string{"gameid"},
	// 	DropDups: true,
	// 	Unique:   false,
	// 	Name:     "game_IDs_index",
	// }))
	// log.Println("Index Built Took:", time.Since(start))
	logger.Println(db.session.BuildInfo())
	// logger.Println(testDB.Pipe([]bson.M{{"$group": {"_id": "$gameid", "dups": {"$push": "$_id"}}}}).All(&stuff))
	var ids []int64
	logger.Println(db.games.Find(nil).Distinct("gameid", &ids))
	logger.Println(len(ids))
	// var games []Game
	// logger.Println(db.games.Find(nil).Select(bson.M{"gameid": 1}).All(&games))
	// logger.Println(len(games))
	for _, gameID := range ids {
		var game Game
		game.GameID = gameID
		var gameCopy Game
		n, err := db.games.Find(bson.M{"gameid": game.GameID}).Count()
		if n == 0 && err == nil {
			fmt.Fprintf(os.Stdout, "Skipped: N: %d  ID: %d\r", n, game.GameID)
			continue
		} else if n > 0 && err == nil {
			fmt.Fprintf(os.Stdout, "Found: N: %d ID: %d\r", n, game.GameID)
			continue
		}
		if err != nil {
			logger.Println("err: error on count", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "Working: %d\r", game.GameID)
		err = db.games.Find(bson.M{"gameid": game.GameID}).One(&gameCopy)
		if err != nil {
			logger.Println("err: Errored on game query", game.GameID, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "Working copy: %d\r", gameCopy.GameID)
		n, err = db.games.Find(bson.M{"gameid": game.GameID}).Count()
		if err != nil {
			logger.Println("err: May have lost game:", gameCopy.GameID, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "n: %d \tWorking copy: %d\r", n, gameCopy.GameID)
		for n > 1 {
			err = db.games.Remove(bson.M{"gameid": game.GameID})
			if err != nil {
				logger.Println("err: errored on game removal", game.GameID)
				break
			}
			logger.Println("Removed: ", game.GameID)
			n, err = db.games.Find(bson.M{"gameid": game.GameID}).Count()
			if err != nil {
				logger.Println("err: errored on game find", game.GameID)
				break
			}
			fmt.Fprintf(os.Stdout, "n: %d \tWorking copy: %d\r", n, gameCopy.GameID)
		}
		n, err = db.games.Find(bson.M{"gameid": game.GameID}).Count()
		if n == 0 && err == nil {
			err = db.games.Insert(gameCopy)
			if err != nil {
				logger.Println("err: May have lost game:", gameCopy.GameID, err)
				continue
			}
		}
		if err != nil {
			logger.Println("err: May have lost game:", gameCopy.GameID, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "Removed Dups for: %d First sum: %s\r", gameCopy.GameID, gameCopy.ParticipantIdentities[0].Player.SummonerName)
	}
	gamesCount, _ := db.games.Count()
	logger.Println("After games de dup loop:", gamesCount, "Games in map:", len(ids))
	logger.Println(db.games.EnsureIndex(mgo.Index{
		Key:      []string{"gameid"},
		DropDups: true,
		Unique:   true,
	}))

}

func (db *lolMongo) GetGameRan() {
	// var games []Game
	// logger.Println(db.games.Find(bson.M{"gameid": 2591856267}).Select(bson.M{"gameid": 1, "participantidentities": 1}).All(&games))
	// logger.Println(games)
	// logger.Println(len(games))
	logger.Println(db.games.Count())
}
