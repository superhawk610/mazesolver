package nodemap

import (
	"image"
)

const (
	Up    = 0 // 0b00
	Down  = 1 // 0b01
	Left  = 2 // 0b10
	Right = 3 // 0b11
)

type Input struct {
	N bool
	E bool
	S bool
	W bool
}

type Connection struct {
	Direction int
	Length    int
	Parent    *Node
	Node      *Node
	Used      bool
}

type Offset struct {
	X, Y int
}

type Node struct {
	Offset      *Offset
	Input       *Input
	IsStart     bool
	IsEnd       bool
	Connections []*Connection
}

func NewNode(x int, y int, input *Input) *Node {
	return &Node{
		Offset: &Offset{X: x, Y: y},
		Input:  input,
	}
}

func NewStartNode(x int, y int) *Node {
	return &Node{
		Offset:  &Offset{X: x, Y: y},
		Input:   &Input{S: true},
		IsStart: true,
	}
}

func NewEndNode(x int, y int) *Node {
	return &Node{
		Offset: &Offset{X: x, Y: y},
		Input:  &Input{N: true},
		IsEnd:  true,
	}
}

func (c *Connection) Use() {
	if c.Used {
		return
	}

	c.Used = true
	for _, con := range c.Node.Connections {
		if con.Node == c.Parent {
			con.Used = true
			return
		}
	}
}

func (c *Connection) OffsetAt(i int) image.Point {
	var shiftUp, shiftRight, shiftDown, shiftLeft int

	if c.Direction == Up {
		shiftUp = 1
	}
	if c.Direction == Right {
		shiftRight = 1
	}
	if c.Direction == Down {
		shiftDown = 1
	}
	if c.Direction == Left {
		shiftLeft = 1
	}
	return image.Point{
		X: c.Parent.Offset.X - (i * shiftLeft) + (i * shiftRight),
		Y: c.Parent.Offset.Y - (i * shiftUp) + (i * shiftDown),
	}
}

func (i *Input) PassthroughX() bool {
	return i.E && i.W
}

func (i *Input) PassthroughY() bool {
	return i.N && i.S
}

func (i *Input) PassthroughXY() bool {
	return i.PassthroughX() && i.PassthroughY()
}

func (n *Node) Critical() bool {
	i := n.Input
	if i.PassthroughX() && !i.N && !i.S {
		return false
	}
	if i.PassthroughY() && !i.E && !i.W {
		return false
	}
	return true
}

func (n *Node) RemainingConnections() int {
	if len(n.Connections) == 0 {
		return 0
	}
	var remaining int
	for _, con := range n.Connections {
		if !con.Used {
			remaining++
		}
	}
	return remaining
}

func (n *Node) Fresh() bool {
	return len(n.Connections)-1 == n.RemainingConnections()
}

func (n *Node) Exhausted() bool {
	return n.RemainingConnections() == 0
}

func (n *Node) DeadEnd() bool {
	if n.IsStart {
		return false
	}
	var inputs int

	if n.Input.N {
		inputs++
	}
	if n.Input.E {
		inputs++
	}
	if n.Input.S {
		inputs++
	}
	if n.Input.W {
		inputs++
	}

	return inputs == 1
}

func (from *Node) Connect(to *Node, direction int, length int) {
	from.Connections = append(
		from.Connections,
		&Connection{
			Direction: direction,
			Length:    length,
			Parent:    from,
			Node:      to,
		},
	)
	to.Connections = append(
		to.Connections,
		&Connection{
			Direction: direction ^ 1, // 0b01
			Length:    length,
			Parent:    to,
			Node:      from,
		},
	)
}

func (n *Node) String() string {
	// special nodes
	if n.IsStart {
		return "S"
	}
	if n.IsEnd {
		return "E"
	}
	if n.DeadEnd() {
		return "x"
	}
	if n.Critical() {
		return "o"
	}

	i := n.Input

	// all 4 inputs
	if i.PassthroughXY() {
		return "┼"
	}

	// just 3 inputs
	if i.PassthroughX() {
		if i.N {
			return "┴"
		}
		if i.S {
			return "┬"
		}
		return "─"
	}
	if i.PassthroughY() {
		if i.E {
			return "├"
		}
		if i.W {
			return "┤"
		}
		return "│"
	}

	// just 2 inputs
	if i.N && i.E {
		return "└"
	}
	if i.E && i.S {
		return "┌"
	}
	if i.S && i.W {
		return "┐"
	}
	if i.W && i.N {
		return "┘"
	}

	// just 1 input
	if i.N {
		return "╵"
	}
	if i.E {
		return "╶"
	}
	if i.S {
		return "╷"
	}
	if i.W {
		return "╴"
	}

	return "?"
}
