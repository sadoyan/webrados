package metadata

import (
	"fmt"
	"math/rand"
	"strings"
	"tools"
	"wrados"
)

func readTheMeta(data []byte) []byte {
	fmt.Println(data)
	return data
}

func cephget(filename string) (string, error) {
	md := strings.Split(filename, "/")
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(md[0])
	//defer ioctx.Destroy()
	if e != nil {
		tools.WriteLogs("Metadata direct read error:", e)
		return "Metadata direct read error:", e
	}
	//_ = ioctx.SetXattr(md[1], "segments", []byte("h.mp4-0,h.mp4-132169727,h.mp4-264339454,h.mp4-396509181,h.mp4-528678908,547579107"))
	ss := make([]byte, 4096)
	ddd, _ := ioctx.GetXattr(md[1], "segments", ss)
	return (string(ss[:ddd])), nil
}

func cephset(filename string, metadata string) (string, error) {
	md := strings.Split(filename, "/")
	randindex := rand.Intn(len(wrados.Rconnect.Connection))
	ioctx, e := wrados.Rconnect.Connection[randindex].OpenIOContext(md[0])
	//defer ioctx.Destroy()
	if e != nil {
		tools.WriteLogs("Metadata direct write error:", e)
		return "Metadata direct write error:", e
	}
	_ = ioctx.SetXattr(md[1], "segments", []byte(metadata))
	return "Done", nil
}
