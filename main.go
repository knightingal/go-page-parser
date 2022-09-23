package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
		DBName: "wvp",
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

func main() {
	initDB()
	initDB2()
	// getDevice()
	router := gin.Default()
	router.GET("/albums", getDevices)
	router.GET("/albums/:id", getDeviceById)

	router.Run("0.0.0.0:8080")
}
