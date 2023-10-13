package passwd

import (
	"testing"
)

func TestBcryptImpl(t *testing.T) {
	impls := []HashComparer{
		BcryptDefaultCost,
	}

	const plain = "abc123"

	for _, impl := range impls {
		hash, err := impl.Hash(plain)
		if err != nil {
			t.Fatalf("expect no error; got an error: %v", err)
		}

		if err := impl.Compare(hash, plain); err != nil {
			t.Errorf("expect password is match")
		}
	}
}
