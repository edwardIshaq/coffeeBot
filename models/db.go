package models

import (
	"database/sql"
	"fmt"
)

// Private package varible, injected in by calling `SetDatabase`
var db *sql.DB

// SetDatabase sets up the DB
func SetDatabase(database *sql.DB) {
	db = database
}

// DemoDB simple demo
func DemoDB() {
	println("we have db? :", db)
	row := db.QueryRow(`SELECT team_id, access_token FROM team WHERE id = $1`, 1)
	var teamID string
	var accessToken string
	row.Scan(&teamID, &accessToken)
	fmt.Printf("accessToken: %s | team: %s", accessToken, teamID)
}

// DBWrapper will be a public interface for other sibling packages
type DBWrapper struct {
	db *sql.DB
}

// NewDBWrapper sets up a new DB
func NewDBWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{
		db: db,
	}
}

func (dbWrapper *DBWrapper) SaveToDB(value interface{}) {
	fmt.Printf("DBWrapper: saving to DB: %v", dbWrapper)
}
