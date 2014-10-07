package wx

import (
	"fmt"
	"github.com/hillbig/fixvec"
	"github.com/hillbig/rsdic"
	"github.com/hillbig/vecstring"
	"sort"
)

type builderImpl struct {
	keys map[string]struct{}
}

func (wxb *builderImpl) Add(key string) {
	wxb.keys[key] = struct{}{}
}

type strs struct {
	str2id map[string]uint64
	ids    []uint64
}

func newStrs() *strs {
	return &strs{
		str2id: make(map[string]uint64),
		ids:    make([]uint64, 0),
	}
}

func (s *strs) add(str string) {
	id, ok := s.str2id[str]
	if !ok {
		id = uint64(len(s.str2id))
		s.str2id[str] = id
	}
	s.ids = append(s.ids, id)
}

func (wxb builderImpl) Build() WX {
	return wxb.build(false)
}

func (wxb builderImpl) Num() uint64 {
	return uint64(len(wxb.keys))
}

func (wxb builderImpl) TotalByteNum() uint64 {
	totalByteNum := uint64(0)
	for key, _ := range wxb.keys {
		totalByteNum += uint64(len(key))
	}
	return totalByteNum
}

func (wxb *builderImpl) build(isRawLead bool) WX {
	keys := make([]string, 0)
	for key, _ := range wxb.keys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	branches := vecstring.NewForWX()
	terminals := rsdic.New()
	leadingIsZeroOrOnes := rsdic.New()
	leadingIsOnes := rsdic.New()
	leadingOnes := make([]byte, 0)
	leadingStrs := newStrs()

	q := newNodeQueue(0)
	if len(keys) > 0 {
		q.push(&node{0, len(keys), 0})
	}
	for q.num > 0 {
		n := q.pop()
		beg, end, depth := n.beg, n.end, n.depth
		bdepth := branchDepth(keys[beg], keys[end-1], depth)
		leading := keys[beg][depth:bdepth]
		if len(leading) <= 1 {
			leadingIsZeroOrOnes.PushBack(true)
			if len(leading) == 0 {
				leadingIsOnes.PushBack(false)
			} else {
				// len(leading) == 1
				leadingIsOnes.PushBack(true)
				leadingOnes = append(leadingOnes, leading[0])
			}
		} else {
			leadingIsZeroOrOnes.PushBack(false)
			leadingStrs.add(leading)
		}
		if len(keys[beg]) == bdepth {
			// keys[beg] is internal node
			terminals.PushBack(true)
			beg++
			if beg == end {
				branches.PushBack("")
				continue
			}
		} else {
			terminals.PushBack(false)
		}

		nextBeg := beg
		c := keys[nextBeg][bdepth]
		branch := make([]byte, 0)
		for i := beg + 1; i <= end; i++ {
			if i == end || c != keys[i][bdepth] {
				branch = append(branch, c)
				q.push(&node{nextBeg, i, bdepth + 1})
				if i == end {
					break
				}
				nextBeg = i
				c = keys[i][bdepth]
			}
		}
		branches.PushBack(string(branch))
	}
	id2key := make([]string, len(leadingStrs.str2id))
	for key, id := range leadingStrs.str2id {
		id2key[id] = key
	}
	leadings := vecstring.New()
	for _, key := range id2key {
		leadings.PushBack(key)
	}
	return &wxImpl{
		branches:            branches,
		terminals:           terminals,
		leadingIDs:          fixvec.NewFromArray(leadingStrs.ids),
		leadings:            leadings,
		leadingIsZeroOrOnes: leadingIsZeroOrOnes,
		leadingIsOnes:       leadingIsOnes,
		leadingOnes:         leadingOnes,
		num:                 uint64(len(keys)),
	}
}

func printString(s string) {
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i])
	}
	fmt.Printf("\n")
}

func branchDepth(s string, u string, depth int) int {
	for i := depth; ; i++ {
		if len(s) == i || len(u) == i { // since keys are sorted, only len(s) == i should be valid. But added len(u) in case.
			return i
		}
		if s[i] != u[i] {
			return i
		}
	}
}
