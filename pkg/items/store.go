package items

type Store interface {
	Store(s string, index int)
}

type MinStore struct {
	Min int
}

func NewMinStore() *MinStore {
	return &MinStore{-1}
}

func (s *MinStore) Init() {
	s.Min = -1
}

func (s *MinStore) Store(_ string, index int) {
	if s.Min < 0 || s.Min > index {
		s.Min = index
	}
}

type IndexStore struct {
	N []int
}

func NewIndexStore() *IndexStore {
	return &IndexStore{}
}

func (s *IndexStore) Init() {
	if len(s.N) > 0 {
		s.N = s.N[:0]
	}
}

func (s *IndexStore) Grow(n int) {
	if n > cap(s.N) {
		newN := make([]int, len(s.N), n)
		copy(newN, s.N)
		s.N = newN
	}
}

func (s *IndexStore) Store(_ string, index int) {
	s.N = append(s.N, index)
}

type StringStore struct {
	S []string
}

func NewStringStore() *StringStore {
	return &StringStore{}
}

func (s *StringStore) Init() {
	if len(s.S) > 0 {
		s.S = s.S[:0]
	}
}

func (s *StringStore) Grow(n int) {
	if n > cap(s.S) {
		newS := make([]string, len(s.S), n)
		copy(newS, s.S)
		s.S = newS
	}
}

func (s *StringStore) Store(sn string, _ int) {
	s.S = append(s.S, sn)
}

type AllStore struct {
	Min   MinStore
	Index IndexStore
	S     StringStore
}

func NewAllStore() *AllStore {
	return &AllStore{}
}

func (s *AllStore) Init() {
	s.Min.Init()
	s.Index.Init()
	s.S.Init()
}

func (s *AllStore) Grow(n int) {
	s.Index.Grow(n)
	s.S.Grow(n)
}

func (s *AllStore) Store(sn string, index int) {
	s.S.Store(sn, index)
	s.Index.Store(sn, index)
	s.Min.Store(sn, index)
}

type Terminated struct {
	Terminate bool
	Query     string // end of chain (resulting seriesByTag)
	Index     int    // resulting seriesByTag index
}

type TreeItem struct {
	Item

	Reverse bool // for suffix or may be other, only last item can be reversed

	Terminated

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*TreeItem // next possible parts slice
}

func LocateChildTreeItem(childs []*TreeItem, node Item, reverse bool) *TreeItem {
	for _, child := range childs {
		if child.Reverse == reverse && child.Equal(node) {
			return child
		}
	}
	return nil
}
