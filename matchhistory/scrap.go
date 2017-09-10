package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hibooboo2/gocat/lol"
)

func scrap() error {
	c := lol.DefaultClient()
	var matchesFarmed, sumsVisited int
	var player lol.Player
	var err error
	allSumsThisSession := make(map[int64]lol.Player)
	for matchesFarmed < 5000000 {
		player, err = c.GetCache().GetPlayerToVisit()
		if err != nil {
			return err
		}
		sumsVisited++
		games, err := c.GetAllGamesLimitPatch(player.AccountID, player.CurrentPlatformID, "7.17.", 20)
		if err != nil {
			log.Println(err)
			continue
		}
		var game *lol.Game
		sums := make(map[int64]lol.Player)
		for _, g := range games {
			id := g.GameID
			if c.HaveMatch(id) {
				fmt.Fprintf(os.Stdout, "\rSum:\t%s\tGame:\t%d\tMatchesFarmed\t%d\tSumsVisited\t%d", player.SummonerName, id, matchesFarmed, sumsVisited)
				continue
			}
			game, err = c.WebMatch(g.GameID, g.PlatformID, false)
			if !game.Cached {
				matchesFarmed++
			}
			fmt.Fprintf(os.Stdout, "\rSum:\t%s\tGame:\t%d\tMatchesFarmed\t%d\tSumsVisited\t%d", player.SummonerName, id, matchesFarmed, sumsVisited)
			if err != nil {
				log.Println("err: Failed to get match:", id, err)
				continue
			}
			for _, sum := range game.ParticipantIdentities {
				_, ok := allSumsThisSession[sum.Player.AccountID]
				if !ok {
					_, ok = sums[sum.Player.AccountID]
					if !ok && sum.Player.AccountID != player.AccountID {
						sums[sum.Player.AccountID] = sum.Player
					}
				}
			}
		}
		for _, sum := range sums {
			c.GetCache().StorePlayer(sum)
		}

	}
	return err
}
