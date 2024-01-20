package components

import (
	_ "image/png"
	"log"

	"github.com/PurityLake/thatsmyspot/data"
	"github.com/PurityLake/thatsmyspot/mapreader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TiledMap struct {
	Width, Height int
	TileW, TileH  int
	Tiles         []int
	TempImage     *ebiten.Image
}

func NewTiledMap(imgFilename, mapFilename, tilesetFilename string) (*TiledMap, error) {
	tiledMap, _, err := ebitenutil.NewImageFromFile(imgFilename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	mapObj, err := mapreader.ReadJson(mapFilename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	tilesetObj, err := mapreader.ReadJson(tilesetFilename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	mapData := mapreader.ParseMapData(mapObj)
	tilesetData := mapreader.ParseTilesetData(tilesetObj)

	propertyList, ok := tilesetData["properties"].Value.([]data.Property)
	if !ok {
		log.Fatal("Could not parse tileset properties")
	}
	tiles := make([][]data.Property, 0)
	for _, prop := range propertyList {
		props, ok := prop.Value.(map[string]data.Property)
		tileProps := make([]data.Property, 0)
		if !ok {
			log.Fatal("Could not parse tile properties")
		}
		for _, p := range props {
			tileProps = append(tileProps, p)
		}
		tiles = append(tiles, tileProps)
	}
	var tilesTypes []int
	for _, layer := range mapData["layers"].Value.([]data.Property) {
		for _, tileIndex := range layer.Value.([]int) {
			if tileIndex == 0 {
				tilesTypes = append(tilesTypes, -1)
				continue
			}
			tile := tiles[tileIndex-1]
			tileType := int(tile[0].Value.(float64))
			tilesTypes = append(tilesTypes, tileType)
		}
	}
	bounds := tiledMap.Bounds()
	return &TiledMap{
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		TileW:     40,
		TileH:     40,
		TempImage: tiledMap,
		Tiles:     tilesTypes,
	}, nil
}
