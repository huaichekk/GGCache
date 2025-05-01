package consistent

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(key []byte) uint32
type Consistent struct {
	ring          []int
	replicas      int
	virtualToNode map[int]string
	addrs         []string
	hash          Hash
}

func NewConsistent(replicas int, hash Hash) *Consistent {
	var h Hash
	if hash == nil {
		h = crc32.ChecksumIEEE
	}
	return &Consistent{
		ring:          make([]int, 0),
		replicas:      replicas,
		virtualToNode: make(map[int]string),
		hash:          h,
	}
}

func (c *Consistent) AddNode(addr string) {
	for i := 0; i < c.replicas; i++ {
		virtualNode := int(c.hash([]byte(strconv.Itoa(i) + addr)))
		c.ring = append(c.ring, virtualNode)
		c.virtualToNode[virtualNode] = addr
	}
	c.addrs = append(c.addrs, addr)
	sort.Ints(c.ring)
	fmt.Println("add node :", c.addrs, c.ring)
}

func (c *Consistent) DeleteNode(addr string) {
	for i := 0; i < c.replicas; i++ {
		virtualNode := int(c.hash([]byte(strconv.Itoa(i) + addr)))
		index := sort.SearchInts(c.ring, virtualNode)
		if index == len(c.ring)-1 {
			c.ring = c.ring[:index]
		} else {
			c.ring = append(c.ring[:index], c.ring[index+1:]...)
		}
		delete(c.virtualToNode, virtualNode)
	}
	for k, v := range c.addrs {
		if v == addr {
			if k == len(c.addrs)-1 {
				c.addrs = c.addrs[:k]
			} else {
				c.addrs = append(c.addrs[:k], c.addrs[k+1:]...)
			}
			break
		}
	}
	fmt.Println("delete node", c.addrs, c.ring)
}

func (c *Consistent) ChooseNode(key string) string {
	if key == "" {
		return ""
	}
	keyHash := int(c.hash([]byte(key)))
	index := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= keyHash
	})
	return c.virtualToNode[c.ring[index%len(c.ring)]]
}
