package main

import (
	"log"

	"github.com/hibooboo2/gocat/lol"
)

var c = lol.NewClient()

func main() {
	// get all game ids for an account
	// get all games for an account
	// collect all summoners from games.
	// repeat for all summoners
	accountID := int64(34178787)
	games, err := c.GetAllGames(accountID, "NA")

	if err != nil {
		log.Fatalln("Failed to get history: ", err)
	}
	sums := make(map[int64]lol.Player)
	for _, game := range games {
		match, _ := c.WebMatch(game.GameID, accountID, game.ParticipantIdentities[0].Player.CurrentPlatformID)
		for _, sum := range match.ParticipantIdentities {
			sums[sum.Player.AccountID] = sum.Player
		}
	}
	log.Println(len(sums))
	for _, player := range sums {
		log.Println(player.SummonerName)
	}
}
