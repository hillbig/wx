// Package wx provides a succinct trie for a set of strings
package wx

import (
	"github.com/hillbig/fixvec"
	"github.com/hillbig/rsdic"
	"github.com/hillbig/vecstring"
)

type WX interface {
	ExactMatch(str string) (id uint64, ok bool)
	LongestPrefixMatch(str string) (ret WXString, ok bool)
	CommonPrefixMatch(str string) (ret []WXString)
	CommonPrefixMatchWithLimit(str string, limit uint64) (ret []WXString)
	PredictiveMatch(str string) (ids []uint64)
	PredictiveMatchWithLimit(str string, limit uint64) (ids []uint64)
	Get(id uint64) string
	Num() uint64
	debugPrint()

	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
}

// WXString represents an internal representation of string in WX
// Decode key by Get(ID)
type WXString struct {
	ID  uint64
	Len uint64
}

type Builder interface {
	Add(key string)
	Build() WX
	Num() uint64
	TotalByteNum() uint64
}

func New() WX {
	return &wxImpl{
		branches:            vecstring.NewForWX(),
		terminals:           rsdic.New(),
		leadingIDs:          fixvec.New(0, 0),
		leadings:            vecstring.New(),
		leadingIsZeroOrOnes: rsdic.New(),
		leadingIsOnes:       rsdic.New(),
		leadingOnes:         make([]byte, 0),
		num:                 0,
	}
}

func NewBuilder() Builder {
	return &builderImpl{make(map[string]struct{})}
}
