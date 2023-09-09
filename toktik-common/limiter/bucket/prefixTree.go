package bucket

// 前缀树

type PrefixTree struct {
	suffix map[string]*PrefixTree
	result interface{}
}

func NewPrefixTree() *PrefixTree {
	return &PrefixTree{suffix: make(map[string]*PrefixTree)}
}

func (t *PrefixTree) Put(prefix []string, v interface{}) {
	root := t
	for _, s := range prefix {
		if root.suffix[s] == nil {
			root.suffix[s] = NewPrefixTree()
		}
		root = root.suffix[s]
	}
	root.result = v
}

func (t *PrefixTree) Get(prefix []string) interface{} {
	root := t
	for _, s := range prefix {
		if root.suffix[s] != nil {
			root = root.suffix[s]
		} else {
			break
		}
	}
	return root.result
}
