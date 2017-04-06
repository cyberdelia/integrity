package integrity

import (
	"os"
	"testing"
)

func TestVerifyCorrect(t *testing.T) {
	correct, err := os.Open("testdata/correct")
	if err != nil {
		t.Fatalf("can't open file: %s", err)
	}
	if err = Verify(correct); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestVerifyCorrupt(t *testing.T) {
	corrupt, err := os.Open("testdata/corrupt")
	if err != nil {
		t.Fatalf("can't open file: %s", err)
	}
	if err = Verify(corrupt); err != ErrChecksum {
		t.Fatalf("unexpected error: %s", err)
	}
}

func BenchmarkVerify(b *testing.B) {
	b.ReportAllocs()
	p, err := os.Open("testdata/correct")
	if err != nil {
		b.Fatalf("can't open file: %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Seek(0, 0)
		if err := Verify(p); err != nil {
			b.Fatalf("unexpected error: %s", err)
		}
	}
}
