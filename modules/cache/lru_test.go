package cache

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func ExampleLRU_Add() {
	lru := newLRU(2)

	lru.Add(1, []byte("foo"))
	lru.Add(2, []byte("bar"))
	lru.Add(3, []byte("baz"))
	lru.debug()

	// Output:
	// 3 baz
	// 2 bar
}

func ExampleLRU_Get() {
	lru := newLRU(2)

	lru.Add(1, []byte("foo"))
	lru.Add(2, []byte("bar"))
	lru.Add(3, []byte("baz"))

	fmt.Println(string(lru.Get(4)))
	fmt.Println(string(lru.Get(3)))

	// Output:
	// baz
}

func ExampleLRU_Remove() {
	lru := newLRU(2)

	lru.Add(1, []byte("foo"))
	lru.Add(2, []byte("bar"))
	lru.Add(3, []byte("baz"))
	lru.Remove(2)
	lru.Remove(3)
	lru.debug()

	// Output:
}

func Benchmark_LRU_Add(b *testing.B) {
	cases := []int{2, 50, 100, 500}
	s := bytes.Repeat([]byte(" "), 1e4)

	for k, c := range cases {
		b.Run(strconv.Itoa(k), func(b *testing.B) {
			lru := newLRU(c)
			for i := 0; i < b.N; i++ {
				lru.Add(i, s)
			}
			for i := 0; i < b.N; i++ {
				lru.Get(i)
			}
			for i := 0; i < b.N; i++ {
				lru.Remove(i)
			}
		})
	}
}
