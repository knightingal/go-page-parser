package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func initFlowDB() {
	cfg := mysql.Config{
		User:                 "knightingal",
		Passwd:               "000000",
		Addr:                 "localhost:3306",
		DBName:               "flow1000",
		AllowNativePasswords: true,
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")

}

func insertSection(section Section) int64 {
	result, error := db.Exec("insert into flow1000section("+
		"album, cover, cover_height, cover_width, create_time, dir_name, name"+
		") values (?,?,?,?,?,?,?)", section.album,
		section.cover.name,
		section.cover.height, section.cover.width, section.timeStamp, section.name, section.name)

	if error != nil {
		log.Fatal(error)
	}

	insertId, _ := result.LastInsertId()
	fmt.Printf("insert %d", insertId)

	return insertId

}

func insertLog(fileName string, msg string) {
	_, error := db.Exec("insert into flow1000log(file_name, msg) values (?,?)", fileName, msg)
	if error != nil {
		log.Fatal(error)
	}
}

func updateLog(fileName string, msg string) {
	_, error := db.Exec("update flow1000log set msg = ? where file_name = ?", msg, fileName)
	if error != nil {
		log.Fatal(error)
	}

}

func insertImg(image Image, sectionId int64) {
	result, error := db.Exec("insert into flow1000img("+
		"name, height, width, in_cover, section_id"+
		") values (?,?,?,?,?)",
		image.name, image.height, image.width, 0, sectionId)

	if error != nil {
		log.Fatal(error)
	}

	insertId, _ := result.LastInsertId()
	fmt.Printf("insert %d\n", insertId)

}
