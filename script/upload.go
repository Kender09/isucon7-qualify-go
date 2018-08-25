package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	avatarMaxBytes = 1 * 1024 * 1024
)

var (
	db *sqlx.DB
)

func init() {
	seedBuf := make([]byte, 8)
	crand.Read(seedBuf)
	rand.Seed(int64(binary.LittleEndian.Uint64(seedBuf)))

	db_host := "localhost"
	db_port := "3306"
	db_user := "isucon"
	db_password := ":isucon"

	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isubata",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Printf("Connecting to error: %s", err)
	}
	log.Printf("Succeeded to connect db.")
}

func writeIcons() error {
	rows, err := db.Query("SELECT name, data FROM image WHERE id <= 1000")
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {
		var name string
		var data []byte
		if err := rows.Scan(&name, &data); err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf(`/home/isucon/static/%s`, name))
		if err != nil {
			return err
		}
		defer file.Close()

		file.Write(data)
	}
	return nil
}

func main() {
	if err := writeIcons(); err != nil {
		log.Printf("err: %s", err)
	}
}
