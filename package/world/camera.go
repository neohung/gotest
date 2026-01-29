package world

type Camera struct {
	X, Y int
	W, H int
}

func (c *Camera) Follow(px, py int, worldW, worldH int) {
	// 依據player x y 置中
	c.X = px - c.W/2
	c.Y = py - c.H/2
	// clamp
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	if c.X+c.W > worldW {
		c.X = worldW - c.W
	}
	if c.Y+c.H > worldH {
		c.Y = worldH - c.H
	}
}
