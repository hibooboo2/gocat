package main

import (
	"log"

	"github.com/hibooboo2/gocat/lol"
)

var c = lol.NewClient()

func main() {
	// repeat for all summoners
	accountID := int64(34178787) //sir yogi bear
	// accountID := int64(44278412)
	// accountID := int64(42795563)
	// accountID := int64(205659322) // Sir fxwright
	games, err := c.GetAllGames(accountID, "NA1")
	// cache := lol.NewLolCache()
	// games, err := cache.GetGames(accountID, "NA1")

	if err != nil {
		log.Fatalln("Failed to get history: ", err)
	}
	sums := make(map[int64]lol.Player)
	var pentas int
	for _, game := range games {
		game, _ := c.WebMatch(game.GameID, game.ParticipantIdentities[0].Player.CurrentPlatformID)
		for _, sum := range game.ParticipantIdentities {
			sums[sum.Player.AccountID] = sum.Player
			if sum.Player.AccountID == accountID {
				for _, particpant := range game.Participants {
					if particpant.ParticipantID == sum.ParticipantID {
						if particpant.Stats.PentaKills > 0 {
							pentas += particpant.Stats.PentaKills
						}
					}
				}
			}
		}

	}
	log.Println(pentas)
	sumVisited := make(map[int64]bool)
	for {
		for _, sum := range sums {
			if sumVisited[sum.AccountID] {
				continue
			}
			games, err := c.GetAllGames(sum.AccountID, sum.CurrentPlatformID)
			if err != nil {
				continue
			}
			for _, game := range games {
				game, err := c.WebMatch(game.GameID, game.PlatformID)
				if err != nil {
					log.Println("err: Failed to get match:", game.GameID, err)
				}
				for _, sum := range game.ParticipantIdentities {
					sums[sum.Player.AccountID] = sum.Player
				}
			}
			sumVisited[sum.AccountID] = true
		}
	}
}
