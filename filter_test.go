package cuckoo

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"testing"
)

func getHashKey(key uint64) uint64 {
	h := fnv.New64a()
	var bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, key)
	h.Write(bs)

	return h.Sum64()
}

func TestNew(t *testing.T) {
	c := New(4)

	var key, i, n uint64
	n = 1000
	for key = 1; key < n; key++ {
		hashKey := getHashKey(key)
		if ok := c.Add(hashKey); ok {
			fmt.Printf("key: %d, num: %d\n", key, c.num)
		} else {
			fmt.Println("===> expand...\n")
			c = c.Expand(1)
			c.Add(hashKey)
		}

		for i = 1; i < n; i++ {
			iKey := getHashKey(i)
			if i <= key && !c.Find(iKey) {
				t.Fatalf("error i=%d", i)
			}
			if i > key && c.Find(iKey) {
				t.Fatalf("error i=%d", i)
			}
		}
	}

}
