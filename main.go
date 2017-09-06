package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/gdamore/tcell"
	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

var s tcell.Screen

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, champ := range sortedChamps() {
		s = append(s, prompt.Suggest{Text: champ.Id, Description: champ.Title})
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
		fmt.Printf("%+v", champ)
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

func drawChamp(champ Champ) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	s.Init()
	defer s.Fini()
	resp, err := http.Get(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/6.24.1/img/champion/%s.png", champ.Id))
	if err != nil {
		// log.Println("Errored getting champ:", err)
		return
	}
	defer resp.Body.Close()

	loadedImage, err := png.Decode(resp.Body)
	if err == nil {
		if loadedImage != nil {
			drawImage(loadedImage, false, 30, 0, 0)
		}
	}
	time.Sleep(time.Second * 3)
}

func drawChampHead(champ Champ) {
	resp, err := http.Get(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/6.24.1/img/champion/%s.png", champ.Id))
	if err != nil {
		// log.Println("Errored getting champ:", err)
		return
	}
	defer resp.Body.Close()

	loadedImage, err := png.Decode(resp.Body)
	if err != nil {
		// Handle error
	}
	if loadedImage != nil {
		drawImage(loadedImage, false, 120, 0, 0)
		// smallImg := resize.Resize(80, 80, loadedImage, resize.Lanczos3)
		// drawImage(smallImg)
		// f, err := os.Create(champ.Key + ".png")
		// defer f.Close()
		// if err == nil {
		// 	png.Encode(f, smallImg)
		// }
	}
}

func drawImage(img image.Image, blackAndWhite bool, size uint, xoff, yoff int) error {
	if img == nil {
		log.Println("Nil img")
		return nil
	}
	if blackAndWhite {
		img = grayscale.Convert(img, grayscale.ToGrayLuminance)
	}
	img = resize.Resize(size, size, img, resize.Lanczos3)
	min := img.Bounds().Min
	max := img.Bounds().Max
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y += 2 {
			st := tcell.StyleDefault
			r, g, b, _ := img.At(x, y).RGBA()
			st = st.Background(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			r, g, b, _ = img.At(x, y+1).RGBA()
			st = st.Foreground(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			s.SetCell(x+xoff, y/2+yoff, st, '▄')
		}
	}
	s.Sync()
	return nil
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
