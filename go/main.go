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
	"github.com/ipreferwater/kafka-bees/kafkabee"
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
	waspImage        *ebiten.Image
	hiveImage       *ebiten.Image
	mplusNormalFont font.Face
)

type cache struct {
	addBees int
	addWasp int
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
}

type Game struct {
	hives      []Hive
	mapCenterX float64
	mapCenterY float64
	worldSpeed float64
	frame      int
}

func (g *Game) Update() error {
	//animate(g)
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

		hiveOpt := &ebiten.DrawImageOptions{}
		hiveOpt.GeoM.Scale(6, 6)
		hiveOpt.GeoM.Translate(hive.position.x, hive.position.y)
		screen.DrawImage(hiveImage, hiveOpt)
		text.Draw(screen, fmt.Sprintf("%d", hive.beesCount), mplusNormalFont, int(hive.position.x)+50, int(hive.position.y)+10, color.RGBA{0, 255, 0, 255})

		for insectType, arrInsects := range hivePointer.insectsToCome {

			if insectType == Bee {

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
						sendDetectionToStream(EuropeanBee, hive.ID, true)
						continue
					}

					getCloserFromHive(beePointer, hive.hiveEntry.x, hive.hiveEntry.y)

				}
			}

			if insectType == Wasp {
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
							//hivePointer.insectsToCome[Wasp] = removeBeeNoOrder(hive.insectsToCome[Wasp], indexWasp)
							//TODO it should not send asianwasp but a simple detection
							waspPointer.waspState = Hunting
							sendDetectionToStream(AsianWasp, hive.ID, true)
							continue
						}
						getCloserFromHiveForHuntingState(waspPointer, hive.hiveEntry.x, hive.hiveEntry.y)
					case Hunting:
						// wait for a bee to kill
					case Leaving:
						// go away with bee victim
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
						sendDetectionToStream(EuropeanBee, hive.ID, false)
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

func getCloserFromHive(insectPointer *Insect, hiveEntryX float64, hiveEntryY float64 ){
	if insectPointer.position.x <= hiveEntryX {
		insectPointer.position.x += 0.5
	}else {
		insectPointer.position.x -= 0.5
	}

	if insectPointer.position.y <= hiveEntryY {
		insectPointer.position.y += randomNumberBeetween(0, 1)
	} else {
		insectPointer.position.y -= randomNumberBeetween(0, 1)
	}
}

func getCloserFromHiveForHuntingState(insectPointer *Insect, hiveEntryX float64, hiveEntryY float64 ){
	if insectPointer.position.x <= randomNumberBeetween(hiveEntryX-10, hiveEntryX+10) {
		insectPointer.position.x += randomNumberBeetween(0, 1)
	}else {
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

func getRandomNumber() float64 {
	return 0
}

func drawBee(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(0.5, 0.5)
	opChar.GeoM.Translate(x, y)

	return beeImage.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), opChar
}

func drawWasp(x, y float64) (*ebiten.Image, *ebiten.DrawImageOptions) {
	opChar := &ebiten.DrawImageOptions{}
	opChar.GeoM.Scale(1.2, 1.2)
	opChar.GeoM.Translate(x, y)

	return waspImage.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), opChar
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fmt.Println("go")
	go kafkabee.Init()
	go kafkabee.InitConsumer()

	g := &Game{
		hives: []Hive{
			{
				ID: 1,
				position: coordinate{
					x: 100,
					y: 100,
				},
				beesCount:    1000,
				beesToAdd:    5,
				beesToRemove: 6,
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
			}},
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
