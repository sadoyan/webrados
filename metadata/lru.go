package metadata

import (
	"container/list"
	"errors"
	"fmt"
	"time"
)

type Node struct {
	Data   string
	KeyPtr *list.Element
	Date   time.Time
}

type LRUCache struct {
	Queue    *list.List
	Items    map[string]*Node
	Capacity int
}

func CacheConstructor(capacity int) LRUCache {
	return LRUCache{Queue: list.New(), Items: make(map[string]*Node), Capacity: capacity}
}

func (l *LRUCache) Get(key string) (string, error) {
	if item, ok := l.Items[key]; ok {
		l.Queue.MoveToFront(item.KeyPtr)
		return item.Data, nil
	}
	return "", errors.New("No item")
}

func (l *LRUCache) Size() int {
	return len(l.Items)
}

func (l *LRUCache) Put(key string, value string) {
	if item, ok := l.Items[key]; !ok {
		if l.Capacity == len(l.Items) {
			back := l.Queue.Back()
			l.Queue.Remove(back)
			delete(l.Items, back.Value.(string))
		}
		l.Items[key] = &Node{Data: value, KeyPtr: l.Queue.PushFront(key)}
		l.Items[key].Date = time.Now()
	} else {
		l.Items[key] = item
		l.Queue.MoveToFront(item.KeyPtr)
	}

}

func (l *LRUCache) Evict() {
	for {
		start := time.Now()
		num := 0
		dif := 300 * time.Second
		then := time.Now().Add(-dif)
		for x := range l.Items {
			dd := l.Items[x].Date
			if dd.Before(then) {
				num = num + 1
				back := l.Queue.Back()
				l.Queue.Remove(back)
				delete(l.Items, back.Value.(string))
			}
		}
		fmt.Println(" ")
		fmt.Println("Evixted element: ", num, "current size:", l.Size(), len(l.Items))
		fmt.Printf("Execution time %s\n", time.Since(start))
		num = 0
		time.Sleep(5 * time.Second)
	}
}
