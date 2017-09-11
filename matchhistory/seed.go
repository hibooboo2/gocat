package main

import (
	"log"

	"github.com/hibooboo2/gocat/lol"
)

func seed(accountID int64) {
	c := lol.DefaultClient()
	log.Println("Seeding....")
	games, err := c.GetAllGamesLimitPatch(accountID, lol.NA1, "7", 700)
	if err != nil {
		log.Fatalln("Failed to get history: ", err)
	}
	sums := make(map[int64]lol.Player)
	for _, game := range games {
		game, _ := c.WebMatch(game.GameID, game.ParticipantIdentities[0].Player.CurrentPlatformID, true)
		if game.Cached {
			continue
		}
		for _, sum := range game.ParticipantIdentities {
			sums[sum.Player.AccountID] = sum.Player
		}
	}
	c.GetCache().VisitedPlayer(accountID)
	delete(sums, accountID)
	for _, sum := range sums {
		c.GetCache().StorePlayer(sum)
	}
	log.Println("Seeded")
}
