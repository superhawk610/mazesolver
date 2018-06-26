package nodemap

type Input struct {
	N bool
	E bool
	S bool
	W bool
}

func (i *Input) Critical() bool {
	if i.PassthroughX() && !i.N && !i.S {
		return false
	}
	if i.PassthroughY() && !i.E && !i.W {
		return false
	}
	return true
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

type Node struct {
	Offset      uint
	Input       *Input
	IsStart     bool
	IsEnd       bool
	Connections []*Node
}

func NewNode(offset uint, input *Input) *Node {
	return &Node{
		Offset: offset,
		Input:  input,
	}
}

func NewStartNode(offset uint) *Node {
	return &Node{
		Offset:  offset,
		Input:   &Input{S: true},
		IsStart: true,
	}
}

func NewEndNode(offset uint) *Node {
	return &Node{
		Offset: offset,
		Input:  &Input{N: true},
		IsEnd:  true,
	}
}

func (n *Node) DeadEnd() bool {
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

func (from *Node) Connect(to *Node) {
	from.Connections = append(
		from.Connections,
		to,
	)
}

func (n *Node) String() string {
	// special nodes
	if n.IsStart {
		return "╷"
	}
	if n.IsEnd {
		return "╵"
	}
	if n.DeadEnd() {
		return "x"
	}
	if n.Input.Critical() {
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
