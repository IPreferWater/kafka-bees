package gui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	beeImage        *ebiten.Image
	waspImage       *ebiten.Image
	hiveImage       *ebiten.Image
	mplusNormalFont font.Face
	bgPixel         background
)

const (
	screenWidth  = 1024
	screenHeight = 600
)

type background struct {
	img        *ebiten.Image
	tree       pixelImage
	grass      pixelImage
	littleTree pixelImage
}

type pixelImage struct {
	x1, y1, x2, y2 int
}

type Game struct {
	hives      []Hive
	mapCenterX float64
	mapCenterY float64
	worldSpeed float64
	frame      int
}

func StartEbiten() {
	g := &Game{
		hives: []Hive{
			{
				ID: 1,
				position: coordinate{
					x: 100,
					y: 100,
				},
				beesCount:     1000,
				beesToAdd:     20,
				beesToRemove:  2,
				waspsCount:    0,
				waspsToAdd:    1,
				waspsToRemove: 1,
				hiveEntry: coordinate{
					x: 152,
					y: 190,
				},
				hiveExit: coordinate{
					x: 180,
					y: 190,
				},
				insectsToCome: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
				insectsToGo: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
			},
			{
				ID: 2,
				position: coordinate{
					x: 400,
					y: 100,
				},
				beesCount:     1000,
				beesToAdd:     20,
				beesToRemove:  2,
				waspsCount:    0,
				waspsToAdd:    1,
				waspsToRemove: 1,
				hiveEntry: coordinate{
					x: 452,
					y: 190,
				},
				hiveExit: coordinate{
					x: 480,
					y: 190,
				},
				insectsToCome: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
				insectsToGo: map[InsecType][]Insect{
					Bee:  make([]Insect, 0),
					Wasp: make([]Insect, 0),
				},
			},
		},
		mapCenterX: 9,
		mapCenterY: 6,
		worldSpeed: 3,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bees-World")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func getCloserFromHive(insectPointer *Insect, hiveEntryX float64, hiveEntryY float64) {
	if insectPointer.position.x <= hiveEntryX {
		insectPointer.position.x += 0.5
	} else {
		insectPointer.position.x -= 0.5
	}

	if insectPointer.position.y <= hiveEntryY {
		insectPointer.position.y += randomNumberBeetween(0, 1)
	} else {
		insectPointer.position.y -= randomNumberBeetween(0, 1)
	}
}

func getCloserFromHiveForHuntingState(insectPointer *Insect, hiveEntryX float64, hiveEntryY float64) {
	if insectPointer.position.x <= randomNumberBeetween(hiveEntryX-10, hiveEntryX+10) {
		insectPointer.position.x += randomNumberBeetween(0, 1)
	} else {
		insectPointer.position.x -= randomNumberBeetween(0, 1)
	}

	if insectPointer.position.y <= hiveEntryY {
		insectPointer.position.y += 0.5
	} else {
		insectPointer.position.y -= 0.5
	}
}

func randomNumberBeetween(max, min float64) float64 {
	//return rand.Intn(max-min) + min
	return min + rand.Float64()*(max-min)
}
func removeBeeNoOrder(bees []Insect, i int) []Insect {
	bees[i] = bees[len(bees)-1]
	return bees[:len(bees)-1]
}

func drawBee(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(0.5, 0.5)
	opChar.GeoM.Translate(x, y)

	return beeImage.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), opChar
}

func drawWasp(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(0.8, 0.8)
	opChar.GeoM.Translate(x, y)

	return waspImage.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), opChar
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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

func getNumberOfTilesToDraw(w, h, tileSize int) (int, int) {
	return divideAndRoundUp(w, tileSize), divideAndRoundUp(h, tileSize)
}

func divideAndRoundUp(a, b int) int {
	res := a / b
	if a%b > 0 {
		res++
	}
	return res
}

func drawBgTiles(x, y float64, pxI pixelImage) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(2, 2)
	opChar.GeoM.Translate(x, y)
	//return bgPixel.img.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), opChar
	return bgPixel.img.SubImage(image.Rect(pxI.x1, pxI.y1, pxI.x2, pxI.y2)).(*ebiten.Image), opChar
}

