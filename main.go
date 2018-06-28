package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/superhawk610/mazesolver/nodemap"
	"golang.org/x/image/bmp"
)

const (
	mazeFile     = "./mazes/easy.bmp"
	solutionFile = "./mazes/solution.bmp"
)

func main() {
	mazeBmp, err := os.Open(mazeFile)
	if err != nil {
		log.Fatal(err)
	}
	maze, err := bmp.Decode(mazeBmp)
	if err != nil {
		log.Fatal(err)
	}
	nm := nodemap.FromMaze(&maze)
	nm.Stat()

	// fmt.Println(nm.Visualize())

	start := time.Now()
	solution := nodemap.SolveDumb(nm)
	// fmt.Println(solution.Path())
	now := time.Now()
	fmt.Printf(
		"Solution found using SolveDumb in %v\n",
		now.Sub(start),
	)

	os.Remove(solutionFile)
	f, err := os.Create(solutionFile)
	if err != nil {
		log.Fatal(err)
	}
	err = nodemap.WriteSolution(f, maze, solution)
	if err != nil {
		log.Fatal(err)
	}
}
