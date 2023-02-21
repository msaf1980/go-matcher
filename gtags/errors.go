package gtags

type ErrQueryInvalid struct {
	Query string
}

func (e ErrQueryInvalid) Error() string {
	return "wrong seriesByTag call: " + e.Query
}

type ErrExprInvalid struct {
	Node string
}

func (e ErrExprInvalid) Error() string {
	return "wrong seriesByTag expr: " + e.Node
}

type ErrExprOverflow struct {
	Node string
}

func (e ErrExprOverflow) Error() string {
	return "overflow seriesByTag expr: " + e.Node
}

type ErrNodeNotTerminated struct {
	Node string
}

func (e ErrNodeNotTerminated) Error() string {
	return "seriesByTag node contains no childs or Terminated: []string{" + e.Node
}

type ErrPathInvalid struct {
	Node   string
	Reason string
}

func (e ErrPathInvalid) Error() string {
	return "invalid path: '" + e.Node + "' " + e.Reason
}
