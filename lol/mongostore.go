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
	logger.Printf("debug:games: %d", n)
	n, _ = session.DB("lol").C("players").Count()
	logger.Printf("debug:left to visit: %d", n)
	n, _ = session.DB("lol").C("playersvisited").Count()
	logger.Printf("debug: visited: %d", n)
	mongo := &lolMongo{
		session:        session,
		db:             session.DB("lol"),
		games:          session.DB("lol").C("games"),
		gamesid:        session.DB("lol").C("gamesid"),
		players:        session.DB("lol").C("players"),
		playersVisited: session.DB("lol").C("playersvisited"),
	}
	mongo.lolCache = &memCache{}
	return mongo, nil
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
	db.lolCache.AddGame(gameID)
	return game, err
}

func (db *lolMongo) CheckGameStored(gameID int64) bool {
	if db.lolCache.HaveGame(gameID) {
		return true
	}
	n, err := db.gamesid.Find(bson.M{"gameid": gameID}).Count()
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
	n, _ := db.gamesid.Find(bson.M{"gameid": game.GameID}).Count()
	if n == 0 {
		_, err := db.gamesid.Upsert(bson.M{"gameid": game.GameID}, bson.M{"gameid": game.GameID})
		if err != nil {
			return err
		}

		err = db.games.Insert(&game)
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
		gid, _ := db.gamesid.Count()
		diff := g - prevCount
		diffs = append(diffs, diff)
		rate := avg(diffs)
		if len(diffs) > 60 {
			diffs = diffs[:30]
		}
		prevCount = g
		p, _ := db.players.Count()
		pv, _ := db.playersVisited.Count()
		fmt.Fprintf(os.Stdout, "GameAddRate %0f/s\t Games: %d\t GameIDs: %d\t Players %d\t PlayersVisited %d\r", rate, g, gid, p, pv)
	}
}

func (db *lolMongo) LoadAllGameIDS() {
	var ids []bson.M
	var err error
	var id int64
	n := 0
	for err == nil {
		err = db.gamesid.Find(nil).Limit(100).Skip(n * 100).All(&ids)
		var gameID int64
		for _, v := range ids {
			id, ok := v["gameid"]
			if !ok {
				logger.Println("info: failed get gameid", v)
				continue
			}
			gameID, ok = id.(int64)
			if !ok {
				logger.Println("info: failed cast gameid", v)
				continue
			}
			db.lolCache.AddGame(gameID)
		}
		if id == gameID {
			break
		}
		id = gameID
		n++
	}
	count := 0
	db.lolCache.(*memCache).games.Range(func(key interface{}, value interface{}) bool {
		count++
		logger.Println(key)
		return true
	})
	logger.Println("info: Found ", count, " games")
	var players []Player
	n = 0
	var pid int64
	for err == nil {
		err = db.playersVisited.Find(nil).Limit(100).Skip(n * 100).All(&players)
		for _, p := range players {
			db.lolCache.VisitedPlayer(p.AccountID)
		}
		n++
		if len(players) == 0 || pid == players[0].AccountID {
			break
		}
		pid = players[0].AccountID
	}
	n = 0
	for err == nil {
		err = db.players.Find(nil).Limit(100).Skip(n * 100).All(&players)
		for _, p := range players {
			db.lolCache.VisitedPlayer(p.AccountID)
		}
		n++
		if len(players) == 0 || pid == players[0].AccountID {
			break
		}
		pid = players[0].AccountID
	}
	playersFound := 0
	db.lolCache.(*memCache).visited.Range(func(key interface{}, value interface{}) bool {
		playersFound++
		return true
	})
	db.lolCache.(*memCache).toVisit.Range(func(key interface{}, value interface{}) bool {
		playersFound++
		return true
	})
	logger.Println("info: Found:", playersFound, " players")
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
