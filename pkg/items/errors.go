package items

type ErrNodeMissmatch struct {
	Type string
	Node string
}

func (e ErrNodeMissmatch) Error() string {
	return "node type '" + e.Type + "'mismatch: " + e.Node
}

type ErrNodeNotEnd struct {
	Node string
}

func (e ErrNodeNotEnd) Error() string {
	return "node contains no childs or terminated: " + e.Node
}

type ErrGlobNotExpanded struct {
	Node string
}

func (e ErrGlobNotExpanded) Error() string {
	return "empty node: " + e.Node
}

type ErrNodeEmpty struct {
	Path string
}

func (e ErrNodeEmpty) Error() string {
	return "empty node in path: " + e.Path
}

type ErrNodeUnclosed struct {
	Segment string
}

func (e ErrNodeUnclosed) Error() string {
	return "glob contain unclosed node segment: " + e.Segment
}
