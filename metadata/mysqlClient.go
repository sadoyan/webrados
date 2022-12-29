package metadata

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
)

var MySQLConnection []*sql.DB

func myget(filename string) (string, error) {
	var meta string
	readQuery := "SELECT filemeta FROM files WHERE filename = ?"
	err := MySQLConnection[rand.Intn(myConns)].QueryRow(readQuery, filename).Scan(&meta)
	return meta, err
}

func myset(key string, val string) error {
	hdrstring := "filename, filemeta"
	zahrinsert := "'" + key + "','" + val + "'"
	insertString := "INSERT INTO files (" + hdrstring + ")" + "VALUES  (" + zahrinsert + ")"
	insert, ierr := MySQLConnection[rand.Intn(myConns)].Query(insertString)
	if ierr != nil {
		log.Println(ierr.Error())
	}
	_ = insert.Close()

	return nil
}

func mydel(key string) error {
	_, err := MySQLConnection[rand.Intn(myConns)].Exec("DELETE FROM files WHERE filename = ?", key)
	if err != nil {
		return err
	} else {
		return nil
	}
}
