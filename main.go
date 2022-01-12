package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1024
	screenHeight = 600
)

const (
	tileSize = 48
)

var (
	tilesImage      *ebiten.Image
	beeImage        *ebiten.Image
	hiveImage       = ebiten.NewImage(100, 100)
	mplusNormalFont font.Face
)

type cache struct {
	addBees int
	addWasp int
}

func init() {
	tilesEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/tiles/room.png")
	if err != nil {
		log.Fatal(err)
	}

	mainCharacterEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/bee.png")
	if err != nil {
		log.Fatal(err)
	}

	tilesImage = tilesEbitenImage

	beeImage = mainCharacterEbitenImage
	hiveImage.Fill(color.White)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type Game struct {
	hives         []Hive
	mapCenterX    float64
	mapCenterY    float64
	worldSpeed    float64
	frame         int
}

func (g *Game) Update() error {
	animate(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.frame++

	secondUpdate := false
	if g.frame%60 == 0 {
		secondUpdate = true
		g.frame = 0
	}

	for indexHive, hive := range g.hives {
		hivePointer := &g.hives[indexHive]

		squareOpt := &ebiten.DrawImageOptions{}
		squareOpt.GeoM.Translate(hive.position.x, hive.position.y)
		screen.DrawImage(hiveImage, squareOpt)
		text.Draw(screen, fmt.Sprintf("%d", hive.beesCount), mplusNormalFont, int(hive.position.x)+20, int(hive.position.y)+50, color.RGBA{0, 255, 0, 255})

		for indexBee, bee := range hive.beesToCome {
			beePointer := &hive.beesToCome[indexBee]
			beeImg, beeOpts := drawBee(bee.position.x, bee.position.y)
			screen.DrawImage(beeImg, beeOpts)

			if bee.position.x >= hive.position.x && bee.position.y >= hive.position.y+10 {
				hivePointer.beesToCome = removeBeeNoOrder(hive.beesToCome, indexBee)
				hivePointer.beesCount++
				continue
			}

			if bee.position.x <= hive.position.x {
				beePointer.position.x += 0.5
			}

			if bee.position.y <= hive.position.y+10 {
				beePointer.position.y += randomNumberBeetween(0,1)
			} else {
				beePointer.position.y -= randomNumberBeetween(0,1)
			}

		}

		for indexBee, bee := range hive.beesToGo {
			beePointer := &hive.beesToGo[indexBee]
			beeImg, beeOpts := drawBee(bee.position.x, bee.position.y)
			screen.DrawImage(beeImg, beeOpts)

			if bee.position.x >= hive.position.x+200 {
				hivePointer.beesToGo = removeBeeNoOrder(hive.beesToGo, indexBee)
				hivePointer.beesCount--
				//TODO send to kafka
				continue
			}

			beePointer.position.x += randomNumberBeetween(0.5,1.5)
			beePointer.position.y += randomNumberBeetween(-1.5,1.5)

		}

		if secondUpdate {
			beesToCome := make([]Bee, hive.beesToAdd)
			for i := 0; i < hive.beesToAdd; i++ {
				beesToCome[i] = Bee{
					position: coordinate{
						//x: hive.position.x - 100,
						x : hive.position.x - randomNumberBeetween(150, 80),
						y: float64(randomNumberBeetween(300, hive.position.y-150)),
					},
				}
			}
			hivePointer.beesToCome = append(hive.beesToCome, beesToCome...)

			beesToGo := make([]Bee, hive.beesToRemove)
			for i := 0; i < hive.beesToRemove; i++ {
				beesToGo[i] = Bee{
					position: coordinate{
						x: hive.position.x + 30,
						y: hive.position.y + 30,
					},
				}
			}
			hivePointer.beesToGo = append(hive.beesToGo, beesToGo...)
		}

	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
}

func randomNumberBeetween(max, min float64) float64 {
	//return rand.Intn(max-min) + min
	return min + rand.Float64()*(max-min)
}
func removeBeeNoOrder(bees []Bee, i int) []Bee {
	bees[i] = bees[len(bees)-1]
	return bees[:len(bees)-1]
}

func getRandomNumber() float64 {
	return 0
}

func drawBee(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(0.1, 0.1)
	opChar.GeoM.Translate(x, y)

	return beeImage.SubImage(image.Rect(0, 0, 432, 432)).(*ebiten.Image), opChar
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		hives: []Hive{
			{
				position: coordinate{
					x: 150,
					y: 150,
				},
				beesCount:    1000,
				beesToAdd:    15,
				beesToRemove: 3,
			}},
		mapCenterX:    9,
		mapCenterY:    6,
		worldSpeed:    3,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bees-World")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
