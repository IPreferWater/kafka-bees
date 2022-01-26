package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	bgPixel background
)

type background struct {
	img   *ebiten.Image
	grass pixelImage
}

type pixelImage struct {
	x1, y1, x2, y2 int
}

func drawBackGround(screen *ebiten.Image) {
	grass, optGrass := drawBgTiles(0, 0, bgPixel.grass)
	x, y := getNumberOfTilesToDraw(screenWidth, screenHeight, 64)
	screen.DrawImage(grass, optGrass)
	for j := 0; j <= y; j++ {
		optGrass.GeoM.Translate(0, float64(j*64))
		//draw the first tile at x=0
		screen.DrawImage(grass, optGrass)
		for i := 0; i <= x; i++ {
			optGrass.GeoM.Translate(64, 0)
			screen.DrawImage(grass, optGrass)
		}
		optGrass.GeoM.Reset()
		optGrass.GeoM.Scale(2, 2)
	}
}

func drawBgTiles(x, y float64, pxI pixelImage) (*ebiten.Image, *ebiten.DrawImageOptions) {

	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(2, 2)
	opChar.GeoM.Translate(x, y)
	return bgPixel.img.SubImage(image.Rect(pxI.x1, pxI.y1, pxI.x2, pxI.y2)).(*ebiten.Image), opChar
}
