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
		Convey("ab PredictiveMatchWithLimit 0 returns nothing", func() {
			ids := wx.PredictiveMatchWithLimit("ab", 0)
			So(len(ids), ShouldEqual, 0)
		})
		Convey("ac PredictiveMatchWithLimit 3 returns nothing", func() {
			ids := wx.PredictiveMatch("ac")
			So(len(ids), ShouldEqual, 0)
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

func TestFailedWX(t *testing.T) {
	m := map[string]struct{}{
		"https://api.github.com/repos/RostovTeam/hackaton/issues/events{/number}":                               struct{}{},
		"https://api.github.com/repos/slickage/rate-limiter/issues/events{/number}":                             struct{}{},
		"https://api.github.com/repos/ConsultingMD/blog/issues/events{/number}":                                 struct{}{},
		"https://api.github.com/repos/OULibraries/bootstrap_subtheme_mdw/issues/events{/number}":                struct{}{},
		"https://api.github.com/repos/VioletGrey/spree_gift_message/issues/events{/number}":                     struct{}{},
		"https://api.github.com/repos/LumenData/AlgorithmsIO-Streaming-Examples/issues/events{/number}":         struct{}{},
		"https://api.github.com/repos/mulesoft/template-sfdc2sfdc-opportunity-migration/issues/events{/number}": struct{}{},
		"https://api.github.com/repos/RealScout/angular-leaflet-directive/issues/events{/number}":               struct{}{},
		"https://api.github.com/repos/engineyard/pundit/issues/events{/number}":                                 struct{}{},
		"https://api.github.com/repos/puppetlabs/csr-generator/issues/events{/number}":                          struct{}{},
		"https://api.github.com/repos/keradgames/keradbot/issues/events{/number}":                               struct{}{},
		"https://api.github.com/repos/DogFoodSoftware/test-repo/issues/events{/number}":                         struct{}{},
		"https://api.github.com/repos/pluspole/PlainText/issues/events{/number}":                                struct{}{},
		"https://api.github.com/repos/USU-Robosub/mk-proto/issues/events{/number}":                              struct{}{},
		"https://api.github.com/repos/AgreeMates/AgreeMates.github.io/issues/events{/number}":                   struct{}{},
		"https://api.github.com/repos/OffTempo/offtempoV4/issues/events{/number}":                               struct{}{},
		"https://api.github.com/repos/AgreeMates/AgreeMates/issues/events{/number}":                             struct{}{},
		"https://api.github.com/repos/colonyamerican/node-googlemaps/issues/events{/number}":                    struct{}{},
		"https://api.github.com/repos/mulesoft/template-sfdc2sfdc-account-broadcast/issues/events{/number}":     struct{}{},
		"https://api.github.com/repos/mulesoft/template-sfdc2sfdc-contact-aggregation/issues/events{/number}":   struct{}{},
		"https://api.github.com/repos/pbdesk/EntityFramework.HandsOn/issues/events{/number}":                    struct{}{},
		"https://api.github.com/repos/godaddy/node-soap-client-utils/issues/events{/number}":                    struct{}{},
		"https://api.github.com/repos/FlorianSoftware/SmartApps-Open/issues/events{/number}":                    struct{}{},
	}
	b := NewBuilder()
	for str, _ := range m {
		b.Add(str)
	}
	w := b.Build()
	Convey("When a random string is generated", t, func() {
		for str, _ := range m {
			ws, ok := w.LongestPrefixMatch(str)
			So(ok, ShouldBeTrue)
			So(ws.Len, ShouldEqual, len(str))
		}
	})
}

func TestMarshallingWX(t *testing.T) {
	num := 1000000
	maxLen := 10
	testNum := 10
	wxb := NewBuilder()
	totalLen := 0
	strs := make(map[string]struct{})
	for i := 0; i < num; i++ {
		l := rand.Int() % maxLen
		strbuf := make([]byte, l)
		for j := 0; j < l; j++ {
			strbuf[j] = byte(rand.Int() % 4)
		}
		totalLen += l
		s := string(strbuf)
		strs[s] = struct{}{}
		wxb.Add(s)
	}
	w := wxb.Build()

	strarray := make([]string, 0)
	for str, _ := range strs {
		strarray = append(strarray, str)
	}
	Convey("When Get is examined", t, func() {
		So(w.Num(), ShouldEqual, uint64(len(strs)))
		for i := 0; i < testNum; i++ {
			ind := uint64(rand.Int31n(int32(w.Num())))
			s := w.Get(uint64(ind))
			_, ok := strs[s]
			So(ok, ShouldBeTrue)
		}
	})
	Convey("When LongestPrefixMatch is examined", t, func() {
		So(w.Num(), ShouldEqual, uint64(len(strs)))
		for i := 0; i < testNum; i++ {
			ind := uint64(rand.Int31n(int32(w.Num())))
			ws, ok := w.LongestPrefixMatch(strarray[ind])
			So(ok, ShouldBeTrue)
			So(ws.Len, ShouldEqual, len(strarray[ind]))
		}
	})

	Convey("When large strings are set", t, func() {
		out, err := w.MarshalBinary()
		So(err, ShouldBeNil)
		wxnew := New()
		err = wxnew.UnmarshalBinary(out)
		So(err, ShouldBeNil)
		So(wxnew.Num(), ShouldEqual, uint64(len(strs)))
		for i := 0; i < testNum; i++ {
			ind := uint64(rand.Int31n(int32(wxnew.Num())))
			s := wxnew.Get(uint64(ind))
			_, ok := strs[s]
			So(ok, ShouldBeTrue)
		}
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
