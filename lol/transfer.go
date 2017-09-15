package lol

import mgo "gopkg.in/mgo.v2"

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
	// var ids []int64
	// logger.Println(db.games.Find(nil).Distinct("gameid", &ids))
	// f, _ := os.Create("gameids")
	// data, _ := json.Marshal(ids)
	// f.Write(data)
	// f.Close()
	// logger.Println(len(ids))
	// var games []Game
	// logger.Println(db.games.Find(nil).All(&games))
	// gamesMap := make(map[int64]Game)
	// f, _ := os.Open("allgames.json")
	// defer f.Close()
	// err := json.NewDecoder(f).Decode(&gamesMap)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, game := range gamesMap {
	// 	err = db.db.C("new_games").Insert(game)
	// 	if err != nil {
	// 		logger.Println("err:", err)
	// 	}
	// 	fmt.Fprintf(os.Stdout, "GameID:\t %d\r", game.GameID)
	// }
	// for _, game := range games {
	// 	gamesMap[game.GameID] = game
	// }
	// dupGamesCount := make(map[int64]int)
	// for _, game := range games {
	// dupGamesCount[game.GameID]++
	// }
	// data, _ = json.Marshal(gamesMap)
	// f, _ = os.Create("allgames.json")
	// f.Write(data)
	// f.Close()
	// os.Exit(0)
	// removed := 0
	// toskip := 0
	// for id, n := range dupGamesCount {
	// 	if n <= 1 {
	// 		logger.Println("Skipped: ", id)
	// 		toskip++
	// 		continue
	// 	} else {
	// 		continue
	// 	}
	// var err error
	// for n > 1 {
	// 	err = db.games.Remove(bson.M{"gameid": id})
	// 	n--
	// 	if err != nil {
	// 		logger.Println("err: errored on game removal", id)
	// 		continue
	// 	}
	// 	removed++
	// 	fmt.Fprintf(os.Stdout, "n: %d \tRemoved: %d\r", n, id)
	// }
	// if err != nil {
	// 	fmt.Fprintf(os.Stdout, "Failed Dups for: %d\r", id)
	// 	continue
	// }
	// fmt.Fprintf(os.Stdout, "Removed: %d Removed Dups for: %d\r", removed, id)
	// }
	// logger.Printf("To SKip: %d Total %d Records TO remove: %d", toskip, len(games), len(games)-toskip)
	// gamesCount, _ := db.games.Count()
	// newGames, _ := db.db.C("new_games").Count()
	// logger.Println("After games de dup loop:", gamesCount, "Games in map:", len(gamesMap), "Removed:", gamesCount-newGames)
	logger.Println(db.games.EnsureIndex(mgo.Index{
		Key:      []string{"gameid"},
		DropDups: true,
		Unique:   true,
	}))
	logger.Println(db.players.EnsureIndex(mgo.Index{
		Key:      []string{"accountid"},
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
