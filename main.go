package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gdamore/tcell"
)

var s tcell.Screen

func completer(d prompt.Document) []prompt.Suggest {
	var lim int
	if d.GetWordBeforeCursor() == "" {
		lim = 5
	}
	s := []prompt.Suggest{}
	for i, champ := range sortedChamps() {
		s = append(s, prompt.Suggest{Text: champ.Id, Description: champ.Title})
		if lim != 0 && i >= lim {
			break
		}
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(in string) {
	fmt.Println("Your input: " + in)
}

func ctrlC(b *prompt.Buffer) {
}

func main() {
	p := prompt.New(executor, completer, prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ControlC, Fn: ctrlC}),
		prompt.OptionTitle("Pick your Champ: "),
		prompt.OptionPrefix("Pick your Champ: "),
		prompt.OptionMaxSuggestion(20))
	champs := getChamps()
	for {
		t := p.Input()
		champ, ok := champs[t]
		if !ok {
			break
		}
		drawChamp(champ)

		fmt.Fprintf(os.Stdout, "%+v\n", champ)
	}
}

func main2() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	s.Init()
	defer s.Fini()

	x := 0
	for _, champ := range sortedChamps() {
		// log.Println(champ.Key)
		drawChampHead(champ)
		time.Sleep(time.Millisecond * 3000)
		x++
		if x > 5 {
			break
		}
	}
	time.Sleep(time.Second * 1)
}

type Champ struct {
	Id    string
	Key   string
	Name  string
	Title string
}

var champsOnce sync.Once
var champsMasterMap map[string]Champ

func getChamps() map[string]Champ {
	champsOnce.Do(func() {
		resp, _ := http.Get("https://ddragon.leagueoflegends.com/realms/na.json")
		var relms struct {
			Cdn string
			V   string
			N   struct {
				Champion string
			}
		}
		err := json.NewDecoder(resp.Body).Decode(&relms)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()

		var champData struct {
			Data map[string]Champ
		}
		resp, _ = http.Get(relms.Cdn + "/" + relms.V + "/data/en_US/champion.json")
		err = json.NewDecoder(resp.Body).Decode(&champData)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(relms)
			fmt.Println(champData.Data)
		}
		resp.Body.Close()
		champsMasterMap = champData.Data
	})
	return champsMasterMap
}

func sortedChamps() []Champ {
	champsMap := getChamps()
	champNames := []string{}
	for name := range champsMap {
		champNames = append(champNames, name)
	}
	sort.Strings(champNames)
	champs := []Champ{}
	for _, champ := range champNames {
		champs = append(champs, champsMap[champ])
	}
	return champs
}
