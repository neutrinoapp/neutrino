package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/services/api/db"
)

func main() {
	s, connectErr := r.Connect(r.ConnectOpts{
		Address: config.Get(config.KEY_RETHINK_ADDR),
	})

	fmt.Println("Preparing rethinkdb")

	if connectErr != nil {
		panic(connectErr)
	}

	dbname := db.DATABASE_NAME

	fmt.Println("Creating database " + dbname)
	_, dbErr := r.DBCreate(dbname).RunWrite(s)
	if dbErr != nil {
		panic(dbErr)
	}

	fmt.Println("Creating table " + db.USERS_TABLE)
	_, usersTableError := r.DB(dbname).TableCreate(db.USERS_TABLE).RunWrite(s)
	if usersTableError != nil {
		panic(usersTableError)
	}

	fmt.Println("Creating table " + db.DATA_TABLE)
	_, dataTableError := r.DB(dbname).TableCreate(db.DATA_TABLE).RunWrite(s)
	if dataTableError != nil {
		panic(dataTableError)
	}

	fmt.Println("Done preparing rethinkdb")
}
