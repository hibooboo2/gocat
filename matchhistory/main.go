package main

import (
	"log"
	"os"

	"github.com/hibooboo2/gocat/lol"
)

var c *lol.Client

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	c, err = lol.NewClient()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer c.Close()
	log.Println("Starting scraping forever...")
	if len(os.Args) != 2 {
		log.Println("Invalid args:", os.Args)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-w":
		c.GetCache().Stats()
	case "seed":
		seed(220448739)
	case "scrap":
		if err := scrap(); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Println(os.Args[1])
	}
}

func seed(accountID int64) {
	log.Println("Seeding....")
	games, err := c.GetAllGamesLimitPatch(accountID, "NA1", "7", 3000)
	if err != nil {
		log.Fatalln("Failed to get history: ", err)
	}
	sums := make(map[int64]lol.Player)
	for _, game := range games {
		game, _ := c.WebMatch(game.GameID, game.ParticipantIdentities[0].Player.CurrentPlatformID, true)
		for _, sum := range game.ParticipantIdentities {
			sums[sum.Player.AccountID] = sum.Player
		}
	}

	c.GetCache().StorePlayer(sums[accountID])
	delete(sums, accountID)
	for _, sum := range sums {
		c.GetCache().StorePlayer(sum)
	}
	log.Println("Seeded")
}
