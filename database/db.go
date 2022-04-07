package database

import (
	"database/sql"
	"log"

	"github.com/6adore/demo-groupcach/cli"
	"github.com/6adore/demo-groupcach/config"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	client := cli.NewSSHClient()
	mysql.RegisterDialContext("mysql+tcp", cli.NewSSHDialer(client).Dial)

	dsn := config.GetDsn()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("SQL connect error: ", err)
	}

	log.Printf("SQL connection success")
	return db
}
