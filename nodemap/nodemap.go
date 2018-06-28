package nodemap

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

type NodeMap struct {
	Nodes  [][]*Node
	Width  int
	Height int
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
		Width:  endX - startX,
		Height: endY - startY,
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
					NewStartNode(x, y),
				)
				continue
			}

			// end node exists at x, y = endY-1
			if y == endY-1 {
				row = append(
					row,
					NewEndNode(x, y),
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
				NewNode(x, y, input),
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
			if !node.Critical() {
				continue
			}

			// connect to the next critical node in the current row
			if node.Input.E {
				for searchNodeIndex := nodeIndex + 1; searchNodeIndex < len(row); searchNodeIndex++ {
					testNode := row[searchNodeIndex]
					if testNode.Critical() {
						node.Connect(testNode, Right, testNode.Offset.X-node.Offset.X)
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
						if testNode.Critical() && testNode.Offset.X == node.Offset.X {
							node.Connect(testNode, Down, searchRowIndex-currentRowIndex)
							break Search
						}
					}
				}
			}
		}
	}

	return nm
}

func (nm *NodeMap) Stat() {
	var totalNodes, criticalNodes, deadNodes int

	for _, row := range nm.Nodes {
		for _, node := range row {
			totalNodes++
			if node.Critical() {
				criticalNodes++
			}
			if node.DeadEnd() {
				deadNodes++
			}
		}
	}

	fmt.Printf(
		"Nodes created: %v\n"+
			"Critical nodes: %v\n"+
			"Dead end nodes: %v\n",
		totalNodes,
		criticalNodes,
		deadNodes,
	)
}

func (nm *NodeMap) Visualize() string {
	var sb strings.Builder

	sb.WriteString("NodeMap\n ")
	for i := 0; i < nm.Width; i++ {
		sb.WriteString(fmt.Sprintf("%v", i%10))
	}
	sb.WriteString("\n")
	for rowIndex, row := range nm.Nodes {
		var (
			node      *Node
			nodeIndex int
		)

		if len(row) > 0 {
			node = row[0]
		}
		for i := 0; i < nm.Width; i++ {
			if node != nil && node.Offset.X == i {

				// display correct node type
				sb.WriteString(node.String())

				// advance node pointer
				nodeIndex++
				if nodeIndex < len(row) {
					node = row[nodeIndex]
				}
			} else if i == 0 {
				sb.WriteString(fmt.Sprintf("%v ", rowIndex%10))
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
					n.Offset.X,
					terminator,
				),
			)
		}
		sb.WriteString(" }\n")
	}
	sb.WriteString("}")

	return sb.String()
}
