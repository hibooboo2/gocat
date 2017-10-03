package lol

type ItemList struct {
	Data map[string]Item `json:"data"`
	Tree []struct {
		Header string   `json:"header"`
		Tags   []string `json:"tags"`
	} `json:"tree"`
	Groups []struct {
		Key             string `json:"key"`
		MaxGroupOwnable string `json:"MaxGroupOwnable"`
	} `json:"groups"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type Item struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	SanitizedDescription string `json:"sanitizedDescription"`
	Image                struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Tags []string `json:"tags"`
	Gold struct {
		Base        int  `json:"base"`
		Total       int  `json:"total"`
		Sell        int  `json:"sell"`
		Purchasable bool `json:"purchasable"`
	} `json:"gold"`
	Plaintext string          `json:"plaintext"`
	Depth     int             `json:"depth"`
	From      []string        `json:"from"`
	Into      []string        `json:"into"`
	Maps      map[string]bool `json:"maps"`
	Effect    struct {
		Effect1Amount string `json:"Effect1Amount"`
		Effect2Amount string `json:"Effect2Amount"`
	} `json:"effect"`
	Stats struct {
		FlatArmorMod int `json:"FlatArmorMod"`
	} `json:"stats"`
}
