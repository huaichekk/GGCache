package consistent

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewConsistent(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.AddNode("6") //06 16 26
	hash.AddNode("4") //04 14 24
	hash.AddNode("2") //02 12 22

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.ChooseNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.AddNode("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.ChooseNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	hash.AddNode("100")
	hash.DeleteNode("100")
}
