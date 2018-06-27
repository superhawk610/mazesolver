package nodemap

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
)

type SolutionImage struct {
	image.Image
	Solution map[image.Point]float64
}

func NewSolutionImage(i image.Image, sol Solution) *SolutionImage {
	mappedSolution := make(map[image.Point]float64)

	totalLength := float64(sol.Length())
	fmt.Println(totalLength)
	for i, node := range sol.Path {
		if node.IsEnd {
			p := node.Connections[0].OffsetAt(0)
			mappedSolution[p] = float64(i) / totalLength
			continue
		}

		for _, con := range node.Connections {
			if con.Node != sol.Path[i+1] {
				continue
			}

			for j := 0; j < con.Length; j++ {
				p := con.OffsetAt(j)
				mappedSolution[p] = float64(i) / totalLength
			}
			break
		}
	}

	return &SolutionImage{
		Image:    i,
		Solution: mappedSolution,
	}
}

func (sm *SolutionImage) At(x, y int) color.Color {
	p := image.Point{
		X: x,
		Y: y,
	}
	if percentage, ok := sm.Solution[p]; ok {
		return color.RGBA{
			R: uint8(255 * (1 - percentage)),
			G: 0,
			B: uint8(255 * percentage),
			A: 255,
		}
	}

	return sm.Image.At(x, y)
}

type Solution struct {
	Path []*Node
}

func (s *Solution) Length() int {
	length := 0

	for _, node := range s.Path {
		if node.IsEnd {
			length++
			break
		}
		for i, con := range node.Connections {
			if con.Node != s.Path[i+1] {
				continue
			}

			fmt.Println(con.Length)
			length += con.Length
		}
	}

	return length
}

var inflections []*Node
var solution Solution

func SolveDumbRecurse(n *Node) bool {
	fmt.Printf(
		"moving to node at {x: %v, y: %v}\n",
		n.Offset.X, n.Offset.Y,
	)

	if len(inflections) == 0 {
		solution.Path = append(
			solution.Path,
			n,
		)
	}

	if n.IsEnd {
		return true
	}

	if n.Fresh() && n.RemainingConnections() > 1 {
		inflections = append(
			inflections,
			n,
		)
	}

	for _, con := range n.Connections {
		if con.Used {
			continue
		}

		con.Use()
		return SolveDumbRecurse(con.Node)
	}

	if len(inflections) > 0 {
		infl := inflections[len(inflections)-1]
		if infl.RemainingConnections() == 1 {
			inflections = inflections[:len(inflections)-1]
		}

		fmt.Printf("\n- returning to inflection at {x: %v, y: %v}\n", infl.Offset.X, infl.Offset.Y)
		return SolveDumbRecurse(infl)
	}

	panic("No solution found.")
	return false
}

func SolveDumb(nm *NodeMap) Solution {
	// get pointer to start node
	n := nm.Nodes[0][0]

	SolveDumbRecurse(n)
	return solution
}

func WriteSolution(w io.Writer, maze image.Image, sol Solution) error {
	output := NewSolutionImage(
		maze,
		sol,
	)
	return bmp.Encode(w, output)
}
