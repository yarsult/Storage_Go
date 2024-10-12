package storage

import "testing"

type pieceOfTest struct {
	name  string
	key   string
	value string
}

func TestSetGet(t *testing.T) {
	pieces := []pieceOfTest{
		{"1st test", "vsem", "privet"},
		{"2nd test", "testing", "tests"},
		{"3rd test", "go", "golang"},
	}
	stor, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}
	for _, p := range pieces {
		t.Run(p.name, func(t *testing.T) {
			stor.Set(p.key, p.value)
			if *stor.Get(p.key) != p.value {
				t.Errorf("not equal values")
			}
		})
	}
}

type pieceOfTestWithKind struct {
	name  string
	key   string
	value string
	kind  string
}

func TestKind(t *testing.T) {
	pieces := []pieceOfTestWithKind{
		{"1st test", "vsem", "privet", "S"},
		{"2nd test", "testing", "tests", "S"},
		{"3rd test", "go", "45678", "D"},
	}
	stor, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}
	for _, p := range pieces {
		t.Run(p.name, func(t *testing.T) {
			stor.Set(p.key, p.value)
			if stor.GetKind(p.key) != p.kind {
				t.Errorf("wrong kind")
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	stor, err := NewStorage()
	if err != nil {
		b.Fatalf("not able to create storage: %v", err)
	}
	stor.Set("vsem", "privet")
	stor.Set("testing", "tests")
	stor.Set("go", "45678")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stor.Get("vsem")
		_ = stor.Get("testing")
		_ = stor.Get("go")
	}
}

func BenchmarkSet(b *testing.B) {
	stor, err := NewStorage()
	if err != nil {
		b.Fatalf("not able to create storage: %v", err)
	}
	for i := 0; i < b.N; i++ {
		stor.Set("vsem", "privet")
		stor.Set("testing", "tests")
		stor.Set("go", "45678")
	}
}

func BenchmarkSetGet(b *testing.B) {
	stor, err := NewStorage()
	if err != nil {
		b.Fatalf("not able to create storage: %v", err)
	}
	stor.Set("vsem", "privet")
	stor.Set("testing", "tests")
	stor.Set("go", "45678")
	for i := 0; i < b.N; i++ {
		_ = stor.Get("vsem")
		_ = stor.Get("testing")
		_ = stor.Get("go")
	}
}

func BenchmarkGetKind(b *testing.B) {
	stor, err := NewStorage()
	if err != nil {
		b.Fatalf("not able to create storage: %v", err)
	}
	stor.Set("vsem", "privet")
	stor.Set("testing", "tests")
	stor.Set("go", "45678")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stor.GetKind("vsem")
		_ = stor.GetKind("testing")
		_ = stor.GetKind("go")
	}
}
