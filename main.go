package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nfnt/resize"
)

var s tcell.Screen

func main() {
	champsMap := getChamps() //"https://na1.api.riotgames.com/lol/static-data/v3/champions?api_key=RGAPI-ed1dfe8d-8adb-4283-8a3a-094e5dddb3df"
	champNames := []string{}
	for name := range champsMap {
		champNames = append(champNames, name)
	}
	sort.Strings(champNames)
	// fmt.Println(loadedImage.Bounds(), loadedImage.At(0, 0))
	// r, g, b, a := loadedImage.At(0, 0).RGBA()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	s.Init()
	defer s.Fini()
	st := tcell.StyleDefault
	st = st.Background(tcell.NewHexColor(0xfea0ab))
	s.SetCell(15, 15, st, 'A')
	s.Show()
	// log.Println(s.Size())

	x := 0
	for _, champName := range champNames {
		champ := champsMap[champName]
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

func drawChampHead(champ Champ) {
	resp, err := http.Get(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/6.24.1/img/champion/%s.png", champ.Key))
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
		drawImage(loadedImage)
		// smallImg := resize.Resize(80, 80, loadedImage, resize.Lanczos3)
		// drawImage(smallImg)
		// f, err := os.Create(champ.Key + ".png")
		// defer f.Close()
		// if err == nil {
		// 	png.Encode(f, smallImg)
		// }
	}
}

func drawImage(img image.Image) error {
	if img == nil {
		log.Println("Nil img")
		return nil
	}
	// img = grayscale.Convert(img, grayscale.ToGrayLuminance)
	img = resize.Resize(80, 80, img, resize.Lanczos3)
	min := img.Bounds().Min
	max := img.Bounds().Max
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			st := tcell.StyleDefault
			r, g, b, _ := img.At(x, y).RGBA()
			st = st.Background(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			r, g, b, _ = img.At(x, y+1).RGBA()
			st = st.Foreground(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			s.SetCell(x, y/2, st, '▄')
			y++
			// termbox.SetCell(x, y, '▄', termbox.Attribute(ansirgb.Convert(&c2).Code), termbox.Attribute(ansirgb.Convert(&c).Code))
			// termbox.SetCell(x+50, y, '▄', termbox.Attribute(ansirgb.Convert(img.At(x, y+1)).Code), termbox.Attribute(ansirgb.Convert(img.At(x, y)).Code))
		}
	}
	// s.Show()
	s.Sync()
	// termbox.Flush()
	return nil
}

type Champ struct {
	Id    int
	Key   string
	Name  string
	Title string
}

func getChamps() map[string]Champ {
	var champData struct {
		Data map[string]Champ
	}
	f, _ := os.Open("champs.json")
	defer f.Close()
	json.NewDecoder(f).Decode(&champData)
	return champData.Data
}
