package main

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/hibooboo2/gocat/lol"
)

func main() {
	url, err := url.Parse(os.Getenv("REQUEST_URL"))
	handleErr(err)
	accountID, err := strconv.Atoi(url.Query().Get("accountID"))
	handleErr(err)

	c, err := lol.NewClient()
	handleErr(err)
	defer c.Close()

	games, err := c.GetAllGamesLimitPatch(int64(accountID), "NA", "7.")
	handleErr(err)

	for _, g := range games {
		c.WebMatch(g.GameID, g.PlatformID)
	}
	log.Printf(`{"found":%d}`, len(games))
}

func handleErr(err error) {
	if err != nil {
		log.Printf(`{"error":"%s"}`, err)
		os.Exit(1)
	}
}
