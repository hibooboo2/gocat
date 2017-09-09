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
	games, err := c.GetAllGamesLimitPatch(p.AccountID, "NA", "7.17")
	handleErr(err)

	var found int
	for _, g := range games {
		game, _ := c.WebMatch(g.GameID, g.PlatformID)
		if !game.Cached {
			found++
		}
	}
	log.Printf(`{"found":%d}`, found)
}

func handleErr(err error) {
	if err != nil {
		log.Printf(`{"error":"%s"}`, err)
		os.Exit(1)
	}
}
