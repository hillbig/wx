wx
=========

wx is a Go library for a succint representation of trie

wx stores a set of strings, and supports PrefixMatch and PredictiveMatch operations.
Each string is assigned to an unique ID from [0, num), and a user can retrieve a string
from ID using Get() method.

Trie tree information and leading edge information is stored using LOUDS representation.

Usage
=====
```
import "github.com/hillbig/wx"

// To use wx, first prepare the builder by wx.NewBuilder
// and then set keys using Add()
wxb := wx.NewBuilder()
wxb.Add("to")
wxb.Add("tea")
wxb.Add("ten")
wxb.Add("ten") // duplicated key is removed
wxb.Add("i")
wxb.Add("in")
wxb.Add("inn")
wxb.Add("we")

w := wxb.Buid()

fmt.Printf("%d\n", w.Num()) // 7 (Note: duplicated "ten" is not added)

// If found, ok = true, and id corresponds to the unique id assigned by wx
ret, ok := w.ExactMatch("tea cup.")

// Can retrieve string using id
fmt.Printf("I like %s\n", ok, w.Get(id)) // I like tea

// LongestPrefixMatch(str) returns the key that matches the longest prefix of str
key := "tea cup"
ret, ok := w.LongestPrefixMatch(key)
fmt.Printf("%s %s", w.Get(ret.ID), key[0:ret.Len]) // tea tea

// CommonPrefixMatch(str) returns all the keys that match prefix of str
rets := w.CommonPrefixMatch("innnnnn")
for _, r := range rets {
	fmt.Printf("%s\n", w.Get(r.ID))
}
// i
// in
// inn

// Predictivematch(str) returns all the keys whose strict prefix matches to str
ids := w.PredictiveMatch("i")
for _ id := range ids {
	fmt.Printf("%s\n", w.Get(id))
}
// in
// inn

// If you limit the number of returns, use WithLimit versions

w.CommonPrefixMatchWithLimit("innnnnn", 2) // i in
w.PredictiveMatchWithLimit("i", 2) // i in

// Encode to binary representation
bytes, err := vs.MarshalBinary()
newvs := vecstring.New()

// Decode from binary presentation
err := newvs.UnmarshalBinary(bytes)
```
