package lol

import (
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
	players        *mgo.Collection
	playersVisited *mgo.Collection
	lolCache
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
	logger.Printf("debug: games: %d", n)
	n, _ = session.DB("lol").C("players").Count()
	logger.Printf("debug: left to visit: %d", n)
	n, _ = session.DB("lol").C("playersvisited").Count()
	logger.Printf("debug: visited: %d", n)
	mongo := &lolMongo{
		session:        session,
		db:             session.DB("lol"),
		games:          session.DB("lol").C("games"),
		players:        session.DB("lol").C("players"),
		playersVisited: session.DB("lol").C("playersvisited"),
	}
	mongo.lolCache = &memCache{}
	return mongo, nil
}

func NewLolMongoWAccess(host string, port int) (*lolMongo, error) {
	db, err := NewLolMongo(host, port)
	if err != nil {
		return nil, err
	}
	return db.(*lolMongo), err
}

func (db *lolMongo) GetGame(gameID int64, currentPlatformID string) (Game, error) {
	var game Game
	// err := db.db.Read("games", fmt.Sprintf("%d_%s", gameID, currentPlatformID), &game)
	err := db.games.Find(bson.M{"gameid": gameID, "platformid": currentPlatformID}).One(&game)
	db.lolCache.AddGame(gameID)
	return game, err
}

func (db *lolMongo) CheckGameStored(gameID int64) bool {
	if db.lolCache.HaveGame(gameID) {
		return true
	}
	n, err := db.games.Find(bson.M{"gameid": gameID}).Count()
	have := err == nil && n > 0
	if have {
		db.lolCache.AddGame(gameID)
	}
	return have
}

func (db *lolMongo) SaveGame(game Game, currentPlatformID string) error {
	if db.lolCache.HaveGame(game.GameID) {
		return nil
	}
	n, _ := db.games.Find(bson.M{"gameid": game.GameID}).Count()
	if n == 0 {
		err := db.games.Insert(&game)
		if err != nil {
			return err
		}
		db.lolCache.AddGame(game.GameID)
		return nil
	}
	return nil
}

func (db *lolMongo) Close() {
	db.session.Close()
}

func (db *lolMongo) StorePlayer(p Player) error {
	if db.lolCache.HaveVisitedPlayer(p.AccountID) {
		return nil
	}
	db.lolCache.Player(p.AccountID)
	n, _ := db.playersVisited.Find(bson.M{"accountid": p.AccountID}).Count()
	if n == 0 {
		return db.players.Insert(&p)
	}
	db.lolCache.VisitedPlayer(p.AccountID)
	return nil
}

func (db *lolMongo) GetPlayerToVisit() int64 {
	id := db.lolCache.GetPlayerToVisit()
	if id != 0 {
		db.lolCache.VisitedPlayer(id)
		db.players.Remove(bson.M{"accountid": id})
		db.playersVisited.Insert(bson.M{"accountid": id})
		return id
	}
	var p Player
	err := db.players.Find(bson.M{"platformid": NA1}).One(&p)
	if err != nil {
		logger.Println("err: Couldnt find a player:", err)
		return 0
	}
	db.VisitPlayer(p)
	return p.AccountID
}

func (db *lolMongo) VisitPlayer(p Player) error {
	if db.lolCache.HaveVisitedPlayer(p.AccountID) {
		return nil
	}
	db.lolCache.VisitedPlayer(p.AccountID)
	err := db.players.Remove(bson.M{"accountid": p.AccountID})
	if err != nil {
		logger.Println("err: While trying to remove player: ", err)
	}
	return db.playersVisited.Insert(&p)
}

func (db *lolMongo) Stats() {
	var diffs []int
	prevCount, _ := db.games.Count()
	for {
		time.Sleep(time.Second)
		g, _ := db.games.Count()
		diff := g - prevCount
		diffs = append(diffs, diff)
		rate := avg(diffs)
		if len(diffs) > 60 {
			diffs = diffs[:30]
		}
		prevCount = g
		p, _ := db.players.Count()
		pv, _ := db.playersVisited.Count()
		fmt.Fprintf(os.Stdout, "GameAddRate %0f/s\t Games: %d\t Players %d\t PlayersVisited %d\r", rate, g, p, pv)
	}
}

func (db *lolMongo) LoadAllGameIDS() {
	start := time.Now()
	var ids []Game
	var err error
	var id int64
	n := 0
	batchSize := 2000
	for err == nil {
		err = db.games.Find(nil).Select(bson.M{"gameid": 1}).Limit(batchSize).Skip(n * batchSize).All(&ids)
		var gameID int64
		for _, v := range ids {
			logger.Println(v.GameID)
			db.lolCache.AddGame(v.GameID)
			gameID = v.GameID
		}
		if id == gameID {
			break
		}
		id = gameID
		n++
		fmt.Fprintf(os.Stdout, "Game IDS found: %d\r", n*batchSize)
	}
	if err != nil {
		logger.Println("err: After gamesid load:", err)
	}
	count := 0
	db.lolCache.(*memCache).games.Range(func(key interface{}, value interface{}) bool {
		count++
		logger.Println(key)
		return true
	})
	logger.Println("info: Found ", count, " games")
	var playersVisited []Player
	n = 0
	for err == nil {
		err = db.playersVisited.Find(nil).Limit(batchSize).Skip(n * batchSize).All(&playersVisited)
		for _, p := range playersVisited {
			db.lolCache.VisitedPlayer(p.AccountID)
		}
		n++
		fmt.Fprintf(os.Stdout, "PlayersVisited found: %d\r", n*batchSize)
		if len(playersVisited) == 0 {
			break
		}
		// pid = players[0].AccountID
		fmt.Fprintf(os.Stdout, "PlayersVisited found: %d \tPlayerID: %d\t PlayerName: %s\r", n*batchSize, playersVisited[0].AccountID, playersVisited[0].SummonerName)
	}
	if err != nil {
		logger.Println("err: After playersvisited load:", err)
	}
	n = 0
	var players []Player
	for err == nil {
		err = db.players.Find(nil).Limit(batchSize).Skip(n * batchSize).All(&players)
		for _, p := range players {
			db.lolCache.Player(p.AccountID)
		}
		n++
		if len(players) == 0 {
			break
		}
		fmt.Fprintf(os.Stdout, "PlayersToVisit found: %d PlayerName: %s PlayerID: %d\r", n*batchSize, players[0].SummonerName, players[0].AccountID)
		// pid = players[0].AccountID
	}
	if err != nil {
		logger.Println("err: After playerstovisit load:", err)
	}
	playersFound := 0
	counter := make(map[interface{}]interface{})
	db.lolCache.(*memCache).visited.Range(func(key interface{}, value interface{}) bool {
		playersFound++
		counter[key] = value
		return true
	})
	db.lolCache.(*memCache).toVisit.Range(func(key interface{}, value interface{}) bool {
		playersFound++
		counter[key] = value
		return true
	})
	logger.Println("info: Found:", playersFound, " players", len(counter))
	logger.Println("info: Took: ", time.Since(start), " to load all playerids and gameids")
}

func avg(vals []int) float64 {
	var avg float64
	for _, val := range vals {
		avg += float64(val)
	}
	avg = avg / float64(len(vals))
	return avg
}
