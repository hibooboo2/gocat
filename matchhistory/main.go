package main

import (
	"log"
	"os"

	"github.com/hibooboo2/gocat/lol"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	defer lol.DefaultClient().Close()
	log.Println("Starting scraping forever...")
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
