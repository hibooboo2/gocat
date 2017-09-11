package main

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/hibooboo2/gocat/lol"
)

func main() {
	lol.SetLogLevel(colog.LError)
	c, err := lol.NewClient()
	handleErr(err)
	defer c.Close()
	p, err := c.GetCache().GetPlayerToVisit()
	handleErr(err)
	games, err := c.GetAllGamesLimitPatch(p.AccountID, NA1, "7.17", 20)
	handleErr(err)

	var found int
	sums := make(map[int64]lol.Player)
	for _, g := range games {
		if c.HaveMatch(g.GameID) {
			continue
		}
		game, _ := c.WebMatch(g.GameID, g.PlatformID, true)
		if !game.Cached {
			found++
		}
		for _, sum := range game.ParticipantIdentities {
			sums[sum.Player.AccountID] = sum.Player
		}
	}
	delete(sums, p.AccountID)
	for _, sum := range sums {
		c.GetCache().StorePlayer(sum)
	}
	log.Printf(`{"found":%d}`, found)
}

func handleErr(err error) {
	if err != nil {
		log.Printf(`{"error":"%s"}`, err)
		os.Exit(1)
	}
}