func (g *Game) Update() error {
	//animate(g)
	return nil
}

func init() {
	beeEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/bee.png")
	if err != nil {
		log.Fatal(err)
	}

	beeImage = beeEbitenImage

	waspEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/wasp.png")
	if err != nil {
		log.Fatal(err)
	}

	waspImage = waspEbitenImage

	hiveEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/hive.png")
	if err != nil {
		log.Fatal(err)
	}
	hiveImage = hiveEbitenImage

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	backgroundEbitenImage, _, err := ebitenutil.NewImageFromFile("./res/tiles/trees-and-bushes.png")

	bgPixel = background{
		img:        backgroundEbitenImage,
		tree:       pixelImage{0, 0, 64, 64},
		grass:      pixelImage{0, 96, 32, 128},
		littleTree: pixelImage{0, 64, 32, 96},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.frame++

	drawBackGround(screen)

	secondUpdate := false
	if g.frame%60 == 0 {
		secondUpdate = true
		g.frame = 0
	}

	for indexHive, hive := range g.hives {
		hivePointer := &g.hives[indexHive]

		hiveOpt := &ebiten.DrawImageOptions{}
		hiveOpt.GeoM.Scale(6, 6)
		hiveOpt.GeoM.Translate(hive.position.x, hive.position.y)
		screen.DrawImage(hiveImage, hiveOpt)
		text.Draw(screen, fmt.Sprintf("%d", hive.beesCount), mplusNormalFont, int(hive.position.x)+50, int(hive.position.y)+10, color.RGBA{0, 255, 0, 255})

		for insectType, arrInsects := range hivePointer.insectsToCome {

			switch insectType {
			case Bee:
				for indexBee, bee := range arrInsects {

					//arrInsects is not updated when we removed an insect, if the indexBee is larger than the arrSize, it means he finished
					if indexBee >= len(hivePointer.insectsToCome[Bee]) {
						break
					}
					beePointer := &hive.insectsToCome[Bee][indexBee]
					beeImg, beeOpts := drawBee(bee.position.x, bee.position.y)
					screen.DrawImage(beeImg, beeOpts)

					if bee.position.x >= hive.hiveEntry.x && bee.position.y >= hive.hiveEntry.y {
						hivePointer.insectsToCome[Bee] = removeBeeNoOrder(hive.insectsToCome[Bee], indexBee)
						hivePointer.beesCount++
						sendDetectionToStream(Bee, hive.ID, true)
						continue
					}

					getCloserFromHive(beePointer, hive.hiveEntry.x, hive.hiveEntry.y)

				}
			case Wasp:
				for indexWasp, wasp := range arrInsects {
					//arrInsects is not updated when we removed an insect, if the indexBee is larger than the arrSize, it means he finished
					if indexWasp >= len(hivePointer.insectsToCome[Wasp]) {
						break
					}
					waspPointer := &hive.insectsToCome[Wasp][indexWasp]

					switch waspPointer.waspState {
					case Approching:
						waspImg, beeOpts := drawWasp(wasp.position.x, wasp.position.y)
						screen.DrawImage(waspImg, beeOpts)

						if wasp.position.x >= hive.hiveEntry.x && wasp.position.y >= hive.hiveEntry.y {

							//TODO it should not send asianwasp but a simple detection
							waspPointer.waspState = Hunting
							sendDetectionToStream(Wasp, hive.ID, true)
							continue
						}
						getCloserFromHiveForHuntingState(waspPointer, hive.hiveEntry.x, hive.hiveEntry.y)

					case Hunting:
						xPosition := randomNumberBeetween(wasp.position.x-1, wasp.position.x+1)
						yPosition := randomNumberBeetween(wasp.position.y-1, wasp.position.y+1)
						waspImg, beeOpts := drawWasp(xPosition, yPosition)
						screen.DrawImage(waspImg, beeOpts)

						indexBee := getInsectInRange(hivePointer.insectsToCome[Bee], xPosition, yPosition, 16)

						if indexBee == -1 {
							break
						}

						waspPointer.waspState = Leaving
						hivePointer.insectsToCome[Bee] = removeBeeNoOrder(hive.insectsToCome[Bee], indexBee)
						sendDetectionToStream(Wasp, hive.ID, false)

					case Leaving:

						waspImg, beeOpts := drawWasp(wasp.position.x, wasp.position.y)
						screen.DrawImage(waspImg, beeOpts)
						waspPointer.position.x += randomNumberBeetween(-1, 1)
						waspPointer.position.y += randomNumberBeetween(0.1, 1)

						beeVictim, beeVictimOpts := drawBee(wasp.position.x+8, wasp.position.y+12)
						screen.DrawImage(beeVictim, beeVictimOpts)

						if waspPointer.position.y >= hivePointer.hiveEntry.y+200 {
							hivePointer.insectsToCome[Wasp] = removeBeeNoOrder(hive.insectsToCome[Wasp], indexWasp)
						}

					default:
						//do nothing
					}

				}
			}

		}

		for insectType, arrInsects := range hivePointer.insectsToGo {
			if insectType == Bee {
				for indexBee, bee := range arrInsects {
					//TODO can I do it better ?
					//arrInsects is not updated when we removed an insect, if the indexBee is larger than the arrSize, it means he finished
					if indexBee >= len(hivePointer.insectsToGo[Bee]) {
						break
					}
					beePointer := &hive.insectsToGo[Bee][indexBee]
					beeImg, beeOpts := drawBee(bee.position.x, bee.position.y)
					screen.DrawImage(beeImg, beeOpts)

					if bee.position.x >= hive.hiveExit.x+200 {
						hivePointer.insectsToGo[Bee] = removeBeeNoOrder(hive.insectsToGo[Bee], indexBee)
						//hivePointer.beesCount--
						sendDetectionToStream(Bee, hive.ID, false)
						continue
					}

					beePointer.position.x += randomNumberBeetween(0.5, 1.5)
					beePointer.position.y += randomNumberBeetween(-1.5, 1.5)

				}
			}
		}

		if secondUpdate {
			//bees
			beesToCome := make([]Insect, hive.beesToAdd)
			for i := 0; i < hive.beesToAdd; i++ {
				beesToCome[i] = Insect{
					position: coordinate{
						//x: hive.position.x - 100,
						x: hive.position.x - randomNumberBeetween(150, 80),
						y: float64(randomNumberBeetween(300, hive.position.y-150)),
					},
				}
			}
			hivePointer.insectsToCome[Bee] = append(hivePointer.insectsToCome[Bee], beesToCome...)

			waspsToCome := make([]Insect, hive.waspsToAdd)
			for i := 0; i < hive.waspsToAdd; i++ {
				waspsToCome[i] = Insect{
					position: coordinate{
						//x: hive.position.x - 100,
						x: hive.hiveEntry.x - randomNumberBeetween(-50, 50),
						y: hive.hiveEntry.y - randomNumberBeetween(200, 220),
					},
					waspState: Approching,
				}
			}
			hivePointer.insectsToCome[Wasp] = append(hivePointer.insectsToCome[Wasp], waspsToCome...)

			beesToGo := make([]Insect, hive.beesToRemove)
			for i := 0; i < hive.beesToRemove; i++ {
				//we count them at immediatly left from the hive
				hivePointer.beesCount--
				beesToGo[i] = Insect{
					position: coordinate{
						x: hive.hiveExit.x,
						y: hive.hiveExit.y,
					},
				}
			}
			hivePointer.insectsToGo[Bee] = append(hivePointer.insectsToGo[Bee], beesToGo...)
		}

	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
}

// return -1 if nothing found
func getInsectInRange(arr []Insect, x, y, rangeInsect float64) int {

	for indexInsect, i := range arr {

		// is be in range of the wasp ?
		if i.position.x+rangeInsect >= x && i.position.x-rangeInsect <= x && i.position.y+rangeInsect >= y && i.position.y-rangeInsect <= y {
			//waspPointer.waspState = Leaving
			// this bee was killed by this wasp
			//hivePointer.insectsToCome[Bee] = removeBeeNoOrder(hive.insectsToCome[Bee], indexBee)
			return indexInsect
		}
	}
	return -1
}
