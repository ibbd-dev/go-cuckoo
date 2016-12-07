package cuckoo

import ()

const (
	// 最大的踢出次数
	kickOutMax = 500

	// 最大的位数
	// 其所表示的空间为：1<<32
	bitsMax = 32

	// 空桶状态
	Empty uint64 = 0
)

type Cuckoo struct {
	num    uint
	size   uint64
	mask   uint64
	table  []uint64
	table2 []uint64
}

func New(bitsNum uint) *Cuckoo {
	size := uint64(1 << bitsNum)
	return &Cuckoo{
		num:    bitsNum,
		mask:   size - 1,
		table:  make([]uint64, size), // hash table，对应c.key()
		table2: make([]uint64, size), // hash table2，对应c.key2()
	}
}

func (c *Cuckoo) Copy() *Cuckoo {
	return c
}

func (c *Cuckoo) Lookup(hashKey uint64) (isExist bool) {
	if c.table[c.key(hashKey)] == hashKey || c.table2[c.key2(hashKey)] == hashKey {
		return true
	}
	return isExist
}

// 增加一个元素
// 如果不成功，则返回false，这是需要先扩容
func (c *Cuckoo) Insert(hashKey uint64) bool {
	for i := 0; i < kickOutMax; i++ {
		if hashKey = c.relocate(hashKey); hashKey == 0 {
			return true
		}
	}

	// 需要扩容
	return false
}

func (c *Cuckoo) Delete(hashKey uint64) {
	key := c.key(hashKey)
	if c.table[key] == hashKey {
		c.table[key] = Empty
		return
	}

	key2 := c.key(hashKey)
	if c.table2[key2] == hashKey {
		c.table2[key2] = Empty
	}
}

// 扩容
// TODO 是否会出现需要连续扩容多次才能满足要求？
// ok 如果为false，则表示插入新值失败，已经到达了空间的上限
func (c *Cuckoo) Expand(step uint, newHashKey uint64) (new *Cuckoo, ok bool) {
	if c.num+step > bitsMax {
		return new, ok
		//panic("error in Expand")
	}

	new = New(c.num + step)
	new.Insert(newHashKey)

	for _, key := range c.table {
		if key != Empty {
			if ok = new.Insert(key); !ok {
				break
			}
		}
	}

	if ok {
		for _, key := range c.table2 {
			if key != Empty {
				if ok = new.Insert(key); !ok {
					break
				}
			}
		}
	}

	if !ok {
		print("=====> Expand more then once")
		return c.Expand(step+1, newHashKey)
	}

	return new, ok
}

//****************** Private ***************************

func (c *Cuckoo) key(hashKey uint64) uint64 {
	return hashKey & c.mask
}

func (c *Cuckoo) key2(hashKey uint64) uint64 {
	return (hashKey >> 32) & c.mask
}

func (c *Cuckoo) relocate(hashKey uint64) (changingKey uint64) {
	key := c.key(hashKey)
	if c.table[key] == Empty {
		c.table[key] = hashKey
		return changingKey
	}
	key2 := c.key2(hashKey)
	if c.table2[key2] == Empty {
		c.table2[key2] = hashKey
		return changingKey
	}

	// 踢出第一个
	//println("relocate")
	changingKey = c.table[key]
	c.table[key] = hashKey
	return changingKey
}
