package web

import (
	"encoding/json"
	"metadata"
	"runtime"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

var m runtime.MemStats

type metrics struct {
	Alloc        uint64 `json:"alloc,int"`
	Total        uint64 `json:"total,int"`
	System       uint64 `json:"memtotal,int"`
	Mallocs      uint64 `json:"malloc,int"`
	Gcnum        uint32 `json:"gcnum,int"`
	NextGC       uint64 `json:"nextgc,int"`
	Frees        uint64 `json:"frees,int"`
	HeapAlloc    uint64 `json:"heapalloc,int"`
	HeapIdle     uint64 `json:"heapidle,int"`
	HeapInuse    uint64 `json:"heapinuse,int"`
	HeapObjects  uint64 `json:"heapobjects,int"`
	HeapReleased uint64 `json:"heapreleased,int"`
	//LastGC        uint64 `json:"lastgc,int"`
	NumForcedGC  uint32 `json:"forcegc,int"`
	PauseTotalNs uint64 `json:"pausetotal,int"`
	Goroutines   int    `json:"goroutines,int"`
	GetCount     int    `json:"getcount,int"`
	PostCount    int    `json:"postcount,int"`
	DelCount     int    `json:"delcount,int"`
	//HeadCount     int    `json:"headcount,int"`
	CacheLen      int   `json:"cacheitems,int"`
	CacheCapacity int   `json:"cachesize,int"`
	CacheHits     int64 `json:"cachehits,int"`
	CacheMiss     int64 `json:"cachemiss,int"`
}

//log.Println("[Cache statistics. Len:", Cache.Len(), "Cap:", Cache.Capacity(), "Hits:", Cache.Stats().Hits, "Misses:", strconv.Itoa(int(Cache.Stats().Misses))+"]")

var momo = &metrics{PostCount: 0, GetCount: 0, DelCount: 0}

func (m *metrics) incrementPost() {
	m.PostCount++
}

func (m *metrics) incrementGet() {
	m.GetCount++
}

func (m *metrics) incrementDel() {
	m.DelCount++
}

//func (m *metrics) incrementHead() {
//	m.HeadCount++
//}

func printStats() (s string) {
	runtime.ReadMemStats(&m)
	u := &metrics{}
	u.Alloc = m.Alloc
	u.Total = m.TotalAlloc
	u.System = m.Sys
	u.Gcnum = m.NumGC
	u.NextGC = m.NextGC
	u.Frees = m.Frees
	u.HeapAlloc = m.HeapAlloc
	u.HeapIdle = m.HeapIdle
	u.HeapInuse = m.HeapInuse
	u.HeapObjects = m.HeapObjects
	u.HeapReleased = m.HeapReleased
	u.PauseTotalNs = m.PauseTotalNs
	u.NumForcedGC = m.NumForcedGC
	u.Goroutines = runtime.NumGoroutine()
	u.PostCount = momo.PostCount
	u.GetCount = momo.GetCount
	//u.HeadCount = momo.HeadCount
	u.DelCount = momo.DelCount
	u.CacheLen = metadata.Cache.Len()
	u.CacheCapacity = metadata.Cache.Capacity()
	u.CacheHits = metadata.Cache.Stats().Hits
	u.CacheMiss = metadata.Cache.Stats().Misses
	result, _ := json.MarshalIndent(u, "", "    ")

	return string(result)
}
