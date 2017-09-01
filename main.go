package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/nfnt/resize"
	"github.com/nsf/termbox-go"
	"github.com/stroborobo/ansirgb"
)

func main() {
	// Read image from file that already exists
	// existingImageFile, err := os.Open("img.png")
	// if err != nil {
	// 	// Handle error
	// }
	// defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	// imageData, imageType, err := image.Decode(existingImageFile)
	// if err != nil {
	// 	panic(err)
	// }
	// if imageType != "png" {
	// 	panic("Only png is allowed")
	// }
	// fmt.Println(imageData)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	// existingImageFile.Seek(0, 0)

	// Alternatively, since we know it is a png already
	// we can call png.Decode() directly
	champsMap := getChamps() //"https://na1.api.riotgames.com/lol/static-data/v3/champions?api_key=RGAPI-ed1dfe8d-8adb-4283-8a3a-094e5dddb3df"
	champNames := []string{}
	for name := range champsMap {
		champNames = append(champNames, name)
	}
	sort.Strings(champNames)
	// fmt.Println(loadedImage.Bounds(), loadedImage.At(0, 0))
	// r, g, b, a := loadedImage.At(0, 0).RGBA()
	if termbox.Init() != nil {
		panic("Failed to init termbox")
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCursor(0, 50)
	termbox.Sync()
	x := 0
	for _, champName := range champNames {
		champ := champsMap[champName]
		termbox.SetCursor(0, 50)
		log.Println(champ.Key)
		drawChampHead(champ)
		time.Sleep(time.Millisecond * 100)
		x++
		if x > 10 {
			break
		}
	}
	time.Sleep(time.Second * 1)
}

func drawChampHead(champ Champ) {
	resp, err := http.Get(fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/6.24.1/img/champion/%s.png", champ.Key))
	if err != nil {
		termbox.SetCursor(0, 50)
		log.Println("Errored getting champ:", err)
		return
	}
	defer resp.Body.Close()

	loadedImage, err := png.Decode(resp.Body)
	if err != nil {
		// Handle error
	}
	if loadedImage != nil {
		smallImg := resize.Resize(80, 80, loadedImage, resize.Lanczos3)
		drawImage(smallImg)
		f, err := os.Create(champ.Key + ".png")
		defer f.Close()
		if err == nil {
			png.Encode(f, smallImg)
		}
	}
}

func drawImage(img image.Image) error {
	if img == nil {
		termbox.SetCursor(0, 50)
		log.Println("Nil img")
		return nil
	}
	// img = grayscale.Convert(img, grayscale.ToGrayLuminance)
	img = resize.Resize(40, 40, img, resize.Lanczos3)
	min := img.Bounds().Min
	max := img.Bounds().Max
	var c, c2 color.RGBA
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c.A = 255
			c.B = uint8(b / 0x101)
			c.G = uint8(g / 0x101)
			c.R = uint8(r / 0x101)
			r, g, b, _ = img.At(x, y+1).RGBA()
			c2.A = 255
			c2.B = uint8(b / 0x101)
			c2.G = uint8(g / 0x101)
			c2.R = uint8(r / 0x101)

			termbox.SetCell(x, y, '▄', termbox.Attribute(ansirgb.Convert(&c2).Code), termbox.Attribute(ansirgb.Convert(&c).Code))
			termbox.SetCell(x+50, y, '▄', termbox.Attribute(ansirgb.Convert(img.At(x, y+1)).Code), termbox.Attribute(ansirgb.Convert(img.At(x, y)).Code))
		}
	}
	termbox.Flush()
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
