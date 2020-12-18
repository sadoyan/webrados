package wrados

import (
	"fmt"
	"github.com/ceph/go-ceph/rados"
)

type Radcon struct {
	Connection *rados.Conn
	Poolnames  map[string]bool
}

var Rconnect = &Radcon{
	Connection: RadoConnect(),
	Poolnames:  map[string]bool{},
}

func RadoConnect() *rados.Conn {

	conn, err := rados.NewConn()
	if err != nil {
		fmt.Println("Error when invoke a new connection:", err)
		return nil
	}

	err = conn.ReadDefaultConfigFile()
	if err != nil {
		fmt.Println("Error when read default config file:", err)
		return nil
	}

	err = conn.Connect()
	if err != nil {
		fmt.Println("Error when connect:", err)
		return nil
	}

	fmt.Println("Connect Ceph cluster OK!")
	return conn
}

func ListPools() {
	pools, _ := Rconnect.Connection.ListPools()
	for p := range pools {
		o := pools[p]
		Rconnect.Poolnames[o] = true
	}
}

func PutData(pool string, name string, input []byte) {
	ioctx, _ := Rconnect.Connection.OpenIOContext(pool)
	_ = ioctx.Write(name, input, 0)

	// read the data back out
	//bytesOut := make([]byte, len(input))
	//f, _ := ioctx.Read("obj", bytesOut, 0)
	//if !bytes.Equal(input, bytesOut) {
	//	fmt.Println("Output is not input!")
	//}
	//fmt.Println(f)
}
