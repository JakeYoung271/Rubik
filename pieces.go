package main

import (
	"sort"
)

const (
	SPACING          = 2.12
	DIST_FROM_ORIGIN = 3.2
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////          center_piece           ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

type center_piece struct {
	m_side  *side
	corners [4]*corner_piece
	edges   [4]*edge_piece
}

func getCenterPieces() [6]*center_piece {
	side1 := newSide(Green, point{DIST_FROM_ORIGIN, 0, 0})
	side2 := newSide(Red, point{0, DIST_FROM_ORIGIN, 0})
	side3 := newSide(White, point{0, 0, DIST_FROM_ORIGIN})
	side4 := newSide(Blue, point{-1 * DIST_FROM_ORIGIN, 0, 0})
	side5 := newSide(Orange, point{0, -1 * DIST_FROM_ORIGIN, 0})
	side6 := newSide(Yellow, point{0, 0, -1 * DIST_FROM_ORIGIN})
	return [6]*center_piece{{&side1, [4]*corner_piece{}, [4]*edge_piece{}}, {&side2, [4]*corner_piece{}, [4]*edge_piece{}}, {&side3, [4]*corner_piece{}, [4]*edge_piece{}}, {&side4, [4]*corner_piece{}, [4]*edge_piece{}}, {&side5, [4]*corner_piece{}, [4]*edge_piece{}}, {&side6, [4]*corner_piece{}, [4]*edge_piece{}}}
}

func (c *center_piece) claim_corners_and_edges(corners [8]*corner_piece, edges [12]*edge_piece) {
	tally := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 3; j++ {
			if corners[i].sides[j].center.equals(c.m_side.center) {
				c.corners[tally] = corners[i]
				tally += 1
				j = 3
			}
		}
	}
	tally = 0
	for i := 0; i < 12; i++ {
		if edges[i].sides[0].center.equals(c.m_side.center) || edges[i].sides[1].center.equals(c.m_side.center) {
			c.edges[tally] = edges[i]
			tally += 1
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////          corner_piece           ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

type corner_piece struct {
	centers [3]*center_piece
	sides   [3]*side
}

func newCornerPiece(centers [3]*center_piece) *corner_piece {
	sides := [3]*side{}
	for i := 0; i < 3; i++ {
		displacement_vector := point{0, 0, 0}
		for s := 0; s < 3; s++ {
			if s == i {
				continue
			}
			displacement_vector.x += SPACING * normalize_number(centers[s].m_side.center.x)
			displacement_vector.y += SPACING * normalize_number(centers[s].m_side.center.y)
			displacement_vector.z += SPACING * normalize_number(centers[s].m_side.center.z)
		}
		copyOfSide := *centers[i].m_side
		sides[i] = &copyOfSide
		sides[i].displace(displacement_vector)
	}
	return &corner_piece{centers, sides}
}

func getCornerPieces(centers [6]*center_piece) [8]*corner_piece {
	corners := [8]*corner_piece{}
	side_corner_counts := [6]int{}
	tally := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			for k := j + 1; k < 6; k++ {
				if are_opposites(centers[i].m_side.center, centers[j].m_side.center) || are_opposites(centers[j].m_side.center, centers[k].m_side.center) || are_opposites(centers[i].m_side.center, centers[k].m_side.center) {
					continue
				}
				corners[tally] = newCornerPiece([3]*center_piece{centers[i], centers[j], centers[k]})
				centers[i].corners[side_corner_counts[i]] = corners[tally]
				centers[j].corners[side_corner_counts[j]] = corners[tally]
				centers[k].corners[side_corner_counts[k]] = corners[tally]
				side_corner_counts[i] += 1
				side_corner_counts[j] += 1
				side_corner_counts[k] += 1
				tally += 1
			}
		}
	}
	return corners
}

func (c *corner_piece) rotate(axis int, forward bool) {
	for i := 0; i < 3; i++ {
		c.sides[i].rotate(axis, forward)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////            edge_piece           ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

type edge_piece struct {
	centers [2]*center_piece
	sides   [2]*side
}

func newEdgePiece(centers [2]*center_piece) *edge_piece {
	sides := [2]*side{}
	for i := 0; i < 2; i++ {
		displacement_vector := point{0, 0, 0}
		displacement_vector.x += SPACING * normalize_number(centers[1-i].m_side.center.x)
		displacement_vector.y += SPACING * normalize_number(centers[1-i].m_side.center.y)
		displacement_vector.z += SPACING * normalize_number(centers[1-i].m_side.center.z)
		copyOfSide := *centers[i].m_side
		sides[i] = &copyOfSide
		sides[i].displace(displacement_vector)
	}
	return &edge_piece{centers, sides}
}

func getEdgePieces(centers [6]*center_piece) [12]*edge_piece {
	edges := [12]*edge_piece{}
	side_edge_counts := [6]int{}
	tally := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			if are_opposites(centers[i].m_side.center, centers[j].m_side.center) {
				continue
			}
			edges[tally] = newEdgePiece([2]*center_piece{centers[i], centers[j]})
			centers[i].edges[side_edge_counts[i]] = edges[tally]
			centers[j].edges[side_edge_counts[j]] = edges[tally]
			side_edge_counts[i] += 1
			side_edge_counts[j] += 1
			tally += 1
		}
	}
	return edges
}

func (e *edge_piece) rotate(axis int, forward bool) {
	for i := 0; i < 2; i++ {
		e.sides[i].rotate(axis, forward)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////           rubiks_cube           ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

type rubiks_cube struct {
	corners [8]*corner_piece
	edges   [12]*edge_piece
	centers [6]*center_piece
}

func (c *rubiks_cube) rotate_side(color color, axis int, forward bool) {
	var spec_center *center_piece
	for _, center := range c.centers {
		if center.m_side.side_color == color {
			spec_center = center
			break
		}
	}
	for _, corner := range spec_center.corners {
		corner.rotate(axis, forward)
	}
	for _, edge := range spec_center.edges {
		edge.rotate(axis, forward)
	}
	spec_center.m_side.rotate(axis, forward)
}

type by_proximity []*side

func (a by_proximity) Len() int      { return len(a) }
func (a by_proximity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a by_proximity) Less(i, j int) bool {
	return sq_distance(scaled(view.normal, 10), average(a[i].edges)) > sq_distance(scaled(view.normal, 10), average(a[j].edges))
}

func (c *rubiks_cube) getVisibleSides(vector point) (visible, toCleanUp []*side) {
	visible = []*side{}
	toCleanUp = []*side{}
	for _, center := range c.centers {
		if center.m_side.center.dot(vector) > 0 {
			visible = append(visible, center.m_side)
		}
	}
	for _, edge := range c.edges {
		for _, side := range edge.sides {
			if side.center.dot(vector) > 0 {
				visible = append(visible, side)
			}
		}
	}
	for _, corner := range c.corners {
		for _, side := range corner.sides {
			if side.center.dot(vector) > 0 {
				visible = append(visible, side)
			}
		}
	}

	sort.Sort(by_proximity(visible))
	return
}
