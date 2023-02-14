package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"knightingal/section"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var db *sql.DB

type Device struct {
	ID           int    `json:"id"`
	DeviceId     string `json:"deviceId"`
	Name         string `json:"Name"`
	ManuFacturer string `json:"manuFacturer"`
}

func initDB() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "303606",
		Addr:   "localhost:3306",
		DBName: "2k",
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

func queryDevice() []Device {
	var devices []Device

	rows, err := db.Query("select id, deviceId, name, manufacturer from device")
	if err != nil {
		_ = fmt.Errorf("query failed %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var device Device
		if err := rows.Scan(&device.ID, &device.DeviceId, &device.Name, &device.ManuFacturer); err != nil {
			_ = fmt.Errorf("query failed %v", err)
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		_ = fmt.Errorf("query failed %v", err)
	}

	return devices
}

func queryDeviceById(id string) (bool, Device) {
	var device Device
	if err := db.QueryRow(
		"select id, deviceId, name, manufacturer from device where id = ?",
		id).Scan(
		&device.ID,
		&device.DeviceId,
		&device.Name,
		&device.ManuFacturer); err != nil {
		_ = fmt.Errorf("query failed %v", err)
		return false, device
	}
	return true, device
}

func getDevices(c *gin.Context) {
	var d = queryDevice()
	c.IndentedJSON(http.StatusOK, d)
}

func getDeviceById(c *gin.Context) {
	id := c.Param("id")
	succ, d := queryDeviceById(id)
	if succ {
		c.IndentedJSON(http.StatusOK, d)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{})
	}
}

func staticFileService(c *gin.Context) {
	fileName := c.Param("fileName")
	baseDir := "C:/Users/knightingal"
	target := baseDir + fileName
	_, err := os.Open(target)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{})
	}
	_, fileName = filepath.Split(fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(target)
}

func main1() {
	initDB()
	// initDB2()
	section.Init(db)
	// getDevice()
	router := gin.Default()
	router.GET("/albums", getDevices)
	router.GET("/section", section.GetSectionList)
	router.GET("/section/:id", section.GetSectionById)
	router.POST("/section", section.PostSection)
	router.GET("/albums/:id", getDeviceById)
	router.GET("/files/*fileName", staticFileService)

	router.Run("0.0.0.0:8080")
}

func main2() {
	file, err := os.Open("C:\\Users\\knightingal\\source\\go_code\\web-service-gin\\index2.html")
	if err != nil {
		fmt.Printf(err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		fmt.Printf(err.Error())
	}

	doc.Find(".f14 > img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		escape, _ := url.QueryUnescape(src)
		escape = strings.ReplaceAll(escape, "|", "")

		fmt.Println(escape)
	})

	dir := os.DirFS("C:\\Users\\knightingal\\source\\go_code\\web-service-gin")

	dirEntityList, err := fs.ReadDir(dir, ".")

	dirNames := make([]string, 0)
	for _, dir := range dirEntityList {
		// fmt.Println(dir.Name())
		dirNames = append(dirNames, dir.Name())
	}
	fmt.Println(dirNames)

	cb := func(src string) (string, bool) {

		filterRet := filter(&dirNames, func(dirName string) bool {
			return strings.Contains(dirName, src)
		})

		if len(*filterRet) == 1 {
			fmt.Println((*filterRet)[0])
			return (*filterRet)[0], true
		}

		return "", false
	}
	const srcString = "2222index22.html11111"
	windowString(srcString, cb)

}

func windowString(src string, process func(string) (string, bool)) {
	srcArray := []rune(src)
	size := len(srcArray)
	stop := false
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			sub1 := srcArray[j : j+size-i]
			fmt.Println(string(sub1))
			_, stop = process(string(sub1))
			if stop {
				break
			}
		}
		if stop {
			break
		}
	}
}

func filter[T any](src *[]T, fn func(T) bool) *[]T {
	ret := make([]T, 0)
	for _, item := range *src {
		if fn(item) {
			ret = append(ret, item)
		}
	}
	return &ret
}
