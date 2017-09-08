package lol

import (
	log "log"
	"os"

	"github.com/comail/colog"
)

var logger *log.Logger

func init() {
	cologer := colog.NewCoLog(os.Stdout, "lolapi:", log.Lshortfile)
	cologer.SetDefaultLevel(colog.LTrace)
	// cologer.SetMinLevel(colog.LDebug)
	logger = log.New(cologer, "", 0)
	logger.Println("Logger initialized for lolapi.")
}
