package metadata

import (
	"container/list"
	"errors"
)

type Node struct {
	Data   string
	KeyPtr *list.Element
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
	} else {
		item.Data = value
		l.Items[key] = item
		l.Queue.MoveToFront(item.KeyPtr)
	}

}

//cache := &metadata.Cache
//cache.Put("Content-Length", w.Header().Get("Content-Length"))
//fmt.Println(cache.Get("Valod"))
