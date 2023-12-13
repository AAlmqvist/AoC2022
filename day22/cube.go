package main

import "fmt"

var cubeEdgeMapping [6]edgeMap = [6]edgeMap{
	// side 0
	{
		sideNbr:    [4]int{2, 3, 4, 1},
		edgeNumber: [4]int{2, 3, 0, 1},
		mirrored:   [4]bool{false, false, false, false},
	},
	// side 1
	{
		sideNbr:    [4]int{2, 0, 4, 5},
		edgeNumber: [4]int{3, 3, 3, 3},
		mirrored:   [4]bool{true, false, false, true},
	},
	// side 2
	{
		sideNbr:    [4]int{5, 3, 0, 1},
		edgeNumber: [4]int{2, 0, 0, 0},
		mirrored:   [4]bool{false, false, false, true},
	},
	// side 3
	{
		sideNbr:    [4]int{2, 5, 4, 0},
		edgeNumber: [4]int{1, 1, 1, 1},
		mirrored:   [4]bool{false, true, true, false},
	},
	// side 4
	{
		sideNbr:    [4]int{0, 3, 5, 1},
		edgeNumber: [4]int{2, 2, 0, 2},
		mirrored:   [4]bool{false, true, false, false},
	},
	// side 5
	{
		sideNbr:    [4]int{4, 3, 2, 1},
		edgeNumber: [4]int{2, 1, 0, 3},
		mirrored:   [4]bool{false, true, false, true},
	},
}

type edgeMap struct {
	sideNbr    [4]int
	edgeNumber [4]int
	mirrored   [4]bool
}

func (em *edgeMap) getNewPos(dir, other, size int) (int, *Pos) {
	newSide := em.sideNbr[dir]
	newDir := (em.edgeNumber[dir] + 2) % 4
	newPos := &Pos{dir: newDir}
	switch newDir {
	case right:
		newPos.x = 0
		newPos.y = other
		if em.mirrored[dir] {
			newPos.y = size - 1 - other
		}
	case down:
		newPos.x = other
		newPos.y = 0
		if em.mirrored[dir] {
			newPos.x = size - 1 - other
		}

	case left:
		newPos.x = size - 1
		newPos.y = other
		if em.mirrored[dir] {
			newPos.y = size - 1 - other
		}

	case up:
		newPos.x = other
		newPos.y = size - 1
		if em.mirrored[dir] {
			newPos.x = size - 1 - other
		}
	}

	return newSide, newPos
}

type side struct {
	size int
	grid [][]int
	// Point (y,x) in cove that the (0,0) relative position on the side
	// in the original rotation corresponds to
	translation []int
	// Rotation compared to original layout
	rotation int
}

// Rotate the grid 90 degrees clockwise
func (s *side) rotate() {
	newGrid := [][]int{}
	for y := range s.grid {
		row := []int{}
		for x := range s.grid[y] {
			row = append(row, s.grid[s.size-1-x][y])
		}
		newGrid = append(newGrid, row)
	}
	s.grid = newGrid
}

// A representation of a cube formed by 6 square sides in a pattern
// as follows:
//			1
// 			|
//	5---4---0---2
//			|
//			3
//
type cube struct {
	p     *Pos
	curr  int // index of side p is currently on
	sides [6]*side
}

