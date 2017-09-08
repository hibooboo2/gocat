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

	var matchesFarmed, sumsVisited int
	var player lol.Player
	var err error
	log.Println("Starting scraping forever...")
	for matchesFarmed < 5000000 {
		player, err = c.GetCache().GetPlayerToVisit()
		if err != nil {
			break
		}
		games, err := c.GetAllGamesLimitPatch(player.AccountID, player.CurrentPlatformID, "7.17.")
		if err != nil {
			continue
		}
		// log.Println("Got games for: ", player.SummonerName, len(games))
		var game *lol.Game
		for _, g := range games {
			id := g.GameID
			game, err = c.WebMatch(g.GameID, g.PlatformID)
			if !game.Cached {
				matchesFarmed++
			}
			fmt.Fprintf(os.Stdout, "\rSum:\t%s\tGame:\t%d\tMatchesFarmed\t%d\tSumsVisited\t%d", player.SummonerName, id, matchesFarmed, sumsVisited)
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
		sumsVisited++
		if err != nil {
			log.Println(err)
		}
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
