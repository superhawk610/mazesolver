package nodemap

import (
	"image"
	"image/color"
)

type SolutionImage struct {
	image.Image
	Solution map[image.Point]float64
}

func NewSolutionImage(i image.Image, sol Solution) *SolutionImage {
	mappedSolution := make(map[image.Point]float64)

	totalLength := float64(sol.Length())
	length := 0
	path := sol.Path()
	for i, node := range path {
		if node.IsEnd {
			p := node.Connections[0].OffsetAt(0)
			mappedSolution[p] = float64(length) / totalLength
			continue
		}

		for _, con := range node.Connections {
			if con.Node != path[i+1] {
				continue
			}

			for j := 0; j < con.Length; j++ {
				p := con.OffsetAt(j)
				mappedSolution[p] = float64(length+j) / totalLength
			}
			length += con.Length
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
