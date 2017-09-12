package main

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/hibooboo2/gocat/lol"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	defer lol.DefaultClient().Close()
	lol.SetLogLevel(colog.LTrace)
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
	case "gameidgen":
		db, err := lol.NewLolMongoWAccess("localhost", 0)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Genning ids..")
		db.GameIDSToIDTable()
		log.Println("Done")
	case "random":
		db, err := lol.NewLolMongoWAccess("dev.jhrb.us", 27217)
		if err != nil {
			log.Fatalln(err)
		}
		db.GetGameRan()
	default:
		log.Println(os.Args[1])
	}
}
