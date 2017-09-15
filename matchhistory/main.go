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
	case "server":
		// Start nats
		// Start vms
		// Start requsting games Based on players in db. If no players start a seed client.
		// Scrap for matchhistory with 1 client. Others will do games.
	case "client":
		// Get games by id. Or get matchhistory by accountid.
	case "-w":
		lol.DefaultClient().GetCache().Stats()
	case "seed":
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
		db, err := lol.NewLolMongoWAccess("192.168.1.170", 27017)
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
