package main

type viewingPlane struct {
	normal point
	vert   point
	horiz  point
}

func (v *viewingPlane) normalize() {
	v.vert.makeUnit()
	v.horiz.makeUnit()
	v.normal.makeUnit()
}

func (v *viewingPlane) rotate(horiz, forward bool) {
	k := 0.02
	if !forward {
		k = -k
	}
	if horiz {
		v.normal.x += v.horiz.x * k
		v.normal.y += v.horiz.y * k
		v.normal.z += v.horiz.z * k
		v.horiz.x -= v.normal.x * k
		v.horiz.y -= v.normal.y * k
		v.horiz.z -= v.normal.z * k
	} else {
		v.normal.x += v.vert.x * k
		v.normal.y += v.vert.y * k
		v.normal.z += v.vert.z * k
		v.vert.x -= v.normal.x * k
		v.vert.y -= v.normal.y * k
		v.vert.z -= v.normal.z * k
	}
	v.normalize()
}
