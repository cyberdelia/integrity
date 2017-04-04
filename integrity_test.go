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
