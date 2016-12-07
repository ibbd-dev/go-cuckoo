package cuckoo

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"testing"
)

var testHashFunc func(uint64) uint64

func fnv64aHashKey(key uint64) uint64 {
	h := fnv.New64a()
	var bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, key)
	h.Write(bs)

	return h.Sum64()
}

func fnv64HashKey(key uint64) uint64 {
	h := fnv.New64()
	var bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, key)
	h.Write(bs)

	return h.Sum64()
}

func md5HashKey(key uint64) uint64 {
	h := md5.New()
	var bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, key)
	res := h.Sum(bs)
	return binary.BigEndian.Uint64(res)
}

func TestNew(t *testing.T) {
	c := New(4)
	testHashFunc = fnv64aHashKey
	testHashFunc = md5HashKey
	testHashFunc = fnv64HashKey

	var key, i, n uint64
	n = 1500000
	for key = 1; key < n; key++ {
		hashKey := testHashFunc(key)
		if ok := c.Insert(hashKey); !ok {
			fmt.Println("===> expand...\n")
			c, ok = c.Expand(1, hashKey)
			if !ok {
				t.Fatalf("too many keys")
			}
			fmt.Printf("key: %d, num: %d\n", key, c.num)
		}

		for i = 0; i < 100; i++ {
			if n-i > 0 {
				iKey := testHashFunc(n - i)
				if n-i <= key && !c.Lookup(iKey) {
					t.Fatalf("error n-i=%d", n-i)
				}
				if n-i > key && c.Lookup(iKey) {
					t.Fatalf("error n-i=%d", n-i)
				}
			}
		}
	}
}
