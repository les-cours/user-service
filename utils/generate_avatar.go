package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func RandomColorGenerator() (string, string) {
	//generate random color in HSL format
	rand.New(rand.NewSource(time.Now().UnixNano()))
	h := rand.Float64() * 360.0
	rand.New(rand.NewSource(time.Now().UnixNano()))
	s := 0.3 + rand.Float64()*(0.7-0.3)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	//convert HSL color to RGB format
	l := 0.3 + rand.Float64()*(0.7-0.3)
	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64

	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r := uint8(math.Round((Rnot + m) * 255))
	g := uint8(math.Round((Gnot + m) * 255))
	b := uint8(math.Round((Bnot + m) * 255))

	//Check if the background color goes better with the white or with the black text
	backgroundBrightness := CalculateBrightness(int(r), int(g), int(b))
	whiteTextBrightness := CalculateBrightness(255, 255, 255)
	blackTextBrightness := CalculateBrightness(0, 0, 0)

	if math.Abs(backgroundBrightness-whiteTextBrightness) > math.Abs(backgroundBrightness-blackTextBrightness) {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b), "#ffffff"
	} else {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b), "#000000"
	}
}

func CalculateBrightness(r, g, b int) float64 {
	return float64((r*299 + g*587 + b*114) / 1000)
}

func GenerateAvatar(firstname, lastname string) (string, error) {
	// Generating SVG file with user initials and random background color
	var userInitials = string(firstname[0]) + string(lastname[0])
	var hexBgColor, hexTextColor = RandomColorGenerator()
	var svg = `
    <svg  
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      width="64px"
      height="64px"
      x="0px"
        y="0px"
      viewBox="0 0 64 64"
      version="1.1"
    >
      <defs>
        <style type="text/css">
          @import url("https://fonts.googleapis.com/css2?family=Inter:wght@700");
        </style>
      </defs>
      <rect fill="` + hexBgColor + `" cx="32" width="64" height="64" cy="32" r="32" />
      <text
        x="50%"
        y="50%"
        style="          
          line-height: 1;
          font-family: 'Inter', sans-serif;
          font-weight: 700;
        "
        alignment-baseline="middle"
        text-anchor="middle"
        font-size="28"
        font-weight="400"
        dy=".1em"
        dominant-baseline="middle"
        fill="` + hexTextColor + `"
      >` + userInitials + `
      </text>
  </svg>	
  `

	return svg, nil
}
