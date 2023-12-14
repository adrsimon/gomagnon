package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var imgSapiens, imgNeanderthal, imgBabySapiens, imgBabyNeanderthal, imgPlains, imgForest, imgWater, imgCaves, imgCow, imgMushroom, imgRock, imgWood *ebiten.Image

func init() {
	sapiens, _, err := ebitenutil.NewImageFromFile("assets/images/sapiens.png")
	if err != nil {
		log.Fatal(err)
	}
	imgSapiens = sapiens

	neanderthal, _, err := ebitenutil.NewImageFromFile("assets/images/neanderthal.png")
	if err != nil {
		log.Fatal(err)
	}
	imgNeanderthal = neanderthal

	babyNean, _, err := ebitenutil.NewImageFromFile("assets/images/baby_neanderthal.png")
	if err != nil {
		log.Fatal(err)
	}
	imgBabyNeanderthal = babyNean

	babySapiens, _, err := ebitenutil.NewImageFromFile("assets/images/baby_sapiens.png")
	if err != nil {
		log.Fatal(err)
	}
	imgBabySapiens = babySapiens

	plains, _, err := ebitenutil.NewImageFromFile("assets/images/plains.png")
	if err != nil {
		log.Fatal(err)
	}
	imgPlains = plains

	forest, _, err := ebitenutil.NewImageFromFile("assets/images/forest.png")
	if err != nil {
		log.Fatal(err)
	}
	imgForest = forest

	water, _, err := ebitenutil.NewImageFromFile("assets/images/water.png")
	if err != nil {
		log.Fatal(err)
	}
	imgWater = water

	caves, _, err := ebitenutil.NewImageFromFile("assets/images/caves.png")
	if err != nil {
		log.Fatal(err)
	}
	imgCaves = caves

	cow, _, err := ebitenutil.NewImageFromFile("assets/images/cow.png")
	if err != nil {
		log.Fatal(err)
	}
	imgCow = cow

	mushroom, _, err := ebitenutil.NewImageFromFile("assets/images/mushroom.png")
	if err != nil {
		log.Fatal(err)
	}
	imgMushroom = mushroom

	rock, _, err := ebitenutil.NewImageFromFile("assets/images/rock.png")
	if err != nil {
		log.Fatal(err)
	}
	imgRock = rock

	wood, _, err := ebitenutil.NewImageFromFile("assets/images/wood.png")
	if err != nil {
		log.Fatal(err)
	}
	imgWood = wood
}
