package tileproperties

type TileProperties struct {
	Properties map[string]interface{}
}

// NewTileProperties () returns a new tile properties interface
func NewTileProperties() *TileProperties {

	return &TileProperties{}
}
