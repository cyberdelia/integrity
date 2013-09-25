package integrity

import (
	"os"
	"testing"
)

func TestVerifyCorrect(t *testing.T) {
	correct, err := os.Open("testdata/correct")
	if err != nil {
		t.Fatalf("Can't open file: %s", err)
	}
	err = Verify(correct)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestVerifyCorrupt(t *testing.T) {
	corrupt, err := os.Open("testdata/corrupt")
	if err != nil {
		t.Fatalf("Can't open file: %s", err)
	}
	err = Verify(corrupt)
	if err != ErrChecksum {
		t.Fatalf("Unexpected error: %s", err)
	}
}
