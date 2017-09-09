package lol

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type lolMongo struct {
	session        *mgo.Session
	db             *mgo.Database
	games          *mgo.Collection
	gamesid        *mgo.Collection
	players        *mgo.Collection
	playersVisited *mgo.Collection
}

var _ lolStorer = &lolMongo{}

func NewLolMongo(host string, port int) (lolStorer, error) {
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 27017
	}

	session, err := mgo.Dial(fmt.Sprintf(`%s:%d`, host, port))
	if err != nil {
		return nil, err
	}
	n, _ := session.DB("lol").C("games").Count()
	logger.Printf("debug:games: %d", n)
	n, _ = session.DB("lol").C("players").Count()
	logger.Printf("debug:left to visit: %d", n)
	n, _ = session.DB("lol").C("playersvisited").Count()
	logger.Printf("debug: visited: %d", n)
	return &lolMongo{
		session,
		session.DB("lol"),
		session.DB("lol").C("games"),
		session.DB("lol").C("gamesid"),
		session.DB("lol").C("players"),
		session.DB("lol").C("playersvisited"),
	}, nil
}

func NewLolMongoWAccess(host string, port int) (*lolMongo, error) {
	db, err := NewLolMongo(host, port)
	return db.(*lolMongo), err
}

func (db *lolMongo) GetGame(gameID int64, currentPlatformID string) (Game, error) {
	var game Game
	// err := db.db.Read("games", fmt.Sprintf("%d_%s", gameID, currentPlatformID), &game)
	n, _ := db.gamesid.Find(bson.M{"gameid": gameID}).Count()
	if n == 0 {
		return game, errors.New("Game Not Found in DB")
	}
	err := db.games.Find(bson.M{"gameid": gameID, "platformid": currentPlatformID}).One(&game)
	return game, err
}

func (db *lolMongo) CheckGameStored(gameID int64) bool {
	n, err := db.gamesid.Find(bson.M{"gameid": gameID}).Count()
	return err == nil && n > 0
}

func (db *lolMongo) SaveGame(game Game, currentPlatformID string) error {
	// return db.db.Write("games", fmt.Sprintf("%d_%s", game.GameID, currentPlatformID), &game)
	n, _ := db.games.Find(bson.M{"gameid": game.GameID, "platformid": currentPlatformID}).Count()
	if n == 0 {
		_, err := db.games.Upsert(bson.M{"gameid": game.GameID}, bson.M{"gameid": game.GameID})
		if err != nil {
			logger.Println("alert: Cailed to upsert id to cache table", err)
		}
		return db.games.Insert(&game)
	}
	return nil
}

func (db *lolMongo) GetGames(accountID int64, currentPlatformID string) ([]Game, error) {
	return nil, errors.New("Not implemented: get games mongo store")
}

func (db *lolMongo) Close() {
	db.session.Close()
}

func (db *lolMongo) StorePlayer(p Player, gotMatches bool) error {
	n, _ := db.playersVisited.Find(bson.M{"accountid": p.AccountID}).Count()
	if n == 0 {
		return db.players.Insert(&p)
	}
	return nil
}
func (db *lolMongo) GetPlayersToVisit() ([]Player, error) {
	var players []Player
	err := db.players.Find(bson.M{"platformid": "NA1"}).Limit(1000).All(&players)
	return players, err
}

func (db *lolMongo) GetPlayerToVisit() (Player, error) {
	var p Player
	err := db.players.Find(bson.M{"platformid": "NA1"}).One(&p)
	db.VisitPlayer(p)
	return p, err
}

func (db *lolMongo) VisitPlayer(p Player) error {
	err := db.players.Remove(bson.M{"accountid": p.AccountID})
	if err != nil {
		return err
	}
	return db.playersVisited.Insert(&p)
}

func (db *lolMongo) Stats() {
	var diffs []int
	prevCount, _ := db.games.Count()
	for {
		time.Sleep(time.Millisecond * 500)
		g, _ := db.games.Count()
		diff := g - prevCount
		diffs = append(diffs, diff)
		rate := avg(diffs) * 2
		if len(diffs) > 60 {
			diffs = diffs[:30]
		}
		prevCount = g
		p, _ := db.players.Count()
		pv, _ := db.playersVisited.Count()
		fmt.Fprintf(os.Stdout, "GameAddRate %0f/s\t Games: %d\t Players %d\t PlayersVisited %d\r", rate, g, p, pv)
	}
}

func avg(vals []int) float64 {
	var avg float64
	for _, val := range vals {
		avg += float64(val)
	}
	avg = avg / float64(len(vals))
	return avg
}

func (db *lolMongo) TransferToAnother(host string, port int) error {
	// db.games.Find(bson.M{"platformid": "NA1"}).Limit(limit).Iter().
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
	batchSize := 100
	totalGames, _ := db.games.Find(nil).Count()
	var count int
	for count < totalGames {
		var games []Game
		err := db.games.Find(nil).Skip(count).Limit(batchSize).All(&games)
		if err != nil {
			logger.Println("err:", err)
			return
		}
		count += len(games)
		logger.Println("Got batch", count)
		b := db.gamesid.Bulk()
		for _, game := range games {
			b.Insert(bson.M{"gameid": game.GameID})
		}
		res, err := b.Run()
		if err != nil {
			logger.Println("err:", err)
			return
		}
		logger.Println("Inserted batch", count)
		if res.Matched+res.Modified != len(games) {
			logger.Printf("May have skipped games. Games: %d Matched+Mod: %d", len(games), res.Matched+res.Modified)
		}
		games = nil
	}
	logger.Println(db.gamesid.Find(nil).Count())
}
