package wx

import (
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
	wxtmp := wxBuilderTmp{
		keys:                keys,
		branches:            vecstring.NewForWX(),
		terminals:           rsdic.New(),
		leadingIsZeroOrOnes: rsdic.New(),
		leadingIsOnes:       rsdic.New(),
		leadingOnes:         make([]byte, 0),
		leadingStrs:         newStrs(),
	}
	if len(keys) > 0 {
		wxtmp.bfs(state{0, len(keys), 0})
	}

	// build trie for leadings recursively
	/*
		if !isRawLead {
			leadBuilder := builderImpl{make(map[string]struct{})}
			for key, _ := range wxtmp.leadingStrs.str2id {
				leadBuilder.Add(key)
			}
			return &wxImpl{
				branches:  wxtmp.branches.Build(),
				terminals: wxtmp.terminals.Build(),
				leadings:  leadBuilder.Build(),
				num:       uint64(len(keys)),
			}
		} else {
	*/
	id2key := make([]string, len(wxtmp.leadingStrs.str2id))
	for key, id := range wxtmp.leadingStrs.str2id {
		id2key[id] = key
	}
	leadingRaws := vecstring.New()
	for _, key := range id2key {
		leadingRaws.PushBack(key)
	}
	return &wxImpl{
		branches:            wxtmp.branches,
		terminals:           wxtmp.terminals,
		leadingIDs:          fixvec.NewFromArray(wxtmp.leadingStrs.ids),
		leadings:            leadingRaws,
		leadingIsZeroOrOnes: wxtmp.leadingIsZeroOrOnes,
		leadingIsOnes:       wxtmp.leadingIsOnes,
		leadingOnes:         wxtmp.leadingOnes,
		num:                 uint64(len(keys)),
	}
	/*
		}
	*/
}

type wxBuilderTmp struct {
	keys                []string
	branches            vecstring.VecStringForWX
	leadingStrs         *strs
	leadingIsZeroOrOnes rsdic.RSDic
	leadingIsOnes       rsdic.RSDic
	leadingOnes         []byte
	terminals           rsdic.RSDic
}

func branchDepth(s string, u string, depth int) int {
	for i := depth; ; i++ {
		if len(s) == i || len(u) == i { // since keys are sorted, len(s) == i only hold. But added len(u) in case.
			return i
		}
		if s[i] != u[i] {
			return i
		}
	}
}

func (wxtmp *wxBuilderTmp) bfs(s state) {
	beg, end, depth := s.beg, s.end, s.depth
	keys := wxtmp.keys
	bdepth := branchDepth(keys[beg], keys[end-1], depth)
	leading := keys[beg][depth:bdepth]
	if len(leading) <= 1 {
		wxtmp.leadingIsZeroOrOnes.PushBack(true)
		if len(leading) == 0 {
			wxtmp.leadingIsOnes.PushBack(false)
		} else {
			// len(leading) == 0
			wxtmp.leadingOnes = append(wxtmp.leadingOnes, leading[0])
			wxtmp.leadingIsOnes.PushBack(true)
		}
	} else {
		wxtmp.leadingIsZeroOrOnes.PushBack(false)
		wxtmp.leadingStrs.add(leading)
	}
	if len(keys[beg]) == bdepth {
		// keys[beg] is internal node
		wxtmp.terminals.PushBack(true)
		beg++
		if beg == end {
			wxtmp.branches.PushBack("")
			return
		}
	} else {
		wxtmp.terminals.PushBack(false)
	}

	nextBeg := beg
	c := keys[nextBeg][bdepth]
	branch := make([]byte, 0)
	nextBFS := make([]state, 0)
	for i := beg + 1; i <= end; i++ {
		if i == end || c != keys[i][bdepth] {
			branch = append(branch, c)
			nextBFS = append(nextBFS, state{nextBeg, i, bdepth + 1})
			if i == end {
				break
			}
			nextBeg = i
			c = keys[i][bdepth]
		}
	}
	wxtmp.branches.PushBack(string(branch))
	for _, s := range nextBFS {
		wxtmp.bfs(s)
	}
}

type state struct {
	beg   int
	end   int
	depth int
}
