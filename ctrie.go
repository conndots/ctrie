package trie

import "fmt"

type NodeType uint8

const (
	NodeTypeRoot    NodeType = 0
	NodeTypeLeaf    NodeType = 2
	NodeTypeDefault NodeType = 1
)

type CTrie struct {
	childrenIdx map[byte]*CTrie
	Size        int

	LeafValue interface{}
	path      string
	nodeType  NodeType
}

func NewCompressedTrie() *CTrie {
	return &CTrie{
		childrenIdx: make(map[byte]*CTrie),
		Size:        0,
		LeafValue:   nil,
		nodeType:    NodeTypeRoot,
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (ct *CTrie) Add(str string, value interface{}) {
	ct.Size++
	if len(ct.path) > 0 || len(ct.childrenIdx) > 0 {
	loopStart:
		for {
			diffSt := 0
			minLen := min(len(ct.path), len(str))
			for diffSt < minLen && str[diffSt] == ct.path[diffSt] {
				diffSt++
			}

			if diffSt < len(ct.path) { //split the node
				child := &CTrie{
					childrenIdx: ct.childrenIdx,
					path:        ct.path[diffSt:],
					Size:        ct.Size,
					LeafValue:   ct.LeafValue,
					nodeType:    ct.nodeType,
				}
				if ct.nodeType == NodeTypeRoot {
					if ct.LeafValue != nil {
						child.nodeType = NodeTypeLeaf
					} else {
						child.nodeType = NodeTypeDefault
					}
				}

				ct.childrenIdx = make(map[byte]*CTrie)
				ct.childrenIdx[ct.path[diffSt]] = child
				ct.path = ct.path[:diffSt]
				if ct.nodeType == NodeTypeLeaf {
					ct.nodeType = NodeTypeDefault
				}
				ct.LeafValue = nil
			}

			if diffSt < len(str) { //str has diff
				str = str[diffSt:]
				c := str[0]

				sub, existed := ct.childrenIdx[c]
				if existed {
					ct = sub
					continue loopStart
				}

				//add to a new child
				child := &CTrie{
					LeafValue: value,
					nodeType:  NodeTypeLeaf,
					path:      str,
					Size:      1,
				}
				if ct.childrenIdx == nil {
					ct.childrenIdx = make(map[byte]*CTrie)
				}
				ct.childrenIdx[c] = child
				return
			} else if diffSt == len(str) {
				if ct.LeafValue == nil {
					ct.LeafValue = value
					if ct.nodeType != NodeTypeRoot {
						ct.nodeType = NodeTypeLeaf
					}
				}
				return
			}
		}
	} else {
		ct.path = str
		ct.LeafValue = value
		ct.nodeType = NodeTypeRoot
	}
}

func (ct *CTrie) GetCandidateLeafs(target string) (candidates []interface{}, fullMatch bool) {
	defer func() {
		//reverse it, because the longest match matters.
		for st, end := 0, len(candidates)-1; st < end; st, end = st+1, end-1 {
			candidates[st], candidates[end] = candidates[end], candidates[st]
		}
	}()
	candidates = make([]interface{}, 0, 2)
	tlen := len(target)

	for {
		if ct.path == target {
			candidates = append(candidates, ct.LeafValue)
			fullMatch = true
			return
		} else if tlen < len(ct.path) {
			fullMatch = false
			return
		} else {
			st := 0
			ctlen := len(ct.path)
			tlen := len(target)
			for st < ctlen && st < tlen && ct.path[st] == target[st] {
				st++
			}
			if st == tlen || st != ctlen {
				fullMatch = false
				return
			}
			target = target[st:]
			if ct.nodeType == NodeTypeLeaf || (ct.LeafValue != nil && ct.nodeType == NodeTypeRoot) {
				candidates = append(candidates, ct.LeafValue)
			}
			c := target[0]
			sub, existed := ct.childrenIdx[c]
			if existed {
				ct = sub
			} else {
				fullMatch = false
				return
			}
		}
	}
}

func (ct *CTrie) print() {
	type Node struct {
		node  *CTrie
		depth int
	}
	queue := make([]*Node, 0, 10)
	queue = append(queue, &Node{node: ct, depth: 0})
	currDepth := 0
	fmt.Println("-------------------------")
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		depth := curr.depth
		if depth > currDepth {
			fmt.Println()
			currDepth = depth
		}
		fmt.Print("[", curr.node.path, " depth(", depth, ") nodeType:", curr.node.nodeType, "]\t")

		for _, sub := range curr.node.childrenIdx {
			queue = append(queue, &Node{node: sub, depth: depth + 1})
		}
	}
	fmt.Println("\n-------------------------")
}
