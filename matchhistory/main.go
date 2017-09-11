package main

import (
	"log"
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
	start := time.Now()
	lol.DefaultClient().GetCache().LoadAllGameIDS()
	log.Println("Took: ", time.Since(start), " to load all playerids and gameids")
	if len(os.Args) != 2 {
		log.Println("Invalid args:", os.Args)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "-w":
		lol.DefaultClient().GetCache().Stats()
	case "seed":
		seed(202988570)
	case "scrap":
		if err := scrap(); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Println(os.Args[1])
	}
}
