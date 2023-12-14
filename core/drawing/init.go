package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var imgSapiens, imgNeanderthal, imgBabySapiens, imgBabyNeanderthal, imgPlains, imgForest, imgWater, imgCaves, imgCow, imgMushroom, imgRock, imgWood, imgMammoth *ebiten.Image

func init() {
	sapiens, _, err := ebitenutil.NewImageFromFile("assets/textures/sapiens.png")
	if err != nil {
		log.Fatal(err)
	}
	imgSapiens = sapiens

	neanderthal, _, err := ebitenutil.NewImageFromFile("assets/textures/neanderthal.png")
	if err != nil {
		log.Fatal(err)
	}
	imgNeanderthal = neanderthal

	babyNean, _, err := ebitenutil.NewImageFromFile("assets/textures/baby_neanderthal.png")
	if err != nil {
		log.Fatal(err)
	}
	imgBabyNeanderthal = babyNean

	babySapiens, _, err := ebitenutil.NewImageFromFile("assets/textures/baby_sapiens.png")
	if err != nil {
		log.Fatal(err)
	}
	imgBabySapiens = babySapiens

	plains, _, err := ebitenutil.NewImageFromFile("assets/textures/plains.png")
	if err != nil {
		log.Fatal(err)
	}
	imgPlains = plains

	forest, _, err := ebitenutil.NewImageFromFile("assets/textures/forest.png")
	if err != nil {
		log.Fatal(err)
	}
	imgForest = forest

	water, _, err := ebitenutil.NewImageFromFile("assets/textures/water.png")
	if err != nil {
		log.Fatal(err)
	}
	imgWater = water

	caves, _, err := ebitenutil.NewImageFromFile("assets/textures/caves.png")
	if err != nil {
		log.Fatal(err)
	}
	imgCaves = caves

	cow, _, err := ebitenutil.NewImageFromFile("assets/textures/cow.png")
	if err != nil {
		log.Fatal(err)
	}
	imgCow = cow

	mushroom, _, err := ebitenutil.NewImageFromFile("assets/textures/mushroom.png")
	if err != nil {
		log.Fatal(err)
	}
	imgMushroom = mushroom

	rock, _, err := ebitenutil.NewImageFromFile("assets/textures/rock.png")
	if err != nil {
		log.Fatal(err)
	}
	imgRock = rock

	wood, _, err := ebitenutil.NewImageFromFile("assets/textures/wood.png")
	if err != nil {
		log.Fatal(err)
	}
	imgWood = wood

	mammoth, _, err := ebitenutil.NewImageFromFile("assets/textures/mammoth.png")
	if err != nil {
		log.Fatal(err)
	}
	imgMammoth = mammoth
}
