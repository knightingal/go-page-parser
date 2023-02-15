package main

import (
	"database/sql"
	"fmt"
	"io"
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

const BASE_DIR = "/mnt/download/"
const TARGET_DIR = "/mnt/linux1000/1024/"

func main() {
	// file, err := os.Open("/mnt/2048/CLImages2/[西川康] お嬢様は戀話がお好き.html")
	file, err := os.Open(BASE_DIR + "輝夜姬想讓人告白_天才們的戀愛頭腦戰_ 早坂愛 2.html")
	if err != nil {
		fmt.Printf(err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		fmt.Printf(err.Error())
	}

	var src string
	var srcDir string

	imgSrcList := make([]string, 0)

	doc.Find(".tpc_content").Each(func(i int, s *goquery.Selection) {
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ = s.Attr("src")
			escape, _ := url.QueryUnescape(src)

			fmt.Println(escape)
			srcDirList := strings.Split(escape, "/")
			srcDir = srcDirList[len(srcDirList)-2]
			src = srcDirList[len(srcDirList)-1]
			imgSrcList = append(imgSrcList, src)
		})
	})

	dir := os.DirFS(BASE_DIR)

	dirEntityList, _ := fs.ReadDir(dir, ".")

	dirNames := make([]string, 0)
	for _, dir := range dirEntityList {
		if dir.IsDir() {
			dirNames = append(dirNames, dir.Name())
		}
	}
	fmt.Println(dirNames)

	cb := func(src string) (string, bool) {

		filterRet := filter(&dirNames, func(dirName string) bool {
			return strings.Contains(dirName, src)
		})

		if len(*filterRet) == 1 {
			fmt.Println("====matched====")
			fmt.Println((*filterRet)[0])
			return (*filterRet)[0], true
		}

		return "", false
	}
	realDir, succ := windowString(srcDir, cb)
	if !succ {
		return
	}
	fmt.Println(realDir)

	os.Mkdir(TARGET_DIR+realDir, 0750)
	for _, imgSrc := range imgSrcList {
		targetFile, _ := os.Create(TARGET_DIR + realDir + "/" + imgSrc)
		srcFile, _ := os.Open(BASE_DIR + realDir + "/" + imgSrc)
		io.Copy(targetFile, srcFile)

	}

}

func windowString(src string, process func(string) (string, bool)) (string, bool) {
	srcArray := []rune(src)
	size := len(srcArray)
	stop := false
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			sub1 := srcArray[j : j+size-i]
			fmt.Println(string(sub1))
			realDir, stop := process(string(sub1))
			if stop {
				return realDir, true
			}
		}
		if stop {
			break
		}
	}
	return "", false
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
