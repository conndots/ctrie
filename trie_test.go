package trie

import "testing"

func getPreparedTrie() *Trie {
	trie := NewTrie()
	trie.Add([]byte("www.google"), 1)
	trie.Add([]byte("www."), 2)
	trie.Add([]byte("www.google.hk."), 3)
	trie.Add([]byte("www.google.us."), 4)
	trie.Add([]byte("www.google.uk"), 5)
	trie.Add([]byte("www.google.uk.wtf"), 6)
	return trie
}

func TestTrieAdd(t *testing.T) {
	trie := NewTrie()
	trie.Add([]byte("www.fuck"), true)
	if trie.Size != 1 {
		t.Errorf("trie add: size expected: 1. got: %v", trie.Size)
	}
	trie.Add([]byte("www.hell"), true)
	if trie.Size != 2 {
		t.Errorf("trie add: size expected: 2, got: %v", trie.Size)
	}
	_, fullMatch := trie.GetCandidateLeafs([]byte("www.fuck"))
	if !fullMatch {
		t.Errorf("trie add: get expected fullMatch=%v, got: %v", true, false)
	}
	_, fullMatch = trie.GetCandidateLeafs([]byte("www.hell"))
	if !fullMatch {
		t.Errorf("trie add: get expected fullMatch=%v, got: %v", true, false)
	}
	cand, fullMatch := trie.GetCandidateLeafs([]byte("www.fuck.com"))
	if fullMatch {
		t.Errorf("trie add: get expected fullMatch=%v, got: %v", false, true)
	}
	if len(cand) != 1 {
		t.Errorf("trie add: get expected len(candidates) == 1, got:%v", len(cand))
	}
}

func TestTrieGetCandidates(t *testing.T) {
	trie := getPreparedTrie()
	candidates, fullMatch := trie.GetCandidateLeafs([]byte("www.google.uk.wtf.fuck"))
	if fullMatch || len(candidates) != 4 {
		t.Errorf("www.google.uk.wtf.fuck expected: fullMatch=%v, candidates length: 4; got: %v, %v", false, fullMatch, candidates)
	}

	if candidates[0] != 6 || candidates[1] != 5 || candidates[2] != 1 || candidates[3] != 2 {
		t.Errorf("www.google.uk.wtf.fuck expected: %v got: %v", []int{6, 5, 1, 2}, candidates)
	}
}

func BenchmarkTrieGetCandidates(b *testing.B) {
	trie := getPreparedTrie()
	for n := 0; n < b.N; n++ {
		trie.GetCandidateLeafs([]byte("www.google.uk.wtf.fuck.hello.what.the.fuck"))
	}
}
