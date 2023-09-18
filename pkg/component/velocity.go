package component

type Velocity struct {
	VX float64
	VY float64
}

func (v *Velocity) ClampVelocity(limit float64) {
	if v.VX > limit {
		v.VX = limit
	}

	if v.VX < -limit {
		v.VX = -limit
	}
	if v.VY > limit {
		v.VY = limit
	}
	if v.VY < -limit {
		v.VY = -limit
	}
}

func (v *Velocity) DiminishVelocity(friction float64) {
	if v.VX < 0 {
		v.VX += friction
	}
	if v.VX > 0 {
		v.VX -= friction
	}
	if v.VY < 0 {
		v.VY += friction
	}
	if v.VY > 0 {
		v.VY -= friction
	}

	if v.VY == friction {
		v.VY = 0
	}
	if v.VX == friction {
		v.VX = 0
	}
}
