package main

type rubiks_cube struct {
	corners [8]*corner_piece
	edges   [12]*edge_piece
	centers [6]*center_piece
}

type center_piece struct {
	m_side *side
}

type corner_piece struct {
	centers [3]*center_piece
	sides   [3]*side
}

func normalize_number(num float64) float64 {
	if num == 0 {
		return 0
	}
	if num > 0 {
		return 1
	}
	return -1
}

func newCornerPiece(centers [3]*center_piece) *corner_piece {
	sides := [3]*side{}
	for i := 0; i < 3; i++ {
		displacement_vector := point{0, 0, 0}
		for s := 0; s < 3; s++ {
			if s == i {
				continue
			}
			displacement_vector.x += 2.1 * normalize_number(centers[s].m_side.center.x)
			displacement_vector.y += 2.1 * normalize_number(centers[s].m_side.center.y)
			displacement_vector.z += 2.1 * normalize_number(centers[s].m_side.center.z)
		}
		copyOfSide := *centers[i].m_side
		sides[i] = &copyOfSide
		sides[i].displace(displacement_vector)
	}
	return &corner_piece{centers, sides}
}

type edge_piece struct {
	centers [2]*center_piece
	sides   [2]*side
}

func newEdgePiece(centers [2]*center_piece) *edge_piece {
	sides := [2]*side{}
	for i := 0; i < 2; i++ {
		displacement_vector := point{0, 0, 0}
		displacement_vector.x += 2.1 * normalize_number(centers[1-i].m_side.center.x)
		displacement_vector.y += 2.1 * normalize_number(centers[1-i].m_side.center.y)
		displacement_vector.z += 2.1 * normalize_number(centers[1-i].m_side.center.z)
		copyOfSide := *centers[i].m_side
		sides[i] = &copyOfSide
		sides[i].displace(displacement_vector)
	}
	return &edge_piece{centers, sides}
}

func (c *rubiks_cube) getVisibleSides(vector point) (visible, toCleanUp []*side) {
	visible = []*side{}
	toCleanUp = []*side{}
	for _, s := range c.centers {
		if s.m_side.center.dot(vector) > 0 {
			visible = append(visible, s.m_side)
		}
	}
	for _, s := range c.corners {
		for _, side := range s.sides {
			if side.center.dot(vector) > 0 {
				visible = append(visible, side)
			}
		}
	}

	for _, s := range c.edges {
		for _, side := range s.sides {
			if side.center.dot(vector) > 0 {
				visible = append(visible, side)
			}
		}
	}

	return
}

func are_opposites(lhs, rhs point) bool {
	return lhs.x == -rhs.x && lhs.y == -rhs.y && lhs.z == -rhs.z
}

func getEdgePieces(centers [6]*center_piece) [12]*edge_piece {
	edges := [12]*edge_piece{}
	tally := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			if are_opposites(centers[i].m_side.center, centers[j].m_side.center) {
				continue
			}
			edges[tally] = newEdgePiece([2]*center_piece{centers[i], centers[j]})
			tally += 1
		}
	}
	return edges
}

func getCornerPieces(centers [6]*center_piece) [8]*corner_piece {
	corners := [8]*corner_piece{}
	tally := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			for k := j + 1; k < 6; k++ {
				if are_opposites(centers[i].m_side.center, centers[j].m_side.center) || are_opposites(centers[j].m_side.center, centers[k].m_side.center) || are_opposites(centers[i].m_side.center, centers[k].m_side.center) {
					continue
				}
				corners[tally] = newCornerPiece([3]*center_piece{centers[i], centers[j], centers[k]})
				tally += 1
			}
		}
	}
	// print(tally)
	return corners
}

func getCenterPieces() [6]*center_piece {
	side1 := newSide(Green, point{3.2, 0, 0})
	side2 := newSide(Red, point{0, 3.2, 0})
	side3 := newSide(White, point{0, 0, 3.2})
	side4 := newSide(Blue, point{-3.2, 0, 0})
	side5 := newSide(Orange, point{0, -3.2, 0})
	side6 := newSide(Yellow, point{0, 0, -3.2})
	return [6]*center_piece{{&side1}, {&side2}, {&side3}, {&side4}, {&side5}, {&side6}}
}
