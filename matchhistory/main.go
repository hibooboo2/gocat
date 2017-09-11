package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hibooboo2/gocat/lol"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	defer lol.DefaultClient().Close()
	// lol.SetLogLevel(colog.LTrace)
	if len(os.Args) != 2 {
		log.Println("Invalid args:", os.Args)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "-w":
		lol.DefaultClient().GetCache().Stats()
	case "seed":
		lol.DefaultClient().GetCache().LoadAllGameIDS()
		seed(202988570)
	case "scrap":
		lol.DefaultClient().GetCache().LoadAllGameIDS()
		if err := scrap(); err != nil {
			log.Fatalln(err)
		}
	case "transfer":
		db, err := lol.NewLolMongoWAccess("", 0)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(db.TransferToAnother("", 27027))
	case "random":
		rand.Seed(time.Now().Unix())
		start := time.Now()
		games := make([]lol.Game, 5000000)
		for i := 0; i < len(games); i++ {
			g := games[i]
			g.GameID = rand.Int63n(70000000)
			games[i] = g
		}
		for _, g := range games {
			fmt.Fprintf(os.Stdout, "ID: %d\r", g.GameID)
		}
		fmt.Printf("\nTook: %v Games: %d", time.Since(start), len(games))
	default:
		log.Println(os.Args[1])
	}
}
