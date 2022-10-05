package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Section struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	BaseDir string `json:"baseDir"`
}

func getSectionList(c *gin.Context) {
	var d = querySectionList()
	c.IndentedJSON(http.StatusOK, d)
}

func postSection(c *gin.Context) {
	var section Section

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&section); err != nil {
		return
	}

	insertSection(section)

	// Add the new album to the slice.
	c.IndentedJSON(http.StatusCreated, nil)
}

func insertSection(section Section) {
	id, _ := uuid.NewUUID()

	section.ID = strings.Replace(id.String(), "-", "", -1)

	result, err := db.Exec("insert into section(id, name, base_dir) values (?, ?, ?)",
		section.ID,
		section.Name,
		section.BaseDir)

	if err != nil {
		_ = fmt.Errorf("insert failed %v", err)
	}

	fmt.Println(result.LastInsertId())

	return
}

func querySectionList() []Section {
	var sectionList []Section

	rows, err := db.Query("select id, name, base_dir from section")
	if err != nil {
		_ = fmt.Errorf("query failed %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var section Section
		if err := rows.Scan(&section.ID, &section.Name, &section.BaseDir); err != nil {
			_ = fmt.Errorf("query failed %v", err)
		}
		sectionList = append(sectionList, section)
	}
	if err := rows.Err(); err != nil {
		_ = fmt.Errorf("query failed %v", err)
	}

	return sectionList
}
