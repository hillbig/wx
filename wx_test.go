package wx

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestEmptyWX(t *testing.T) {
	Convey("When an empty key/val is given", t, func() {
		wxb := NewBuilder()
		wx := wxb.Build()
		Convey("Num should be 0", func() {
			So(wxb.TotalByteNum(), ShouldEqual, 0)
			So(wx.Num(), ShouldEqual, 0)
		})
	})
}

func TestOneWX(t *testing.T) {
	wxb := NewBuilder()
	key := "abcdefg"
	wxb.Add(key)
	wx := wxb.Build()
	Convey("When one key/val is given", t, func() {
		Convey("Num should be 1", func() {
			So(wxb.Num(), ShouldEqual, 1)
			So(wx.Num(), ShouldEqual, 1)
		})
		Convey("Get(0) should return abcdefg", func() {
			So(wx.Get(0), ShouldEqual, key)
		})
		Convey("abcdefg should be found", func() {
			wxstr, ok := wx.LongestPrefixMatch(key)
			So(wxstr.ID, ShouldEqual, 0)
			So(ok, ShouldEqual, true)
		})
		Convey("abcd should not be found", func() {
			_, ok := wx.LongestPrefixMatch("abcd")
			So(ok, ShouldEqual, false)
		})
	})
}

func TestMutipleWX(t *testing.T) {
	wxb := NewBuilder()
	wxb.Add("abc")
	wxb.Add("abe")
	wxb.Add("abe")
	wxb.Add("a")
	wx := wxb.Build()
	Convey("When multiple key/vals are given", t, func() {
		Convey("Num should be 3", func() {
			So(wx.Num(), ShouldEqual, 3)
		})
		Convey("abc should be found", func() {
			id, ok := wx.ExactMatch("abc")
			So(ok, ShouldEqual, true)
			So(wx.Get(id), ShouldEqual, "abc")
		})
		Convey("bbc should not be found", func() {
			_, ok := wx.ExactMatch("bbc")
			So(ok, ShouldEqual, false)
		})
		Convey("ab should not be found", func() {
			_, ok := wx.ExactMatch("ab")
			So(ok, ShouldEqual, false)
		})
		Convey("abe should not be found", func() {
			_, ok := wx.ExactMatch("abd")
			So(ok, ShouldEqual, false)
		})
		Convey("abc CommonPrefixMatch returns a and abc", func() {
			rets := wx.CommonPrefixMatch("abc")
			So(len(rets), ShouldEqual, 2)
			So(wx.Get(rets[0].ID), ShouldEqual, "a")
			So(wx.Get(rets[1].ID), ShouldEqual, "abc")
		})
		Convey("abc CommonPrefixMatchWithLimit(0) return 0", func() {
			rets := wx.CommonPrefixMatchWithLimit("abc", 0)
			So(len(rets), ShouldEqual, 0)
		})
		Convey("abc CommonPrefixMatchWithLimit(1) return 1", func() {
			rets := wx.CommonPrefixMatchWithLimit("abc", 1)
			So(len(rets), ShouldEqual, 1)
			So(wx.Get(rets[0].ID), ShouldEqual, "a")
		})
		Convey("ab PredictiveMatch returns abc abe", func() {
			ids := wx.PredictiveMatch("ab")
			So(len(ids), ShouldEqual, 2)
			So(wx.Get(ids[0]), ShouldEqual, "abc")
			So(wx.Get(ids[1]), ShouldEqual, "abe")
		})
		Convey("ab PredictiveMatchWithLimit returns a abc", func() {
			ids := wx.PredictiveMatchWithLimit("ab", 1)
			So(len(ids), ShouldEqual, 1)
			So(wx.Get(ids[0]), ShouldEqual, "abc")
		})
		Convey("MarshalBinary", func() {
			out, err := wx.MarshalBinary()
			So(err, ShouldBeNil)
			newwx := New()
			err = newwx.UnmarshalBinary(out)
			So(err, ShouldBeNil)
			So(newwx.Num(), ShouldEqual, 3)
			rets := newwx.CommonPrefixMatch("abc")
			So(len(rets), ShouldEqual, 2)
			So(newwx.Get(rets[0].ID), ShouldEqual, "a")
			So(newwx.Get(rets[1].ID), ShouldEqual, "abc")
		})
	})
}

func BenchmarkWXBuild(b *testing.B) {
	num := 1000000
	maxLen := 10
	wxb := NewBuilder()
	totalLen := 0
	for i := 0; i < num; i++ {
		l := rand.Int() % maxLen
		strbuf := make([]byte, l)
		for j := 0; j < l; j++ {
			strbuf[j] = byte(rand.Int() % 256)
		}
		totalLen += l
		wxb.Add(string(strbuf))
	}
	dummy := uint64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wx := wxb.Build()
		dummy += wx.Num()
	}
	if dummy == 777 {
		dummy = 7777 // I'm very luckey. suppress optimization
	}
}

func BenchmarkWXExactMatch(b *testing.B) {
	num := 1000000
	maxLen := 10
	wxb := NewBuilder()
	origs := make([]string, num)
	for i := 0; i < num; i++ {
		l := rand.Int() % maxLen
		strbuf := make([]byte, l)
		for j := 0; j < l; j++ {
			strbuf[j] = byte(rand.Int() % 256)
		}
		origs[i] = string(strbuf)
		wxb.Add(string(strbuf))
	}
	dummy := uint64(0)
	wx := wxb.Build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ind := rand.Int() % num
		_, ok := wx.ExactMatch(origs[ind])
		if ok {
			dummy++
		}
	}
	if dummy == 777 {
		dummy = 7777 // I'm very luckey. suppress optimization
	}
}
