package lol

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type lolMongo struct {
	session        *mgo.Session
	db             *mgo.Database
	games          *mgo.Collection
	players        *mgo.Collection
	playersVisited *mgo.Collection
}

var _ lolStorer = &lolMongo{}

func NewLolMongo() (lolStorer, error) {
	session, err := mgo.Dial("localhost:27017")
	// session, err := mgo.Dial("linode.jhrb.us:27017")
	if err != nil {
		return nil, err
	}
	n, _ := session.DB("lol").C("games").Count()
	logger.Printf("trace:games: %d", n)
	n, _ = session.DB("lol").C("players").Count()
	logger.Printf("trace:left to visit: %d", n)
	n, _ = session.DB("lol").C("playersvisited").Count()
	logger.Printf("trace: visited: %d", n)
	return &lolMongo{
		session,
		session.DB("lol"),
		session.DB("lol").C("games"),
		session.DB("lol").C("players"),
		session.DB("lol").C("playersvisited"),
	}, nil
}

func (db *lolMongo) GetGame(gameID int64, currentPlatformID string) (Game, error) {
	var game Game
	// err := db.db.Read("games", fmt.Sprintf("%d_%s", gameID, currentPlatformID), &game)
	err := db.games.Find(bson.M{"gameid": gameID, "platformid": currentPlatformID}).One(&game)
	return game, err
}

func (db *lolMongo) SaveGame(game Game, currentPlatformID string) error {
	// return db.db.Write("games", fmt.Sprintf("%d_%s", game.GameID, currentPlatformID), &game)
	n, _ := db.games.Find(bson.M{"gameid": game.GameID, "platformid": currentPlatformID}).Count()
	if n == 0 {
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

type PlayerWithVisited struct {
	Player
	Visited bool
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

func (db *lolMongo) VisitPlayer(p Player) error {
	err := db.players.Remove(bson.M{"accountid": p.AccountID})
	if err != nil {
		return err
	}
	return db.playersVisited.Insert(&p)
}
