// Package wx provides a succinct trie containing a set of strings
package wx

import (
	"github.com/hillbig/fixvec"
	"github.com/hillbig/rsdic"
	"github.com/hillbig/vecstring"
)

// WX represents a trie containing a set of strings
// Given a query Q of length m, WX supports various operations on trie in O(m) time.
// WX is built by using Builder function
type WX interface {
	// ExactMatch examines whether the query exactly matches to any string in WX,
	// and returns (id, true) if found, or return (0, false)
	ExactMatch(str string) (id uint64, ok bool)

	// LongestPrefixMatch returns the longest key that matches the prefix of query.
	LongestPrefixMatch(str string) (ret WXString, ok bool)

	// CommonPrefixMatch returns the all keys that match the prefix of the query
	CommonPrefixMatch(str string) (ret []WXString)

	// CommonPrefixMatch returns the all keys that match the prefix of the query
	// limit indicates the maximum number of matched keys.
	CommonPrefixMatchWithLimit(str string, limit uint64) (ret []WXString)

	// PredictiveMatch returns the all keys whose prefix matches to the query
	PredictiveMatch(str string) (ids []uint64)

	// PredictiveMatchWithLimit returns the all keys whose prefix matches to the query
	// limit indicates the maximum number of matched keys
	PredictiveMatchWithLimit(str string, limit uint64) (ids []uint64)

	// Get returns the key with ID.
	Get(id uint64) string

	// Num returns the number of keys
	Num() uint64

	// MarshalBinary encodes WX into a binary form and returns the result.
	MarshalBinary() ([]byte, error)

	// UnmarshalBinary decodes WX from a binary form generated MarshalBinary
	UnmarshalBinary([]byte) error

	debugPrint()
}

// WXString represents an internal representation of string in WX
// Decode key by Get(ID)
type WXString struct {
	ID  uint64
	Len uint64
}

// Builder is used for building WX.
type Builder interface {
	// Add adds a key to the WX
	// A duplicatad　key is removed.
	Add(key string)

	// Build builds a WX (Builder is not changed)
	Build() WX

	// Num returns the number of registered keys
	Num() uint64

	// TotalByteNum returns the number of bytes of keys.
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