func (c *cube) move(steps int) {
	for i := 0; i < steps; i++ {
		currSide := c.sides[c.curr]
		hitRock := false
		switch c.p.dir {
		case right:
			// Still on same side
			if c.p.x < currSide.size-1 {
				if currSide.grid[c.p.y][c.p.x+1] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.p.x += 1
			} else {
				nextSide, newPos := cubeEdgeMapping[c.curr].getNewPos(c.p.dir, c.p.y, currSide.size)
				if c.sides[nextSide].grid[newPos.y][newPos.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.curr = nextSide
				c.p = newPos
			}
		case down:
			// Still on same side
			if c.p.y < currSide.size-1 {
				if currSide.grid[c.p.y+1][c.p.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.p.y += 1
			} else {
				nextSide, newPos := cubeEdgeMapping[c.curr].getNewPos(c.p.dir, c.p.x, currSide.size)
				if c.sides[nextSide].grid[newPos.y][newPos.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.curr = nextSide
				c.p = newPos
			}
		case left:
			// Still on same side
			if c.p.x > 0 {
				if currSide.grid[c.p.y][c.p.x-1] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.p.x -= 1
			} else {
				nextSide, newPos := cubeEdgeMapping[c.curr].getNewPos(c.p.dir, c.p.y, currSide.size)
				if c.sides[nextSide].grid[newPos.y][newPos.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.curr = nextSide
				c.p = newPos
			}
		case up:
			// Still on same side
			if c.p.y > 0 {
				if currSide.grid[c.p.y-1][c.p.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.p.y -= 1
			} else {
				nextSide, newPos := cubeEdgeMapping[c.curr].getNewPos(c.p.dir, c.p.x, currSide.size)
				if c.sides[nextSide].grid[newPos.y][newPos.x] == 1 {
					// Hit a rock
					hitRock = true
					break
				}
				c.curr = nextSide
				c.p = newPos
			}
		}
		if hitRock {
			break
		}
	}
}

// Use the current side the position is on to rotate and then
// translate into the positions corresponding position in the cove
// and then use the score from the first part
func (c *cube) orig_pos() *Pos {
	o := c.p.copy()
	// Find original rotaton and relative position
	for i := 0; i < c.sides[c.curr].rotation; i++ {
		// x corresponds to col, y to row
		x := o.x
		y := o.y
		o.x = y
		o.y = c.sides[c.curr].size - 1 - x
		// rotate direction
		o.dir = wrapAround(o.dir-1, 4)
	}
	// Translate back to original position as well
	o.y += c.sides[c.curr].translation[0]
	o.x += c.sides[c.curr].translation[1]

	// We're in original place and rotation, return the score
	return o
}

func cubeFromMap(cove [][]int, start *Pos) *cube {
	cube := &cube{p: &Pos{}}
	for i := 0; i < 6; i++ {
		cube.sides[i] = &side{}
	}
	// Map the sides of the cube given the cove-layout and starting position.
	// The starting position should always be on side 0 in the cube.
	sideSize := findGridSize(cove)
	// Make a condensed map with only a single integer (0/1) depending on
	// if the grid is a side or not.
	sideGrid := [][]int{}
	visited := [][]bool{}
	for i := 0; i < len(cove)/sideSize; i++ {
		visitRow := []bool{}
		gridRow := []int{}
		for j := 0; j < len(cove[0])/sideSize; j++ {
			if cove[i*sideSize][j*sideSize] < 0 {
				gridRow = append(gridRow, -1)
				visitRow = append(visitRow, false)
				continue
			}
			gridRow = append(gridRow, 0)
			visitRow = append(visitRow, false)
		}
		sideGrid = append(sideGrid, gridRow)
		visited = append(visited, visitRow)
	}

	// find index of grid corresponding to 0 in the map
	g_row, g_col := start.y/sideSize, start.x/sideSize
	// nodes in the search are (sideNbr, grid_row, grid_col, rotation)
	nodes := [][]int{}
	nodes = append(nodes, []int{0, g_row, g_col, 0})

	for len(nodes) > 0 {
		// Get all the meta info
		curr := nodes[0]
		nodes = nodes[1:]
		ind, r, c, rot := curr[0], curr[1], curr[2], curr[3]
		visited[r][c] = true
		// Create the side and fill and rotate its grid
		currSide := &side{
			size:        sideSize,
			translation: []int{r * sideSize, c * sideSize},
			rotation:    rot,
		}
		innerGrid := [][]int{}
		for i := 0; i < sideSize; i++ {
			innerGrid = append(innerGrid, cove[r*sideSize+i][c*sideSize:(c+1)*sideSize])
		}
		currSide.grid = innerGrid
		for i := 0; i < rot; i++ {
			currSide.rotate()
		}
		// Point to the side from the slice on the cube
		cube.sides[ind] = currSide

		// Look for neighboring sides and add them to the nodes to explore
		for dir := right; dir <= up; dir++ {
			var dr, dc int
			switch dir {
			case right:
				dc = 1
			case down:
				dr = 1
			case left:
				dc = -1
			case up:
				dr = -1
			}
			// Check we're inbounds
			if r+dr < 0 || r+dr >= len(sideGrid) || c+dc < 0 || c+dc >= len(sideGrid[0]) {
				continue
			}
			if sideGrid[r+dr][c+dc] >= 0 && !visited[r+dr][c+dc] {
				relSide := (rot + dir) % 4
				neighbor := cubeEdgeMapping[ind].sideNbr[relSide]
				relRot := wrapAround(cubeEdgeMapping[ind].edgeNumber[relSide]-(dir+2)%4, 4)
				nodes = append(
					nodes,
					[]int{neighbor,
						r + dr,
						c + dc,
						relRot},
				)
				sideGrid[r+dr][c+dc] = neighbor
			}
		}
	}
	if DEBUG {
		for si, side := range cube.sides {
			fmt.Printf("%d => %v %d\n", si, side.translation, side.rotation)
		}
	}
	return cube
}

// The ratio of the sides being mapped onto 2D is either (4,3) or (3,4)
func findGridSize(cove [][]int) int {
	if len(cove)/4 == len(cove[0])/3 {
		return len(cove) / 4
	}
	return len(cove) / 3
}
