package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	bgPixel background
	mapCoordinateFlowers []coordinate
)

type background struct {
	img   *ebiten.Image
	grass pixelImage
}

type pixelImage struct {
	x1, y1, x2, y2 int
}

func drawBackGround(screen *ebiten.Image) {
	//background tiles
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

	//flowers
	drawFlowers(screen)

}

//init will have a randomized value once, otherwise on each frame the flower are mooving
func initCoordinateFlowers(){
	for j := 0; j <= 3; j++ {
		for i := 0; i <= 6; i++ {
			mapCoordinateFlowers = append(mapCoordinateFlowers, coordinate{
				x: randomNumberBeetween(float64(i*50)-10,float64(i*50)+10),
				y: randomNumberBeetween(float64(j*50)-10,float64(j*50)+10),
			})
		}
	}
}

func drawFlowers(screen *ebiten.Image) {

	for _,c := range mapCoordinateFlowers {
		flowerOpt := &ebiten.DrawImageOptions{}
		flowerOpt.GeoM.Translate(c.x, c.y)
			flowerOpt.GeoM.Scale(3, 3)
			screen.DrawImage(flowerImage, flowerOpt)
	}
}

func drawBgTiles(x, y float64, pxI pixelImage) (*ebiten.Image, *ebiten.DrawImageOptions) {

	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(2, 2)
	opChar.GeoM.Translate(x, y)
	return bgPixel.img.SubImage(image.Rect(pxI.x1, pxI.y1, pxI.x2, pxI.y2)).(*ebiten.Image), opChar
}
