package main

import (
	"fmt"
	"github.com/superhawk610/mazesolver/nodemap"
	"golang.org/x/image/bmp"
	"log"
	"os"
)

const mazeFile = "./mazes/easy.bmp"

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

	fmt.Println(nm.Visualize())
}
