package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hibooboo2/lol"
)

func scrap() error {
	c := lol.DefaultClient()
	var matchesFarmed, sumsVisited, fromCache int
	var accountID int64
	var err error
	allSumsThisSession := make(map[int64]lol.Player)
	allGamesThisSession := make(map[int64]struct{})
	log.Println("Starting scraping forever...")
	for matchesFarmed < 5000000 {
		accountID = c.GetCache().GetPlayerToVisit()
		if accountID == 0 {
			return errors.New("No players to visit")
		}
		sumsVisited++
		games, err := c.GetAllGamesLimitPatch(accountID, lol.NA1, "7.17.", 1000)
		if err != nil {
			log.Println(err)
			continue
		}
		var game *lol.Game
		sums := make(map[int64]lol.Player)
		for _, g := range games {
			id := g.GameID
			_, have := allGamesThisSession[id]
			if have || c.HaveMatch(id) {
				fromCache++
				fmt.Fprintf(os.Stdout, "\rSum:\t%d\tGame:\t%d\tMatchesFarmed\t %d\tMatchesFromCache\t %d\tSumsVisited\t%d", accountID, id, matchesFarmed, fromCache, sumsVisited)
				continue
			}
			game, err = c.WebMatch(g.GameID, g.PlatformID, false)
			if !game.Cached {
				matchesFarmed++
			}
			fmt.Fprintf(os.Stdout, "\rSum:\t%d\tGame:\t%d\tMatchesFarmed\t %d\tMatchesFromCache\t %d\tSumsVisited\t%d", accountID, id, matchesFarmed, fromCache, sumsVisited)
			if err != nil {
				log.Println("err: Failed to get match:", id, err)
				continue
			}
			for _, sum := range game.ParticipantIdentities {
				_, ok := allSumsThisSession[sum.Player.AccountID]
				if !ok {
					_, ok = sums[sum.Player.AccountID]
					if !ok && sum.Player.AccountID != accountID {
						sums[sum.Player.AccountID] = sum.Player
					}
				}
			}
			allGamesThisSession[id] = struct{}{}
		}
		for _, sum := range sums {
			c.GetCache().StorePlayer(sum)
		}
	}
	return err
}
