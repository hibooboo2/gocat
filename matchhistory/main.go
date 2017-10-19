package main

import (
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/hibooboo2/lol"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	// defer lol.DefaultClient().Close()
	lol.SetLogLevel(colog.LDebug)
	if len(os.Args) < 2 {
		log.Println("Invalid args:", os.Args)
		os.Exit(0)
	}
	switch os.Args[1] {
	case "exp":
		// sum := lol.DefaultClient().Summoners().ByName("Sir Yogi Bear")
		// log.Println(sum)
		// log.Printf("%##v\n", lol.DefaultClient().Mastery().All(sum.ID)[0])
		log.Println(lol.DefaultClient().Mastery().All(lol.DefaultClient().Summoners().ByName(os.Args[2]).ID))
	case "pentas":
		c, err := lol.NewClient(lol.NA)
		if err != nil {
			log.Fatalln(err)
		}
		g := c.Spectator().GameSummonerName(os.Args[2])
		summonerStats := func(summonerName string) {
			sum := lol.DefaultClient().Summoners().ByName(summonerName)
			accID := sum.AccountID
			log.Println(sum)
			games, err := lol.DefaultClient().GetAllGames(accID, lol.NA1)
			if err != nil {
				log.Println(err)
			}
			var kills, doubles, triples, quads, pentas int
			var wins, loss float64
			for _, g := range games {
				for _, p := range g.ParticipantIdentities {
					if p.Player.AccountID == accID {
						for _, part := range g.Participants {
							if part.ParticipantID == p.ParticipantID {
								pentas += part.Stats.PentaKills
								doubles += part.Stats.DoubleKills
								kills += part.Stats.Kills
								triples += part.Stats.TripleKills
								quads += part.Stats.QuadraKills
								if part.Stats.Win {
									wins++
								} else {
									loss++
								}
							}
						}
					}
				}
			}
			log.Printf("%s Kills: %d Doubles: %d Triples: %d Quads: %d Pentas: %d ", sum.Name, kills, doubles, triples, quads, pentas)
			log.Printf(" Wins: %0f Loss: %0f Ration: %f", wins, loss, float64(wins/(loss+wins)))
		}

		if g == nil || g.GameID == 0 {
			log.Println("No active game for:", os.Args[2])
			summonerStats(os.Args[2])
			return
		}
		var sums []string
		for _, p := range g.Participants {
			sums = append(sums, p.SummonerName)
		}
		for _, summonerName := range sums {
			summonerStats(summonerName)
		}
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
		// db, err := lol.NewLolMongo("", 0)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// log.Println(db.TransferToAnother("", 27027))
	case "gameidgen":
		db, err := lol.NewLolMongo("192.168.1.170", 27017)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Genning ids..")
		db.EnsureIndexes()
		log.Println("Done")
	default:
		log.Println(os.Args[1])
	}
}
