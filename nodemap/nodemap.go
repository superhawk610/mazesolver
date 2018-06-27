package nodemap

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

type NodeMap struct {
	Nodes  [][]*Node
	Width  uint
	Height uint
}

func isWall(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

func FromMaze(maze *image.Image) *NodeMap {
	startX := (*maze).Bounds().Min.X
	startY := (*maze).Bounds().Min.Y
	endX := (*maze).Bounds().Max.X
	endY := (*maze).Bounds().Max.Y
	nm := &NodeMap{
		Width:  uint(endX - startX),
		Height: uint(endY - startY),
	}

	// create nodes at critical points
	for y := startY; y < endY; y++ {
		row := []*Node{}
		for x := startX; x < endX; x++ {

			// left/right vertical walls cannot contain nodes
			if x == 0 || x == endX-1 {
				continue
			}

			// walls cannot contain nodes
			if isWall((*maze).At(x, y)) {
				continue
			}

			// start node exists at x, y = 0
			if y == 0 {
				row = append(
					row,
					NewStartNode(uint(x)),
				)
				continue
			}

			// end node exists at x, y = endY-1
			if y == endY-1 {
				row = append(
					row,
					NewEndNode(uint(x)),
				)
				continue
			}

			// everything else must either be a node or path connector
			input := &Input{
				N: !isWall((*maze).At(x, y-1)),
				E: !isWall((*maze).At(x+1, y)),
				S: !isWall((*maze).At(x, y+1)),
				W: !isWall((*maze).At(x-1, y)),
			}
			row = append(
				row,
				NewNode(uint(x), input),
			)
		}
		nm.Nodes = append(
			nm.Nodes,
			row,
		)
	}

	// connect nodes
	for currentRowIndex, row := range nm.Nodes {
		for nodeIndex, node := range row {
			// for every critical node, connect it to joining critical nodes

			// connect to the next critical node in the current row
			if node.Input.E {
				for searchNodeIndex := nodeIndex + 1; searchNodeIndex < len(row); searchNodeIndex++ {
					testNode := row[searchNodeIndex]
					if testNode.Critical() {
						node.Connect(testNode, Right)
						break
					}
				}
			}

			// check for nodes below and connect first critical match
			if node.Input.S {
			Search:
				for searchRowIndex := currentRowIndex + 1; searchRowIndex < len(nm.Nodes); searchRowIndex++ {
					searchRow := nm.Nodes[searchRowIndex]
					if len(searchRow) == 0 {
						continue
					}

					for _, testNode := range searchRow {
						if testNode.Critical() && testNode.Offset == node.Offset {
							node.Connect(testNode, Down)
							break Search
						}
					}
				}
			}
		}
	}

	return nm
}

func (nm *NodeMap) Visualize() string {
	var sb strings.Builder

	sb.WriteString("NodeMap\n")
	for _, row := range nm.Nodes {
		var (
			node      *Node
			nodeIndex int
		)

		if len(row) > 0 {
			node = row[0]
		}
		for i := uint(0); i < nm.Width; i++ {
			if node != nil && node.Offset == i {

				// display correct node type
				sb.WriteString(node.String())

				// advance node pointer
				nodeIndex++
				if nodeIndex < len(row) {
					node = row[nodeIndex]
				}

			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (nm *NodeMap) String() string {
	var sb strings.Builder

	sb.WriteString("NodeMap {\n")
	for _, row := range nm.Nodes {
		sb.WriteString("  Row { ")
		for i, n := range row {
			terminator := ","
			if i == len(row)-1 {
				terminator = ""
			}
			sb.WriteString(
				fmt.Sprintf(
					"%v%s",
					n.Offset,
					terminator,
				),
			)
		}
		sb.WriteString(" }\n")
	}
	sb.WriteString("}")

	return sb.String()
}
