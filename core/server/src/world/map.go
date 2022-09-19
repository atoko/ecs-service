package world

type TileType byte
const (
	OPEN TileType = iota
	WALL
)

type Map struct {
	Tiles []TileType
	Width int
	Height int
}

func (m *Map) xy_idx(x int, y int) int {

	return (y * m.Width) + x
}

func NewBorderedMap(width int, height int) *Map {
	tiles := make([]TileType, (width*height)+1)
	m := &Map{
		Tiles: tiles,
		Width: width,
		Height: height,
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 || x == width - 1 {
				tiles[m.xy_idx(x, y)] = WALL
			} else if y == 0 || y == height - 1 {
				tiles[m.xy_idx(x, y)] = WALL
			} else {
				tiles[m.xy_idx(x, y)] = OPEN
			}
		}
	}

	return m
}