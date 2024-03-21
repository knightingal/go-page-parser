package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func initFlowDB() {
	cfg := mysql.Config{
		User:                 "root",
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
		"album, cover, cover_height, cover_width, create_time, dir_name, name, client_status"+
		") values (?,?,?,?,?,?,?,?)", section.destAlbum,
		section.cover.binName,
		section.cover.height, section.cover.width, section.timeStamp, section.name, section.name, section.clientStatus)

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

func updateComment(fileName string, comment string) {
	r := db.QueryRow("select comment from flow1000log where file_name = ?", fileName)
	existComment := ""
	r.Scan(&existComment)
	existComment = existComment + "," + comment
	ret, error := db.Exec("update flow1000log set comment = ? where file_name = ?", existComment, fileName)
	if error != nil {
		log.Fatal(error)
	}
	log.Println(ret)
}

func updateLog(fileName string, msg string) {
	_, error := db.Exec("update flow1000log set msg = ? where file_name = ?", msg, fileName)
	if error != nil {
		log.Fatal(error)
	}
}

func checkSuccLog(fileName string) bool {
	r := db.QueryRow("select count(file_name) from flow1000log where file_name = ? and msg = 'succ'", fileName)
	var count int
	r.Scan(&count)
	return count > 0
}

func checkExistLog(fileName string) bool {
	r := db.QueryRow("select count(file_name) from flow1000log where file_name = ? ", fileName)
	var count int
	r.Scan(&count)
	return count > 0
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

func insertImgBin(image Image, sectionId int64) {
	result, error := db.Exec("insert into flow1000img("+
		"name, height, width, in_cover, section_id"+
		") values (?,?,?,?,?)",
		image.binName, image.height, image.width, 0, sectionId)

	if error != nil {
		log.Fatal(error)
	}

	insertId, _ := result.LastInsertId()
	fmt.Printf("insert %d\n", insertId)

}
