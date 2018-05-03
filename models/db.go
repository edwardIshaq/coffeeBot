package models

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

// SetDatabase sets up the DB
func SetDatabase(database *sql.DB) {
	db = database
}

// DemoDB simple demo
func DemoDB() {
	println("we have db? :", db)
	row := db.QueryRow(`SELECT "teamID", "accessToken" FROM team WHERE id = $1`, 1)
	var teamID string
	var accessToken string
	row.Scan(&teamID, &accessToken)
	fmt.Printf("accessToken: %s | team: %s", accessToken, teamID)
}
