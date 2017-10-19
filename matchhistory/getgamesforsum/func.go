package main

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/hibooboo2/lol"
)

func main() {
	lol.SetLogLevel(colog.LError)
	c, err := lol.NewClient(lol.NA)
	handleErr(err)
	defer c.GetCache().Close()
	p := c.GetCache().GetPlayerToVisit()
	if p == 0 {
		log.Fatalln("No player gotten")
	}
	games, err := c.GetAllGamesLimitPatch(p, lol.NA1, "7.17", 20)
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
	delete(sums, p)
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
