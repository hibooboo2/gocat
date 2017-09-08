package main

import (
	"fmt"
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
	var matchesFarmed int
	player, err := c.GetCache().GetPlayerToVisit()
	for err != nil && matchesFarmed > 5000000 {
		log.Println(player.SummonerName)
		games, err := c.GetAllGamesLimitPatch(player.AccountID, player.CurrentPlatformID, "7.17.")
		if err != nil {
			continue
		}
		// log.Println("Got games for: ", player.SummonerName, len(games))
		var game *lol.Game
		for _, g := range games {
			id := g.GameID
			game, err = c.WebMatch(g.GameID, g.PlatformID)
			matchesFarmed++
			// log.Println("Farmed", matchesFarmed)
			if err != nil {
				log.Println("err: Failed to get match:", id, err)
				continue
			}
			for _, sum := range game.ParticipantIdentities {
				if sum.Player.AccountID != player.AccountID {
					c.GetCache().StorePlayer(sum.Player, false)
				}
			}
			// log.Println("Got game: ", game.GameID)
		}
		err = c.GetCache().VisitPlayer(player)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stdout, "\rTotal Games for Sum: %d Sum: %s", len(games), player.SummonerName)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func seed() {
	// accountID := int64(34178787) //sir yogi bear
	accountID := int64(44278412)
	// accountID := int64(42795563)
	// accountID := int64(205659322) // Sir fxwright
	games, err := c.GetAllGamesLimitPatch(accountID, "NA1", "7.17")
	if err != nil {
		log.Fatalln("Failed to get history: ", err)
	}
	var thisSum lol.Player
	for _, game := range games {
		game, _ := c.WebMatch(game.GameID, game.ParticipantIdentities[0].Player.CurrentPlatformID)
		for _, sum := range game.ParticipantIdentities {
			if sum.Player.AccountID != accountID {
				c.GetCache().StorePlayer(sum.Player, false)
			} else {
				thisSum = sum.Player
			}
		}
	}
	c.GetCache().StorePlayer(thisSum, true)
}
