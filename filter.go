package cuckoo

import ()

const (
	keyMask = 1<<32 - 1

	loopMax = 20

	bitsMax = 32

	Empty uint64 = 0
)

type Cuckoo struct {
	num      uint
	size     uint64
	mask     uint64
	buckets  []uint64
	buckets2 []uint64
}

func New(bitsNum uint) *Cuckoo {
	size := uint64(1 << bitsNum)
	return &Cuckoo{
		num:      bitsNum,
		size:     size,
		mask:     size - 1,
		buckets:  make([]uint64, size),
		buckets2: make([]uint64, size),
	}
}

func (c *Cuckoo) Copy() *Cuckoo {
	return c
}

func (c *Cuckoo) Find(hashKey uint64) (isExist bool) {
	if c.buckets[c.key(hashKey)] == hashKey || c.buckets2[c.key2(hashKey)] == hashKey {
		return true
	}
	return isExist
}

// 增加一个元素
// 如果不成功，则返回false，这是需要先扩容
func (c *Cuckoo) Add(hashKey uint64) bool {
	for i := 0; i < loopMax; i++ {
		if hashKey = c.change(hashKey); hashKey == 0 {
			return true
		}
	}

	// 需要扩容
	return false
}

func (c *Cuckoo) Delete(hashKey uint64) {
	key := c.key(hashKey)
	if c.buckets[key] == hashKey {
		c.buckets[key] = Empty
		return
	}

	key2 := c.key(hashKey)
	if c.buckets2[key2] == hashKey {
		c.buckets2[key2] = Empty
	}
}

// 扩容
// TODO 是否会出现需要连续扩容多次才能满足要求？
func (c *Cuckoo) Expand(step uint) *Cuckoo {
	if step > 2 || c.num+step > bitsMax {
		panic("error in Expand")
	}

	new := New(c.num + step)

	var ok bool
	for _, key := range c.buckets {
		if key != Empty {
			if ok = new.Add(key); !ok {
				break
			}
		}
	}

	if ok {
		for _, key := range c.buckets2 {
			if key != Empty {
				if ok = new.Add(key); !ok {
					break
				}
			}
		}
	}

	if !ok {
		return c.Expand(step + 1)
	}

	return new
}

//****************** Private ***************************

func (c *Cuckoo) key(hashKey uint64) uint64 {
	return hashKey & c.mask
}

func (c *Cuckoo) key2(hashKey uint64) uint64 {
	return (hashKey >> 32) & c.mask
}

func (c *Cuckoo) change(hashKey uint64) (changingKey uint64) {
	key := c.key(hashKey)
	if c.buckets[key] == Empty {
		c.buckets[key] = hashKey
		return changingKey
	}
	key2 := c.key2(hashKey)
	if c.buckets2[key2] == Empty {
		c.buckets2[key2] = hashKey
		return changingKey
	}

	// 踢出第一个
	//println("change")
	changingKey = c.buckets[key]
	c.buckets[key] = hashKey
	return changingKey
}
