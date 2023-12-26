package main

import (
	"fmt"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
)

func main() {
	r := iradix.New[int]()
	r, _, _ = r.Insert([]byte("foo"), 1)
	r, _, _ = r.Insert([]byte("fooaa"), 2)
	r, _, _ = r.Insert([]byte("fooaabb"), 3)
	fmt.Println(r.Get([]byte("fooaa")))                   //输出2
	s, i, e := r.Root().LongestPrefix([]byte("fooaa123")) //输出2, 因为最长前缀匹配到了fooaa
	fmt.Println(string(s), i, e)
	s, i, e = r.Root().Maximum()
	fmt.Println(string(s), i, e) //输出3, 因为fooaabb是key最长的
	it := r.Root().Iterator()
	it.SeekPrefix([]byte("fooa")) //将指针定到fooa前缀
	for key, _, ok := it.Next(); ok; key, _, ok = it.Next() {
		fmt.Println("===", string(key))
	}
	// 输出
	// === fooaa
	// === fooaabb
	it = r.Root().Iterator()
	it.SeekLowerBound([]byte("foo"))
	for key, _, ok := it.Next(); ok; key, _, ok = it.Next() {
		fmt.Println("+++", string(key))
	}
	// 输出
	// +++ foo
	// +++ fooaa
	// +++ fooaabb

	tn := r.Txn() //启动事务，但要注意这个事务不是线程安全的，只能在一个goroutine中使用
	fmt.Println(tn.Delete([]byte("fooaa")))
	r = tn.Commit()

	it = r.Root().Iterator()
	it.SeekLowerBound([]byte("foo"))
	for key, _, ok := it.Next(); ok; key, _, ok = it.Next() {
		fmt.Println("---", string(key))
	}
	// 输出
	// --- foo
	// --- fooaabb
}
