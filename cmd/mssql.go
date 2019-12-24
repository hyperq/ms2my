package cmd

import (
	"database/sql"
	"fmt"
)

// NewMssql new db
func NewMssql(username, dbname, ip, password string, port int) (msdb *sql.DB, err error) {
	//ini := config.GetConfig(filename)
	msdb, err = sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", ip, dbname,
		username, password, port))
	return
}
