package wx

import (
	"fmt"
	"math"

	"github.com/hillbig/fixvec"
	"github.com/hillbig/rsdic"
	"github.com/hillbig/vecstring"

	"github.com/ugorji/go/codec"
)

type wxImpl struct {
	branches            vecstring.VecStringForWX
	terminals           rsdic.RSDic
	leadingIDs          fixvec.FixVec
	leadings            vecstring.VecString
	leadingIsZeroOrOnes rsdic.RSDic
	leadingIsOnes       rsdic.RSDic
	leadingOnes         []byte
	num                 uint64
}

func (wx wxImpl) ExactMatch(str string) (id uint64, found bool) {
	ret, found := wx.LongestPrefixMatch(str)
	if found == true && int(ret.Len) == len(str) {
		return ret.ID, found
	} else {
		return 0, false
	}
}

func (wx wxImpl) LongestPrefixMatch(str string) (WXString, bool) {
	rets := wx.CommonPrefixMatch(str)
	if len(rets) > 0 {
		return rets[len(rets)-1], true
	} else {
		return WXString{}, false
	}
}

func (wx wxImpl) CommonPrefixMatch(str string) []WXString {
	return wx.CommonPrefixMatchWithLimit(str, math.MaxUint64)
}

func (wx wxImpl) CommonPrefixMatchWithLimit(str string, limit uint64) (rets []WXString) {
	rets = make([]WXString, 0)
	if limit == 0 {
		return
	}
	nodeind := uint64(0)
	strind := uint64(0)
	strlen := uint64(len(str))
	for {
		l, ok := wx.prefixMatchLeading(nodeind, str[strind:])
		strind += l
		if !ok {
			return
		}
		ok, id := wx.terminals.BitAndRank(nodeind)
		if ok {
			rets = append(rets, WXString{id, strind})
			if uint64(len(rets)) == limit {
				return
			}
		}
		if strind == strlen {
			return
		}
		nodeind, ok = wx.branches.FindZeroRank(nodeind, str[strind])
		if !ok {
			return
		}
		nodeind++
		strind++
	}
	// don't come here
}

func (wx wxImpl) getLeading(nodeind uint64) string {
	isZeroOrOne, rank := wx.leadingIsZeroOrOnes.BitAndRank(nodeind)
	if isZeroOrOne {
		isOne, rank := wx.leadingIsOnes.BitAndRank(rank)
		if isOne {
			s := make([]byte, 1)
			s[0] = wx.leadingOnes[rank]
			return string(s)
		} else {
			return ""
		}
	} else {
		id := wx.leadingIDs.Get(rank)
		return wx.leadings.Get(id)
	}
}

func (wx wxImpl) prefixMatchLeading(nodeind uint64, str string) (uint64, bool) {
	isZeroOrOne, rank := wx.leadingIsZeroOrOnes.BitAndRank(nodeind)
	if isZeroOrOne {
		isOne, rank := wx.leadingIsOnes.BitAndRank(rank)
		if isOne {
			if len(str) > 0 && wx.leadingOnes[rank] == str[0] {
				return 1, true
			} else {
				return 0, false
			}
		} else {
			return 0, true
		}
	} else {
		id := wx.leadingIDs.Get(rank)
		return wx.leadings.PrefixMatch(id, str)
	}
}

func (wx wxImpl) PredictiveMatch(str string) (ids []uint64) {
	return wx.PredictiveMatchWithLimit(str, math.MaxUint64)
}

func (wx wxImpl) PredictiveMatchWithLimit(str string, limit uint64) (ids []uint64) {
	ids = make([]uint64, 0)
	if limit == 0 {
		return
	}
	nodeind := uint64(0)
	strind := uint64(0)
	strlen := uint64(len(str))
	for {
		l, ok := wx.prefixMatchLeading(nodeind, str[strind:])
		strind += l
		if !ok {
			if strind != strlen {
				return
			} else {
				break
			}
		}
		if strind == strlen {
			break
		}
		nodeind, ok = wx.branches.FindZeroRank(nodeind, str[strind])
		if !ok {
			return
		}
		nodeind++
		strind++
	}
	// strind == strlen
	wx.enumerateDescendant(nodeind, limit, &ids)
	return
}

