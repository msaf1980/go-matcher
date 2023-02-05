package gglob

type ErrNodeMissmatch struct {
	typ  NodeType
	node string
}

func (e ErrNodeMissmatch) Error() string {
	return "node type " + e.typ.String() + " mismatch: " + e.node
}

type ErrNodeNotEnd struct {
	node string
}

func (e ErrNodeNotEnd) Error() string {
	return "node contains no childs or terminated: " + e.node
}

type ErrGlobNotExpanded struct {
	node string
}

func (e ErrGlobNotExpanded) Error() string {
	return "inner node contains 0 inners: " + e.node
}

type ErrNodeEmpty struct {
	path string
}

func (e ErrNodeEmpty) Error() string {
	return "empty node in path: " + e.path
}

type ErrNodeUnclosed struct {
	segment string
}

func (e ErrNodeUnclosed) Error() string {
	return "glob contain unclosed node segment: " + e.segment
}
