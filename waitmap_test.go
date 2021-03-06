package waitmap

import (
	"testing"
)

func TestFlatGet(t *testing.T) {
	m := New()
	m.Set(1, 2)
	m.Set(2, 3)
	if m.Get(1) != 2 {
		t.Errorf("Wrong result from Get(1)")
	}
	if m.Get(2) != 3 {
		t.Errorf("Wrong result from Get(2)")
	}
}

func TestFlatCheck(t *testing.T) {
	m := New()
	if m.Check(1) {
		t.Errorf("Check(1) returned true; false expected")
	}

	m.Set(1, 2)
	if !m.Check(1) {
		t.Errorf("Check(1) returned false; true expected")
	}
	if m.Check(2) {
		t.Errorf("Check(2) returned true; false expected")
	}
}

func TestSimple(t *testing.T) {
	m := New()
	var v interface{}
	c := make(chan interface{})
	go func() {
		v = m.Get(1)
		c <- nil
	}()
	go func() {
		m.Set(1, "hi")
	}()
	<-c
	if v != "hi" {
		t.Errorf("m.Get(1) returned incorrect valeu")
	}
}

func TestMultipleSet(t *testing.T) {
	m := New()
	a := m.Set(1, "HI")
	b := m.Set(1, "HI")
	if a != true {
		t.Errorf("m.Set(1) inappropriately returned false")
	}
	if b == true {
		t.Errorf("m.Set(1) inappropriately returned true")
	}
}

func BenchmarkRawMapGet(b *testing.B) {
	m := map[string]string{
		"foo": "bar",
	}
	for i := 0; i < b.N; i++ {
		_, _ = m["foo"]
	}
}

func BenchmarkFlatGoodCheck(b *testing.B) {
	m := New()
	m.Set(0, 0)
	for i := 0; i < b.N; i++ {
		m.Check(0)
	}
}

func BenchmarkFlatFailedCheck(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Check(0)
	}
}

func BenchmarkFlatWaitedGet(b *testing.B) {
	m := New()
	go func() { m.Get(0) }()
	m.Set(0, "ho")
	for i := 0; i < b.N; i++ {
		m.Get(0)
	}
}

func BenchmarkFlatSimpleGet(b *testing.B) {
	m := New()
	m.Set(0, "ho")
	for i := 0; i < b.N; i++ {
		m.Get(0)
	}
}

func BenchmarkRawSimpleSet(b *testing.B) {
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		m[0] = i
	}
}

func BenchmarkRawIncrementalSet(b *testing.B) {
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}

func BenchmarkFlatSimpleSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(0, i)
	}
}

func BenchmarkFlatIncrementalSet(b *testing.B) {
	m := New()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}
