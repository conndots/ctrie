package trie

type Trie struct {
	Branches  map[byte]*Trie
	End       bool
	LeafValue interface{}
	Size      uint
}

func newNode(end bool, leaf interface{}) *Trie {
	return &Trie{
		Branches:  make(map[byte]*Trie),
		End:       end,
		LeafValue: leaf,
		Size:      0,
	}
}

func NewTrie() *Trie {
	return &Trie{
		Branches:  make(map[byte]*Trie),
		End:       false,
		LeafValue: nil,
		Size:      0,
	}
}

func (t *Trie) Add(text []byte, leafValue interface{}) {
	curr := t
	t.Size++
	for idx, by := range text {
		branch, ok := curr.Branches[by]
		if !ok {
			if idx == len(text)-1 {
				branch = newNode(true, leafValue)
			} else {
				branch = newNode(false, nil)
			}
			curr.Branches[by] = branch
		}
		if idx == len(text)-1 {
			branch.End = true
			branch.LeafValue = leafValue
		}
		curr = branch
		curr.Size++
	}
}

// GetCandidateLeafs ~ 结果按照最长匹配近似度排序.我们假设的是target是完整的，我们要获取所有能够前缀匹配到target的leafs。并且按照最长匹配的原则排序
func (t *Trie) GetCandidateLeafs(target []byte) ([]interface{}, bool) {
	candidates := make([]interface{}, 0, 5)

	curr := t
	fullMatch := false
	for idx, by := range target {
		branch, ok := curr.Branches[by]
		if !ok {
			break
		}
		if branch.End && branch.LeafValue != nil {
			candidates = append(candidates, branch.LeafValue)
			if idx == len(target)-1 {
				fullMatch = true
			}
		}

		curr = branch
	}

	//reverse it, because the longest match matters.
	for st, end := 0, len(candidates)-1; st < end; st, end = st+1, end-1 {
		candidates[st], candidates[end] = candidates[end], candidates[st]
	}

	return candidates, fullMatch
}
