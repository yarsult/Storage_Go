package storage

import "testing"

type pieceOfTest struct {
	key   string
	value string
}

func TestSetGet(t *testing.T) {
	pieces := []pieceOfTest{
		{"vsem", "privet"},
		{"testing", "tests"},
		{"go", "golang"},
	}
	stor, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}
	for _, p := range pieces {
		t.Run("1", func(t *testing.T) {
			stor.Set(p.key, p.value)
			res, _ := stor.Get(p.key)
			if res != p.value {
				t.Errorf("not equal values")
			}
		})
	}
}

type pieceOfTestWithKind struct {
	key   string
	value string
	kind  string
}

func TestKind(t *testing.T) {
	pieces := []pieceOfTestWithKind{
		{"vsem", "privet", "S"},
		{"testing", "tests", "S"},
		{"go", "45678", "D"},
	}
	stor, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}
	for _, p := range pieces {
		t.Run("2", func(t *testing.T) {
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
		_, _ = stor.Get("vsem")
		_, _ = stor.Get("testing")
		_, _ = stor.Get("go")
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
		_, _ = stor.Get("vsem")
		_, _ = stor.Get("testing")
		_, _ = stor.Get("go")
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
