package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/neutrinoapp/neutrino/src/common/config"
	"github.com/neutrinoapp/neutrino/src/common/db"
)

func main() {
	s, connectErr := r.Connect(r.ConnectOpts{
		Address: config.Get(config.KEY_RETHINK_ADDR),
	})

	fmt.Println("Preparing rethinkdb")

	if connectErr != nil {
		fmt.Println(connectErr)
	}

	dbname := db.DATABASE_NAME

	fmt.Println("Creating database " + dbname)
	_, dbErr := r.DBCreate(dbname).RunWrite(s)
	if dbErr != nil {
		fmt.Println(dbErr)
	}

	fmt.Println("Creating table " + db.USERS_TABLE)
	_, usersTableError := r.DB(dbname).TableCreate(db.USERS_TABLE).RunWrite(s)
	if usersTableError != nil {
		fmt.Println(usersTableError)
	}

	fmt.Println("Creating indexes for " + db.USERS_TABLE)
	_, createEmailIndexError := r.DB(dbname).Table(db.USERS_TABLE).IndexCreate(db.USERS_TABLE_EMAIL_INDEX).RunWrite(s)
	if createEmailIndexError != nil {
		fmt.Println(createEmailIndexError)
	}

	fmt.Println("Creating table " + db.DATA_TABLE)
	_, dataTableError := r.DB(dbname).TableCreate(db.DATA_TABLE).RunWrite(s)
	if dataTableError != nil {
		fmt.Println(dataTableError)
	}

	fmt.Println("Creating index for table " + db.DATA_TABLE)
	_, createDataEmailAppidIndexError := r.DB(dbname).Table(db.DATA_TABLE).
		IndexCreateFunc(db.DATA_TABLE_APPIDTYPE_INDEX, func(row r.Term) interface{} {
			return []interface{}{row.Field(db.APP_ID_FIELD), row.Field(db.TYPE_FIELD)}
		}).RunWrite(s)
	if createDataEmailAppidIndexError != nil {
		fmt.Println(createDataEmailAppidIndexError)
	}

	fmt.Println("Creating table " + db.APPS_TABLE)
	_, createAppTableError := r.DB(dbname).TableCreate(db.APPS_TABLE).RunWrite(s)
	if createAppTableError != nil {
		fmt.Println(createAppTableError)
	}

	fmt.Println("Creating table " + db.APPS_USERS_TABLE)
	_, createAppsUsersTableError := r.DB(dbname).TableCreate(db.APPS_USERS_TABLE).RunWrite(s)
	if createAppsUsersTableError != nil {
		fmt.Println(createAppsUsersTableError)
	}

	fmt.Println("Creating index for table " + db.APPS_USERS_TABLE)
	_, createAppsUsersTableIndexError := r.DB(dbname).Table(db.APPS_USERS_TABLE).
		IndexCreateFunc(db.APPS_USERS_TABLE_EMAILAPPID_INDEX, func(row r.Term) interface{} {
			return []interface{}{row.Field(db.EMAIL_FIELD), row.Field(db.APP_ID_FIELD)}
		}).RunWrite(s)
	if createAppsUsersTableIndexError != nil {
		fmt.Println(createAppsUsersTableIndexError)
	}

	fmt.Println("Done preparing rethinkdb")
}
