package drawing

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var imgSapiens, imgNeanderthal, imgBabySapiens, imgBabyNeanderthal, imgPlains, imgForest, imgWater, imgDeepWater, imgCaves, imgCow, imgMushroom, imgRock, imgWood, imgHutSapiens, imgHutNeanderthal, imgHutAbandonned *ebiten.Image

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

	deepwater, _, err := ebitenutil.NewImageFromFile("assets/textures/deepwater.png")
	if err != nil {
		log.Fatal(err)
	}
	imgDeepWater = deepwater

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

	hutSapiens, _, err := ebitenutil.NewImageFromFile("assets/textures/hutSapiens.png")
	if err != nil {
		log.Fatal(err)
	}
	imgHutSapiens = hutSapiens

	hutNeanderthal, _, err := ebitenutil.NewImageFromFile("assets/textures/hutNeanderthal.png")
	if err != nil {
		log.Fatal(err)
	}
	imgHutNeanderthal = hutNeanderthal

	hutAbandonned, _, err := ebitenutil.NewImageFromFile("assets/textures/hutAbandonned.png")
	if err != nil {
		log.Fatal(err)
	}
	imgHutAbandonned = hutAbandonned
}
