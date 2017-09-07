package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/gdamore/tcell"
	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"
)

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
