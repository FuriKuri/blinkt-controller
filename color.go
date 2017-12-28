package main

var Colors map[string]LedColor

type LedColor struct {
	Red   int32 `json:"red"`
	Green int32 `json:"green"`
	Blue  int32 `json:"blue"`
	Led   int32 `json:"led"`
}

func init() {
	Colors = map[string]LedColor{
		"red":    LedColor{Red: 255, Blue: 0, Green: 0, Led: 0},
		"brown":  LedColor{Red: 139, Blue: 19, Green: 69, Led: 1},
		"orange": LedColor{Red: 255, Blue: 0, Green: 69, Led: 2},
		"yellow": LedColor{Red: 255, Blue: 0, Green: 255, Led: 3},
		"green":  LedColor{Red: 0, Blue: 0, Green: 255, Led: 4},
		"blue":   LedColor{Red: 0, Blue: 255, Green: 0, Led: 5},
		"violet": LedColor{Red: 128, Blue: 128, Green: 0, Led: 6},
		"grey":   LedColor{Red: 255, Blue: 100, Green: 255, Led: 7},
	}
}
