package nodemap

import (
	"image"
	"io"

	"golang.org/x/image/bmp"
)

type Solution struct {
	segments [][]*Node
}

func (s *Solution) Path() []*Node {
	var path []*Node

	for _, seg := range s.segments {
		path = append(
			path,
			seg...,
		)
	}

	return path
}

func (s *Solution) StepIn() {
	s.segments = append(
		s.segments,
		make([]*Node, 0),
	)
}

func (s *Solution) Walk(n *Node) {
	if len(s.segments) == 0 {
		s.StepIn()
	}
	currentPath := &s.segments[len(s.segments)-1]

	*currentPath = append(
		*currentPath,
		n,
	)
}

func (s *Solution) StepBack() {
	s.segments = s.segments[:len(s.segments)-1]
}

func (s *Solution) Length() int {
	length := 0
	path := s.Path()

	for i, node := range path {
		if node.IsEnd {
			length++
			break
		}
		for _, con := range node.Connections {
			if con.Node != path[i+1] {
				continue
			}

			length += con.Length
		}
	}

	return length
}

var inflections []*Node
var solution Solution

func SolveDumbRecurse(n *Node) bool {
	// fmt.Printf(
	// 	"moving to node at (x: %v, y: %v)\n",
	// 	n.Offset.X,
	// 	n.Offset.Y,
	// )

	// fmt.Printf(
	// 	"%v\n",
	// 	solution.segments,
	// )

	if n.IsEnd {
		solution.Walk(n)
		return true
	}

	if n.Fresh() && n.RemainingConnections() > 1 {
		inflections = append(
			inflections,
			n,
		)
		solution.StepIn()
	}

	for _, con := range n.Connections {
		if con.Used {
			continue
		}

		con.Use()
		solution.Walk(n)
		return SolveDumbRecurse(con.Node)
	}

	if len(inflections) > 0 {
		infl := inflections[len(inflections)-1]
		if infl.RemainingConnections() == 1 {
			inflections = inflections[:len(inflections)-1]
			solution.StepBack()
		}

		// fmt.Printf(
		// 	"returning to inflection point at (x: %v, y: %v)\n",
		// 	infl.Offset.X,
		// 	infl.Offset.Y,
		// )
		return SolveDumbRecurse(infl)
	}

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