func (wx wxImpl) enumerateDescendant(nodeind uint64, limit uint64, ids *[]uint64) {
	ok, id := wx.terminals.BitAndRank(nodeind)
	if ok {
		*ids = append(*ids, id)
		if uint64(len(*ids)) == limit {
			return
		}
	}
	l, offset := wx.branches.LenAndOffset(nodeind)
	for i := uint64(0); i < l; i++ {
		wx.enumerateDescendant(offset+i+1, limit, ids)
		if uint64(len(*ids)) == limit {
			break
		}
	}
}

func (wx wxImpl) Get(id uint64) string {
	nodeind := wx.terminals.Select(id, true)
	strs := make([]string, 0)
	for {
		str := wx.getLeading(nodeind)
		strs = append(strs, str)
		if nodeind == 0 {
			break
		}
		s := make([]byte, 1)
		s[0] = wx.branches.GetByte(nodeind - 1)
		strs = append(strs, string(s))
		nodeind = wx.branches.IthCharInd(nodeind - 1)
	}

	ret := make([]byte, 0)
	for i := len(strs) - 1; i >= 0; i-- {
		ret = append(ret, []byte(strs[i])...)
	}
	return string(ret)
}

func (wx wxImpl) Num() uint64 {
	return wx.num
}

func (wx wxImpl) MarshalBinary() (out []byte, err error) {
	var bh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&out, &bh)
	err = enc.Encode(wx.branches)
	if err != nil {
		return
	}
	err = enc.Encode(wx.terminals)
	if err != nil {
		return
	}
	err = enc.Encode(wx.leadingIDs)
	if err != nil {
		return
	}
	err = enc.Encode(wx.leadings)
	if err != nil {
		return
	}
	err = enc.Encode(wx.leadingIsZeroOrOnes)
	if err != nil {
		return
	}
	err = enc.Encode(wx.leadingIsOnes)
	if err != nil {
		return
	}
	err = enc.Encode(wx.leadingOnes)
	if err != nil {
		return
	}
	err = enc.Encode(wx.num)
	if err != nil {
		return
	}
	return
}

func (wx *wxImpl) UnmarshalBinary(in []byte) (err error) {
	var bh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(in, &bh)
	err = dec.Decode(&wx.branches)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.terminals)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.leadingIDs)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.leadings)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.leadingIsZeroOrOnes)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.leadingIsOnes)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.leadingOnes)
	if err != nil {
		return
	}
	err = dec.Decode(&wx.num)
	if err != nil {
		return
	}
	return nil
}

func debugPrintVecString(name string, vs vecstring.VecString) {
	fmt.Printf("%s num=%d\n", name, vs.Num())
	num := vs.Num()
	for i := uint64(0); i < num; i++ {
		s := vs.Get(i)
		fmt.Printf("%d:%s\n", i, s)
	}
}

func debugPrintRSDic(name string, rs rsdic.RSDic) {
	fmt.Printf("%s num=%d\n", name, rs.Num())
	num := rs.Num()
	for i := uint64(0); i < num; i++ {
		if rs.Bit(i) {
			fmt.Printf("1")
		} else {
			fmt.Printf("0")
		}
	}
	fmt.Printf("\n")
}

func debugPrintFixVec(name string, fv fixvec.FixVec) {
	fmt.Printf("%s num=%d blen=%d\n", name, fv.Num(), fv.Blen())
	num := fv.Num()
	for i := uint64(0); i < num; i++ {
		fmt.Printf("%d ", fv.Get(i))
	}
	fmt.Printf("\n")
}

func debugPrintBytes(name string, bytes []byte) {
	fmt.Printf("%s num=%d\n", name, len(bytes))
	for _, v := range bytes {
		fmt.Printf("%c ", v)
	}
	fmt.Printf("\n")
}

func (wx wxImpl) debugPrint() {
	fmt.Printf("num=%d\n", wx.num)
	debugPrintVecString("branches", wx.branches)
	debugPrintRSDic("terminals", wx.terminals)
	debugPrintFixVec("leadingIDs", wx.leadingIDs)
	debugPrintVecString("leadings", wx.leadings)
	debugPrintRSDic("leadingIsZeroOrOnes", wx.leadingIsZeroOrOnes)
	debugPrintRSDic("leadingIsOnes", wx.leadingIsOnes)
	debugPrintBytes("leadingOnes", wx.leadingOnes)
}
