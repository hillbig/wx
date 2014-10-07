package wx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQueue(t *testing.T) {
	q := newNodeQueue(0)
	Convey("When a empty queue is examined", t, func() {
		So(q.num, ShouldEqual, 0)
		So(q.beg, ShouldEqual, 0)
	})
	q.push(&node{0, 1, 2})
	Convey("When a queue with one node is examined", t, func() {
		So(q.num, ShouldEqual, 1)
	})
	n := q.pop()
	Convey("When a value is popped", t, func() {
		So(n.beg, ShouldEqual, 0)
		So(n.end, ShouldEqual, 1)
		So(n.depth, ShouldEqual, 2)
	})
	q.push(&node{3, 4, 5})
	n = q.pop()
	Convey("When a value is popped", t, func() {
		So(n.beg, ShouldEqual, 3)
		So(n.end, ShouldEqual, 4)
		So(n.depth, ShouldEqual, 5)
		q.push(&node{1, 2, 3})
		q.push(&node{1, 2, 3})
		q.push(&node{4, 2, 3})
		n = q.pop()
		n = q.pop()
		n = q.pop()
		So(n.beg, ShouldEqual, 4)
	})
}

func TestManyQueue(t *testing.T) {
	q := newNodeQueue(0)
	for i := 0; i < 100; i++ {
		q.push(&node{i, i + 1, i + 2})
	}
	Convey("When a queue with 100 nodes is examined", t, func() {
		So(q.num, ShouldEqual, 100)
		for i := 0; i < 100; i++ {
			n := q.pop()
			So(n.beg, ShouldEqual, i)
		}
		So(q.num, ShouldEqual, 0)
	})
}
